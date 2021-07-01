package media

import (
	"github.com/Nistagram-Organization/nistagram-media/src/repositories/media"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
)

type MediaService interface {
	Create(*model.Media) (*model.Media, rest_error.RestErr)
}

type mediaService struct {
	mediaRepository media.MediaRepository
}

func NewMediaService(mediaRepository media.MediaRepository) MediaService {
	return &mediaService{
		mediaRepository: mediaRepository,
	}
}

func (s *mediaService) Create(media *model.Media) (*model.Media, rest_error.RestErr) {
	return s.mediaRepository.Create(media)
}