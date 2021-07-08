package application

import (
	"github.com/Nistagram-Organization/nistagram-media/src/datasources/mysql"
	media2 "github.com/Nistagram-Organization/nistagram-media/src/repositories/media"
	"github.com/Nistagram-Organization/nistagram-media/src/services/media"
	"github.com/Nistagram-Organization/nistagram-media/src/services/media_grpc_service"
	"github.com/Nistagram-Organization/nistagram-media/src/utils/image_utils"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/prometheus_handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	router        = gin.Default()
	requestsCount = prometheus_handler.GetHttpRequestsCounter()
	requestsSize  = prometheus_handler.GetHttpRequestsSize()
	uniqueUsers   = prometheus_handler.GetUniqueClients()
)

func configureCORS() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))
}

func setupDatabase() (datasources.DatabaseClient, error) {
	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		return nil, err
	}
	if err := database.Migrate(
		&model.Media{},
	); err != nil {
		return nil, err
	}
	return database, nil
}

func registerPrometheusMiddleware() {
	prometheus.Register(requestsCount)
	prometheus.Register(requestsSize)
	prometheus.Register(uniqueUsers)

	router.Use(prometheus_handler.PrometheusMiddleware(requestsCount, requestsSize, uniqueUsers))
}

func StartApplication() {
	configureCORS()
	registerPrometheusMiddleware()

	database, err := setupDatabase()
	if err != nil {
		panic(err)
	}

	port := ":8089"
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	mediaRepository := media2.NewMediaRepository(database)
	imageUtilsService := image_utils.NewImageUtilsService()
	mediaService := media.NewMediaService(mediaRepository)
	mediaGrpcService := media_grpc_service.NewMediaGrpcService(mediaService, imageUtilsService, "temp")

	grpcS := grpc.NewServer()
	proto.RegisterMediaServiceServer(grpcS, mediaGrpcService)

	router.GET("/metrics", prometheus_handler.PrometheusGinHandler())

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
