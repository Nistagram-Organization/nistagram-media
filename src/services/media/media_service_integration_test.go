package media

import (
	"github.com/Nistagram-Organization/nistagram-media/src/datasources/mysql"
	"github.com/Nistagram-Organization/nistagram-media/src/repositories/media"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type MediaServiceIntegrationTestsSuite struct {
	suite.Suite
	service MediaService
	db      *gorm.DB
}

func (suite *MediaServiceIntegrationTestsSuite) SetupSuite() {
	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}
	if err := database.Migrate(
		&model.Media{},
	); err != nil {
		panic(err)
	}

	suite.db = database.GetClient()

	mediaRepository := media.NewMediaRepository(database)
	suite.service = NewMediaService(mediaRepository)
}

func (suite *MediaServiceIntegrationTestsSuite) SetupTest() {
	mediaEntity := model.Media{
		ID:   1,
		Path: "path",
	}

	tx := suite.db.Begin()
	tx.Create(mediaEntity)
	tx.Commit()
}

func TestMediaServiceIntegrationTestsSuite(t *testing.T) {
	suite.Run(t, new(MediaServiceIntegrationTestsSuite))
}

func (suite *MediaServiceIntegrationTestsSuite) TestIntegrationMediaService_Create() {
	mediaEntity := model.Media{
		Path: "temp/slika.jpg",
	}

	retMedia, createErr := suite.service.Create(&mediaEntity)

	assert.Equal(suite.T(), &mediaEntity, retMedia)
	assert.Equal(suite.T(), nil, createErr)
}

func (suite *MediaServiceIntegrationTestsSuite) TestIntegrationMediaService_Get() {
	retMedia, createErr := suite.service.Get(1)

	assert.Equal(suite.T(), uint(1), retMedia.ID)
	assert.Equal(suite.T(), nil, createErr)
}
