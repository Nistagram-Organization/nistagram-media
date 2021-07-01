package application

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-media/src/datasources/mysql"
	media2 "github.com/Nistagram-Organization/nistagram-media/src/repositories/media"
	"github.com/Nistagram-Organization/nistagram-media/src/services/media"
	"github.com/Nistagram-Organization/nistagram-media/src/services/media_grpc_service"
	"github.com/Nistagram-Organization/nistagram-media/src/utils/image_utils"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	router = gin.Default()
)

func StartApplication() {
	router.Use(cors.Default())

	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}
	if err := database.Migrate(
		&model.Media{},
	); err != nil {
		panic(err)
	}

	port := ":8089"
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1%s", port))
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	grpcS := grpc.NewServer()
	proto.RegisterMediaServiceServer(grpcS,
		media_grpc_service.NewMediaGrpcService(
			media.NewMediaService(
				media2.NewMediaRepository(database)),
			image_utils.NewImageUtilsService(),
			"temp",
		),
	)

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
