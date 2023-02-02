package service

import (
	"context"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/model"
	log "github.com/dsoprea/go-logging"
	"go.uber.org/zap"
)

type FsService interface {
	FList(ctx context.Context, path string, refresh ...bool) ([]model.Obj, error)
	GetStorage(path string) (driver.Driver, error)
	Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error)
}

type fsService struct {
}

// the param named path of functions in this package is a mount path
// So, the purpose of this package is to convert mount path to actual path
// then pass the actual path to the op package

func (f *fsService) FList(ctx context.Context, path string, refresh ...bool) ([]model.Obj, error) {
	res, err := MyService.FsListService().FsList(ctx, path, refresh...)
	if err != nil {
		logger.Info("failed list", zap.Any("path", path), zap.Any("err", err))
		return nil, err
	}
	return res, nil
}

// func (f *fsService) Get(ctx context.Context, path string) (model.Obj, error) {
// 	res, err := get(ctx, path)
// 	if err != nil {
// 		log.Errorf("failed get %s: %+v", path, err)
// 		return nil, err
// 	}
// 	return res, nil
// }

func (f *fsService) Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error) {
	res, file, err := MyService.FsLinkService().Link(ctx, path, args)
	if err != nil {
		log.Errorf("failed link %s: %+v", path, err)
		return nil, nil, err
	}
	return res, file, nil
}

// func (f *fsService) MakeDir(ctx context.Context, path string, lazyCache ...bool) error {
// 	err := makeDir(ctx, path, lazyCache...)
// 	if err != nil {
// 		log.Errorf("failed make dir %s: %+v", path, err)
// 	}
// 	return err
// }

// func (f *fsService) Move(ctx context.Context, srcPath, dstDirPath string, lazyCache ...bool) error {
// 	err := move(ctx, srcPath, dstDirPath, lazyCache...)
// 	if err != nil {
// 		log.Errorf("failed move %s to %s: %+v", srcPath, dstDirPath, err)
// 	}
// 	return err
// }

// func (f *fsService) Copy(ctx context.Context, srcObjPath, dstDirPath string, lazyCache ...bool) (bool, error) {
// 	res, err := _copy(ctx, srcObjPath, dstDirPath, lazyCache...)
// 	if err != nil {
// 		log.Errorf("failed copy %s to %s: %+v", srcObjPath, dstDirPath, err)
// 	}
// 	return res, err
// }

// func (f *fsService) Rename(ctx context.Context, srcPath, dstName string, lazyCache ...bool) error {
// 	err := rename(ctx, srcPath, dstName, lazyCache...)
// 	if err != nil {
// 		log.Errorf("failed rename %s to %s: %+v", srcPath, dstName, err)
// 	}
// 	return err
// }

// func (f *fsService) Remove(ctx context.Context, path string) error {
// 	err := remove(ctx, path)
// 	if err != nil {
// 		log.Errorf("failed remove %s: %+v", path, err)
// 	}
// 	return err
// }

// func PutDirectly(ctx context.Context, dstDirPath string, file *model.FileStream, lazyCache ...bool) error {
// 	err := putDirectly(ctx, dstDirPath, file, lazyCache...)
// 	if err != nil {
// 		log.Errorf("failed put %s: %+v", dstDirPath, err)
// 	}
// 	return err
// }

// func (f *fsService) PutAsTask(dstDirPath string, file *model.FileStream) error {
// 	err := putAsTask(dstDirPath, file)
// 	if err != nil {
// 		log.Errorf("failed put %s: %+v", dstDirPath, err)
// 	}
// 	return err
// }

func (f *fsService) GetStorage(path string) (driver.Driver, error) {
	storageDriver, _, err := MyService.StoragePath().GetStorageAndActualPath(path)
	if err != nil {
		return nil, err
	}
	return storageDriver, nil
}

// func (f *fsService) Other(ctx context.Context, args model.FsOtherArgs) (interface{}, error) {
// 	res, err := other(ctx, args)
// 	if err != nil {
// 		log.Errorf("failed remove %s: %+v", args.Path, err)
// 	}
// 	return res, err
// }

// func get(ctx context.Context, path string) (model.Obj, error) {
// 	path = utils.FixAndCleanPath(path)
// 	// maybe a virtual file
// 	if path != "/" {
// 		virtualFiles := op.GetStorageVirtualFilesByPath(stdpath.Dir(path))
// 		for _, f := range virtualFiles {
// 			if f.GetName() == stdpath.Base(path) {
// 				return f, nil
// 			}
// 		}
// 	}
// 	storage, actualPath, err := op.GetStorageAndActualPath(path)
// 	if err != nil {
// 		// if there are no storage prefix with path, maybe root folder
// 		if path == "/" {
// 			return &model.Object{
// 				Name:     "root",
// 				Size:     0,
// 				Modified: time.Time{},
// 				IsFolder: true,
// 			}, nil
// 		}
// 		return nil, errors.WithMessage(err, "failed get storage")
// 	}
// 	return op.Get(ctx, storage, actualPath)
// }

func NewFsService() FsService {
	return &fsService{}
}
