package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

type FileInfo struct {
	init             bool
	uploaded         []bool
	uploadedChunkNum int64
}

type FileUploadService struct {
	uploadStatus sync.Map
	lock         sync.RWMutex
}

func NewFileUploadService() *FileUploadService {
	return &FileUploadService{
		uploadStatus: sync.Map{},
		lock:         sync.RWMutex{},
	}
}

func (s *FileUploadService) TestChunk(c echo.Context) error {
	// s.lock.RLock()
	// defer s.lock.RUnlock()

	identifier := c.QueryParam("identifier")
	chunkNumber, err := strconv.ParseInt(c.QueryParam("chunkNumber"), 10, 64)
	if err != nil {
		return err
	}
	fileInfoTemp, ok := s.uploadStatus.Load(identifier)

	if !ok {
		return c.NoContent(http.StatusNoContent)
	}

	fileInfo := fileInfoTemp.(*FileInfo)

	if !fileInfo.init {
		return c.NoContent(http.StatusNoContent)
	}

	// 这里返回的应该得是 permanentErrors，不是 404. 不然前端会上传失败而不是重传块。
	// 梁哥应该得改一下。
	if !fileInfo.uploaded[chunkNumber-1] {
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusOK)
}

func (s *FileUploadService) UploadFile(c echo.Context) error {
	path := filepath.Join(c.FormValue("path"), c.FormValue("relativePath"))

	// handle the request
	chunkNumber, err := strconv.ParseInt(c.FormValue("chunkNumber"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	chunkSize, err := strconv.ParseInt(c.FormValue("chunkSize"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	currentChunkSize, err := strconv.ParseInt(c.FormValue("currentChunkSize"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	totalChunks, err := strconv.ParseInt(c.FormValue("totalChunks"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	totalSize, err := strconv.ParseInt(c.FormValue("totalSize"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	identifier := c.FormValue("identifier")
	fileName := c.FormValue("filename")
	bin, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	s.lock.Lock()
	fileInfoTemp, ok := s.uploadStatus.Load(identifier)
	var fileInfo *FileInfo

	file, err := os.OpenFile(path+"/"+fileName+".tmp", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		s.lock.Unlock()
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !ok {
		// file, err := os.Create(path + "/" + fileName + ".tmp")

		if err != nil {
			s.lock.Unlock()
			return c.JSON(http.StatusInternalServerError, err)
		}

		// pre allocate file size
		fmt.Println("truncate", totalSize)
		if err != nil {
			s.lock.Unlock()
			return c.JSON(http.StatusInternalServerError, err)
		}

		// file info init
		fileInfo = &FileInfo{
			init:             true,
			uploaded:         make([]bool, totalChunks),
			uploadedChunkNum: 0,
		}
		s.uploadStatus.Store(identifier, fileInfo)
	} else {
		fileInfo = fileInfoTemp.(*FileInfo)
	}

	s.lock.Unlock()

	_, err = file.Seek((chunkNumber-1)*chunkSize, io.SeekStart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	src, err := bin.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer src.Close()

	buf := make([]byte, int(currentChunkSize))
	_, err = io.CopyBuffer(file, src, buf)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	s.lock.Lock()
	// handle file after write a chunk
	// handle single chunk upload twice
	if !fileInfo.uploaded[chunkNumber-1] {
		fileInfo.uploadedChunkNum++
		fileInfo.uploaded[chunkNumber-1] = true
	}

	// handle file after write all chunk
	if fileInfo.uploadedChunkNum == totalChunks {
		file.Close()
		os.Rename(path+"/"+fileName+".tmp", path+"/"+fileName)

		// remove upload status info after upload complete
		s.uploadStatus.Delete(identifier)
	}
	s.lock.Unlock()

	return c.NoContent(http.StatusOK)
}
