package service

import (
	"context"
	stdpath "path"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/singleflight"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/Xhofe/go-cache"

	log "github.com/dsoprea/go-logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type FsListService interface {
	FsList(ctx context.Context, path string, refresh ...bool) ([]model.Obj, error)
	Key(storage driver.Driver, path string) string
	Get(ctx context.Context, storage driver.Driver, path string) (model.Obj, error)
	GetUnwrap(ctx context.Context, storage driver.Driver, path string) (model.Obj, error)
	List(ctx context.Context, storage driver.Driver, path string, args model.ListArgs, refresh ...bool) ([]model.Obj, error)
}

type fsListService struct {
}

var listCache = cache.NewMemCache(cache.WithShards[[]model.Obj](64))
var listG singleflight.Group[[]model.Obj]

// List files
func (fl *fsListService) FsList(ctx context.Context, path string, refresh ...bool) ([]model.Obj, error) {

	virtualFiles := MyService.Storages().GetStorageVirtualFilesByPath(path)
	storage, actualPath, err := MyService.StoragePath().GetStorageAndActualPath(path)
	if err != nil && len(virtualFiles) == 0 {
		return nil, errors.WithMessage(err, "failed get storage")
	}

	var _objs []model.Obj
	if storage != nil {
		_objs, err = fl.List(ctx, storage, actualPath, model.ListArgs{
			ReqPath: path,
		}, refresh...)
		if err != nil {
			log.Errorf("%+v", err)
			if len(virtualFiles) == 0 {
				return nil, errors.WithMessage(err, "failed get objs")
			}
		}
	}

	om := model.NewObjMerge()

	objs := om.Merge(virtualFiles, _objs...)
	return objs, nil
}

func (fl *fsListService) Key(storage driver.Driver, path string) string {
	return stdpath.Join(storage.GetStorage().MountPath, utils.FixAndCleanPath(path))
}

// Get object from list of files
func (fl *fsListService) Get(ctx context.Context, storage driver.Driver, path string) (model.Obj, error) {
	path = utils.FixAndCleanPath(path)
	logger.Info("get", zap.String("path", path))

	// is root folder
	if utils.PathEqual(path, "/") {
		var rootObj model.Obj
		switch r := storage.GetAddition().(type) {
		case driver.IRootId:
			rootObj = &model.Object{
				ID:       r.GetRootId(),
				Name:     op.RootName,
				Size:     0,
				Modified: storage.GetStorage().Modified,
				IsFolder: true,
			}
		case driver.IRootPath:
			rootObj = &model.Object{
				Path:     r.GetRootPath(),
				Name:     op.RootName,
				Size:     0,
				Modified: storage.GetStorage().Modified,
				IsFolder: true,
			}
		default:
			if storage, ok := storage.(driver.Getter); ok {
				obj, err := storage.GetRoot(ctx)
				if err != nil {
					return nil, errors.WithMessage(err, "failed get root obj")
				}
				rootObj = obj
			}
		}
		if rootObj == nil {
			return nil, errors.Errorf("please implement IRootPath or IRootId or Getter method")
		}
		return &model.ObjWrapName{
			Name: op.RootName,
			Obj:  rootObj,
		}, nil
	}

	// not root folder
	dir, name := stdpath.Split(path)
	files, err := fl.List(ctx, storage, dir, model.ListArgs{})
	if err != nil {
		return nil, errors.WithMessage(err, "failed get parent list")
	}
	for _, f := range files {
		// TODO maybe copy obj here
		if f.GetName() == name {
			return f, nil
		}
	}
	logger.Info("cant find obj with name", zap.Any("name", name))
	return nil, errors.WithStack(errors.New("object not found"))
}

func (fl *fsListService) GetUnwrap(ctx context.Context, storage driver.Driver, path string) (model.Obj, error) {
	obj, err := fl.Get(ctx, storage, path)
	if err != nil {
		return nil, err
	}
	return model.UnwrapObjs(obj), err
}

// List files in storage, not contains virtual file
func (fl *fsListService) List(ctx context.Context, storage driver.Driver, path string, args model.ListArgs, refresh ...bool) ([]model.Obj, error) {
	if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
		return nil, errors.Errorf("storage not init: %s", storage.GetStorage().Status)
	}
	path = utils.FixAndCleanPath(path)
	logger.Info("op.List", zap.Any("path", path))
	key := fl.Key(storage, path)
	if !utils.IsBool(refresh...) {
		if files, ok := listCache.Get(key); ok {
			logger.Info("op.List", zap.Any("use cache", path))
			return files, nil
		}
	}
	dir, err := fl.GetUnwrap(ctx, storage, path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed get dir")
	}
	logger.Info("op.List", zap.Any("dir", dir))
	if !dir.IsDir() {
		return nil, errors.WithStack(errors.New("not a folder"))
	}
	objs, err, _ := listG.Do(key, func() ([]model.Obj, error) {
		files, err := storage.List(ctx, dir, args)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list objs")
		}
		// set path
		for _, f := range files {
			if s, ok := f.(model.SetPath); ok && f.GetPath() == "" && dir.GetPath() != "" {
				s.SetPath(stdpath.Join(dir.GetPath(), f.GetName()))
			}
		}
		// warp obj name
		model.WrapObjsName(files)
		// call hooks
		go func(reqPath string, files []model.Obj) {
			for _, hook := range op.ObjsUpdateHooks {
				hook(args.ReqPath, files)
			}
		}(args.ReqPath, files)

		// sort objs
		if storage.Config().LocalSort {
			model.SortFiles(files, storage.GetStorage().OrderBy, storage.GetStorage().OrderDirection)
		}
		model.ExtractFolder(files, storage.GetStorage().ExtractFolder)

		if !storage.Config().NoCache {
			if len(files) > 0 {
				logger.Info("set cache", zap.Any("key", key), zap.Any("files", files))
				listCache.Set(key, files, cache.WithEx[[]model.Obj](time.Minute*time.Duration(storage.GetStorage().CacheExpiration)))
			} else {
				logger.Info("del cache", zap.Any("key", key))
				listCache.Del(key)
			}
		}
		return files, nil
	})
	return objs, err
}

func NewFsListService() FsListService {
	return &fsListService{}
}
