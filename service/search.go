package service

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/types"
)

type SearchService interface {
	SearchList(key string) ([]model.SearchFileInfo, error)
}

type searchService struct {
}

func (s *searchService) SearchList(key string) ([]model.SearchFileInfo, error) {
	pathName := "/Users/liangjianli/go/CasaOSNew/searchTest"
	resArr := []model.SearchFileInfo{}
	files, _ := ioutil.ReadDir(pathName)
	for _, file := range files {
		if file.IsDir() {
			tempArr, err := s.SearchList(pathName + "/" + file.Name())
			if err != nil {
				resArr = append(resArr, tempArr...)
			}
		} else {
			if strings.Contains(file.Name(), key) {
				resArr = append(resArr, model.SearchFileInfo{Path: pathName, Name: file.Name(), Type: GetSearchType(path.Ext(file.Name()))})
			}
			fmt.Println(pathName + "/" + file.Name())
		}
	}
	return resArr, nil
}

func GetSearchType(ext string) int {
	var reType int = types.UNKNOWN
	switch ext {
	case ".png":
		reType = types.PICTURE
	case ".mp4":
		reType = types.MEDIA
	case ".mp3":
		reType = types.MUSIC
	default:
		reType = types.UNKNOWN
	}
	return reType
}

func NewSearchService() SearchService {
	return &searchService{}
}
