package media

import (
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/media"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type MediaRepositoryMock struct {
	mock.Mock
}

func (m *MediaRepositoryMock) Create(media *model.Media) (*model.Media, rest_error.RestErr) {
	args := m.Called(media)
	if args.Get(1) == nil {
		return args.Get(0).(*model.Media), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (m *MediaRepositoryMock) Get(u uint) (*model.Media, rest_error.RestErr) {
	panic("implement me")
}

