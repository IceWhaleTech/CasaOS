package service

import (
	"context"

	"github.com/IceWhaleTech/CasaOS/internal/op"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/pkg/errors"
)

type FsLinkService interface {
	Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error)
}

type fsLinkService struct {
}

func (f *fsLinkService) Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error) {
	storage, actualPath, err := MyService.StoragePath().GetStorageAndActualPath(path)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed get storage")
	}
	return op.Link(ctx, storage, actualPath, args)
}
func NewFsLinkService() FsLinkService {
	return &fsLinkService{}
}
