package media

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(*media.Media) (*media.Media, rest_error.RestErr)
	Get(uint) (*media.Media, rest_error.RestErr)
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(databaseClient datasources.DatabaseClient) MediaRepository {
	return &mediaRepository{
		databaseClient.GetClient(),
	}
}

func (m *mediaRepository) Create(media *media.Media) (*media.Media, rest_error.RestErr) {
	if err := m.db.Create(media).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to save media", err)
	}
	return media, nil
}

func (m *mediaRepository) Get(id uint) (*media.Media, rest_error.RestErr) {
	mediaEntity := media.Media{
		ID: id,
	}

	if err := m.db.Take(&mediaEntity, mediaEntity.ID).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get media with id %d", mediaEntity.ID))
	}

	return &mediaEntity, nil
}
