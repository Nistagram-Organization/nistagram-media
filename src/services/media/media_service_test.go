package media

import (
	"github.com/Nistagram-Organization/nistagram-media/src/repositories/media"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MediaServiceUnitTestsSuite struct {
	suite.Suite
	mediaRepositoryMock *media.MediaRepositoryMock
	service                MediaService
}

func TestMediaServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(MediaServiceUnitTestsSuite))
}

func (suite *MediaServiceUnitTestsSuite) SetupSuite() {
	suite.mediaRepositoryMock = new(media.MediaRepositoryMock)
	suite.service = NewMediaService(suite.mediaRepositoryMock)
}

func (suite *MediaServiceUnitTestsSuite) TestNewMediaService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *MediaServiceUnitTestsSuite) TestMediaService_Create() {
	mediaEntity := model.Media{
		Path: "temp/slika.jpg",
	}

	suite.mediaRepositoryMock.On("Create", &mediaEntity).Return(&mediaEntity, nil)

	retMedia, createErr := suite.service.Create(&mediaEntity)

	assert.Equal(suite.T(), &mediaEntity, retMedia)
	assert.Equal(suite.T(), nil, createErr)
}


