package media

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(*media.Media) (*media.Media, rest_error.RestErr)
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
