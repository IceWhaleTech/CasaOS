package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
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

func (s *FileUploadService) TestChunk(
	c echo.Context,
	identifier string,
	chunkNumber int64,
) error {
	fileInfoTemp, ok := s.uploadStatus.Load(identifier)

	if !ok {
		return fmt.Errorf("file not found")
	}

	fileInfo := fileInfoTemp.(*FileInfo)

	if !fileInfo.init {
		return fmt.Errorf("file not init")
	}

	// return StatusNoContent instead of 404
	// the is require by frontend
	if !fileInfo.uploaded[chunkNumber-1] {
		return fmt.Errorf("file not found")
	}

	return nil
}

func (s *FileUploadService) UploadFile(
	c echo.Context,
	path string,
	chunkNumber int64,
	chunkSize int64,
	currentChunkSize int64,
	totalChunks int64,
	totalSize int64,
	identifier string,
	fileName string,
	bin *multipart.FileHeader,
) error {
	s.lock.Lock()
	fileInfoTemp, ok := s.uploadStatus.Load(identifier)
	var fileInfo *FileInfo

	file, err := os.OpenFile(path+"/"+fileName+".tmp", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		s.lock.Unlock()
		return err
	}

	if !ok {

		if err != nil {
			s.lock.Unlock()
			return err
		}

		// pre allocate file size
		fmt.Println("truncate", totalSize)
		if err != nil {
			s.lock.Unlock()
			return err
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
		return err
	}

	src, err := bin.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := make([]byte, int(currentChunkSize))
	_, err = io.CopyBuffer(file, src, buf)

	if err != nil {
		fmt.Println(err)
		return err
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

	return nil
}
