package media_grpc_service

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-media/src/services/media"
	"github.com/Nistagram-Organization/nistagram-media/src/utils/image_utils"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
)

type mediaGrpcService struct {
	proto.MediaServiceServer
	mediaService      media.MediaService
	imageUtilsService image_utils.ImageUtilsService
	tempFolder        string
}

func NewMediaGrpcService(mediaService media.MediaService, imageUtilsService image_utils.ImageUtilsService, tempFolder string) proto.MediaServiceServer {
	return &mediaGrpcService{
		proto.UnimplementedMediaServiceServer{},
		mediaService,
		imageUtilsService,
		tempFolder,
	}
}

func (s *mediaGrpcService) SaveMedia(ctx context.Context, saveMediaRequest *proto.SaveMediaRequest) (*proto.SaveMediaResponse, error) {
	mediaMessage := saveMediaRequest.GetImage()

	imagePath, err := s.imageUtilsService.SaveImage(mediaMessage.ImageBase64, s.tempFolder)
	if err != nil {
		return nil, err
	}

	mediaEntity := &model.Media{
		Path: imagePath,
	}

	mediaEntity, err = s.mediaService.Create(mediaEntity)
	if err != nil {
		return nil, err
	}

	res := proto.SaveMediaResponse{Id: uint64(mediaEntity.ID)}

	return &res, nil
}
