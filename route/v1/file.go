package v1

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/internal/conf"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/model"

	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"

	"github.com/IceWhaleTech/CasaOS/internal/sign"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// @Summary 读取文件
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "路径"
// @Success 200 {string} string "ok"
// @Router /file/read [get]
func GetFilerContent(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{
			Success: common_err.INVALID_PARAMS,
			Message: common_err.GetMsg(common_err.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(common_err.SERVICE_ERROR, model.Result{
			Success: common_err.FILE_DOES_NOT_EXIST,
			Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST),
		})
		return
	}
	// 文件读取任务是将文件内容读取到内存中。
	info, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{
			Success: common_err.FILE_READ_ERROR,
			Message: common_err.GetMsg(common_err.FILE_READ_ERROR),
			Data:    err.Error(),
		})
		return
	}
	result := string(info)

	c.JSON(common_err.SUCCESS, model.Result{
		Success: common_err.SUCCESS,
		Message: common_err.GetMsg(common_err.SUCCESS),
		Data:    result,
	})
}

func GetLocalFile(c *gin.Context) {
	path := c.Query("path")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: common_err.INVALID_PARAMS,
			Message: common_err.GetMsg(common_err.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{
			Success: common_err.FILE_DOES_NOT_EXIST,
			Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST),
		})
		return
	}
	c.File(path)
}

// @Summary download
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param format query string false "Compression format" Enums(zip,tar,targz)
// @Param files query string true "file list eg: filename1,filename2,filename3 "
// @Success 200 {string} string "ok"
// @Router /file/download [get]
func GetDownloadFile(c *gin.Context) {
	t := c.Query("format")

	files := c.Query("files")

	if len(files) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{
			Success: common_err.INVALID_PARAMS,
			Message: common_err.GetMsg(common_err.INVALID_PARAMS),
		})
		return
	}
	list := strings.Split(files, ",")
	for _, v := range list {
		if !file.Exists(v) {
			c.JSON(common_err.SERVICE_ERROR, model.Result{
				Success: common_err.FILE_DOES_NOT_EXIST,
				Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST),
			})
			return
		}
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	// handles only single files not folders and multiple files
	if len(list) == 1 {

		filePath := list[0]
		info, err := os.Stat(filePath)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{
				Success: common_err.FILE_DOES_NOT_EXIST,
				Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST),
			})
			return
		}
		if !info.IsDir() {

			// 打开文件
			fileTmp, _ := os.Open(filePath)
			defer fileTmp.Close()

			// 获取文件的名称
			fileName := path.Base(filePath)
			c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url2.PathEscape(fileName))
			c.File(filePath)
			return
		}
	}

	extension, ar, err := file.GetCompressionAlgorithm(t)
	if err != nil {
		c.JSON(common_err.CLIENT_ERROR, model.Result{
			Success: common_err.INVALID_PARAMS,
			Message: common_err.GetMsg(common_err.INVALID_PARAMS),
		})
		return
	}

	err = ar.Create(c.Writer)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{
			Success: common_err.SERVICE_ERROR,
			Message: common_err.GetMsg(common_err.SERVICE_ERROR),
			Data:    err.Error(),
		})
		return
	}
	defer ar.Close()
	commonDir := file.CommonPrefix(filepath.Separator, list...)

	currentPath := filepath.Base(commonDir)

	name := "_" + currentPath
	name += extension
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(name))
	for _, fname := range list {
		err = file.AddFile(ar, fname, commonDir)
		if err != nil {
			log.Printf("Failed to archive %s: %v", fname, err)
		}
	}
}

func GetDownloadSingleFile(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{
			Success: common_err.INVALID_PARAMS,
			Message: common_err.GetMsg(common_err.INVALID_PARAMS),
		})
		return
	}
	fileName := path.Base(filePath)
	// c.Header("Content-Disposition", "inline")
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url2.PathEscape(fileName))

	storage, _ := service.MyService.FsService().GetStorage(filePath)
	if storage != nil {
		if shouldProxy(storage, fileName) {
			Proxy(c)
			return
		} else {
			link, _, err := service.MyService.FsService().Link(c, filePath, model.LinkArgs{
				IP:     c.ClientIP(),
				Header: c.Request.Header,
				Type:   c.Query("type"),
			})
			if err != nil {
				c.JSON(common_err.SERVICE_ERROR, model.Result{
					Success: common_err.SERVICE_ERROR,
					Message: common_err.GetMsg(common_err.SERVICE_ERROR),
					Data:    err.Error(),
				})
				return

			}
			c.Header("Referrer-Policy", "no-referrer")
			c.Header("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
			c.Redirect(302, link.URL)
			return
		}
	}

	fileTmp, err := os.Open(filePath)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{
			Success: common_err.FILE_DOES_NOT_EXIST,
			Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST),
		})
		return
	}
	defer fileTmp.Close()

	c.File(filePath)
}

// @Summary 获取目录列表
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/dirpath [get]
func DirPath(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
		return
	}
	req.Validate()
	info := service.MyService.System().GetDirPath(req.Path)
	shares := service.MyService.Shares().GetSharesList()
	sharesMap := make(map[string]string)
	for _, v := range shares {
		sharesMap[v.Path] = fmt.Sprint(v.ID)
	}
	// if len(info) <= (req.Page-1)*req.Size {
	// 	c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "page out of range"})
	// 	return
	// }
	forEnd := req.Index * req.Size
	if forEnd > len(info) {
		forEnd = len(info)
	}
	for i := (req.Index - 1) * req.Size; i < forEnd; i++ {
		if v, ok := sharesMap[info[i].Path]; ok {
			ex := make(map[string]interface{})
			shareEx := make(map[string]string)
			shareEx["shared"] = "true"
			shareEx["id"] = v
			ex["share"] = shareEx
			info[i].Extensions = ex
		}
	}
	// Hide the files or folders in operation
	fileQueue := make(map[string]string)
	if len(service.OpStrArr) > 0 {
		for _, v := range service.OpStrArr {
			v, ok := service.FileQueue.Load(v)
			if !ok {
				continue
			}
			vt := v.(model.FileOperate)
			for _, i := range vt.Item {
				lastPath := i.From[strings.LastIndex(i.From, "/")+1:]
				fileQueue[vt.To+"/"+lastPath] = i.From
			}
		}
	}

	pathList := []ObjResp{}
	for i := (req.Index - 1) * req.Size; i < forEnd; i++ {
		if info[i].Name == ".temp" && info[i].IsDir {
			continue
		}
		if _, ok := fileQueue[info[i].Path]; !ok {
			t := ObjResp{}
			t.IsDir = info[i].IsDir
			t.Name = info[i].Name
			t.Modified = info[i].Date
			t.Size = info[i].Size
			t.Path = info[i].Path
			t.Extensions = info[i].Extensions
			pathList = append(pathList, t)

		}
	}
	flist := FsListResp{
		Content: pathList,
		Total:   int64(len(info)),
		// Readme:   "",
		// Write:    true,
		// Provider: "local",
		Index: req.Index,
		Size:  req.Size,
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: flist})
}

// @Summary rename file or dir
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param oldpath body string true "path of old"
// @Param newpath body string true "path of new"
// @Success 200 {string} string "ok"
// @Router /file/rename [put]
func RenamePath(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	op := json["old_path"]
	np := json["new_path"]
	if len(op) == 0 || len(np) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	success, err := service.MyService.System().RenameFile(op, np)
	c.JSON(common_err.SUCCESS, model.Result{Success: success, Message: common_err.GetMsg(success), Data: err})
}

// @Summary create folder
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path body string true "path of folder"
// @Success 200 {string} string "ok"
// @Router /file/mkdir [post]
func MkdirAll(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	path := json["path"]
	var code int
	if len(path) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	// decodedPath, err := url.QueryUnescape(path)
	// if err != nil {
	// 	c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	// 	return
	// }
	code, _ = service.MyService.System().MkdirAll(path)
	c.JSON(common_err.SUCCESS, model.Result{Success: code, Message: common_err.GetMsg(code)})
}

// @Summary create file
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path body string true "path of folder (path need to url encode)"
// @Success 200 {string} string "ok"
// @Router /file/create [post]
func PostCreateFile(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	path := json["path"]
	var code int
	if len(path) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	// decodedPath, err := url.QueryUnescape(path)
	// if err != nil {
	// 	c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	// 	return
	// }
	code, _ = service.MyService.System().CreateFile(path)
	c.JSON(common_err.SUCCESS, model.Result{Success: code, Message: common_err.GetMsg(code)})
}

// @Summary upload file
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "file path"
// @Param file formData file true "file"
// @Success 200 {string} string "ok"
// @Router /file/upload [get]
func GetFileUpload(c *gin.Context) {
	relative := c.Query("relativePath")
	fileName := c.Query("filename")
	chunkNumber := c.Query("chunkNumber")
	totalChunks, _ := strconv.Atoi(c.DefaultQuery("totalChunks", "0"))
	path := c.Query("path")
	dirPath := ""
	hash := file.GetHashByContent([]byte(fileName))
	tempDir := filepath.Join(path, ".temp", hash+strconv.Itoa(totalChunks)) + "/"
	if fileName != relative {
		dirPath = strings.TrimSuffix(relative, fileName)
		tempDir += dirPath
		file.MkDir(path + "/" + dirPath)
	}
	tempDir += chunkNumber
	if !file.CheckNotExist(tempDir) {
		c.JSON(200, model.Result{Success: 200, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
		return
	}

	c.JSON(204, model.Result{Success: 204, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary upload file
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "file path"
// @Param file formData file true "file"
// @Success 200 {string} string "ok"
// @Router /file/upload [post]
func PostFileUpload(c *gin.Context) {
	f, _, _ := c.Request.FormFile("file")
	relative := c.PostForm("relativePath")
	fileName := c.PostForm("filename")
	totalChunks, _ := strconv.Atoi(c.DefaultPostForm("totalChunks", "0"))
	chunkNumber := c.PostForm("chunkNumber")
	dirPath := ""
	path := c.PostForm("path")

	hash := file.GetHashByContent([]byte(fileName))

	if len(path) == 0 {
		logger.Error("path should not be empty")
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	tempDir := filepath.Join(path, ".temp", hash+strconv.Itoa(totalChunks)) + "/"

	if fileName != relative {
		dirPath = strings.TrimSuffix(relative, fileName)
		tempDir += dirPath
		if err := file.MkDir(path + "/" + dirPath); err != nil {
			logger.Error("error when trying to create `"+path+"/"+dirPath+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}
	}

	path += "/" + relative

	if !file.CheckNotExist(tempDir + chunkNumber) {
		if err := file.RMDir(tempDir + chunkNumber); err != nil {
			logger.Error("error when trying to remove existing `"+tempDir+chunkNumber+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}
	}

	if totalChunks > 1 {
		if err := file.IsNotExistMkDir(tempDir); err != nil {
			logger.Error("error when trying to create `"+tempDir+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}

		out, err := os.OpenFile(tempDir+chunkNumber, os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			logger.Error("error when trying to open `"+tempDir+chunkNumber+"` for creation", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}

		defer out.Close()

		if _, err := io.Copy(out, f); err != nil { // recommend to use https://github.com/iceber/iouring-go for faster copy
			logger.Error("error when trying to write to `"+tempDir+chunkNumber+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}

		fileNum, err := ioutil.ReadDir(tempDir)
		if err != nil {
			logger.Error("error when trying to read number of files under `"+tempDir+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}

		if totalChunks == len(fileNum) {
			if err := file.SpliceFiles(tempDir, path, totalChunks, 1); err != nil {
				logger.Error("error when trying to splice files under `"+tempDir+"`", zap.Error(err))
				c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
				return
			}

			if err := file.RMDir(tempDir); err != nil {
				logger.Error("error when trying to remove `"+tempDir+"`", zap.Error(err))
				c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
				return
			}
		}
	} else {
		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			logger.Error("error when trying to open `"+path+"` for creation", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: err.Error()})
			return
		}

		defer out.Close()

		if _, err := io.Copy(out, f); err != nil { // recommend to use https://github.com/iceber/iouring-go for faster copy
			logger.Error("error when trying to write to `"+path+"`", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary copy or move file
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param body body model.FileOperate true "type:move,copy"
// @Success 200 {string} string "ok"
// @Router /file/operate [post]
func PostOperateFileOrDir(c *gin.Context) {
	list := model.FileOperate{}
	c.ShouldBind(&list)

	if len(list.Item) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if list.To == list.Item[0].From[:strings.LastIndex(list.Item[0].From, "/")] {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SOURCE_DES_SAME, Message: common_err.GetMsg(common_err.SOURCE_DES_SAME)})
		return
	}

	var total int64 = 0
	for i := 0; i < len(list.Item); i++ {

		size, err := file.GetFileOrDirSize(list.Item[i].From)
		if err != nil {
			continue
		}
		list.Item[i].Size = size
		total += size
	}

	list.TotalSize = total
	list.ProcessedSize = 0

	uid := uuid.NewV4().String()
	service.FileQueue.Store(uid, list)
	service.OpStrArr = append(service.OpStrArr, uid)
	if len(service.OpStrArr) == 1 {
		go service.ExecOpFile()
		go service.CheckFileStatus()

		go service.MyService.Notify().SendFileOperateNotify(false)

	}

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary delete file
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param body body string true "paths eg ["/a/b/c","/d/e/f"]"
// @Success 200 {string} string "ok"
// @Router /file/delete [delete]
func DeleteFile(c *gin.Context) {
	paths := []string{}
	c.ShouldBind(&paths)
	if len(paths) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	//	path := c.Query("path")

	//	paths := strings.Split(path, ",")

	for _, v := range paths {
		err := os.RemoveAll(v)
		if err != nil {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_DELETE_ERROR, Message: common_err.GetMsg(common_err.FILE_DELETE_ERROR), Data: err})
			return
		}
	}

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary update file
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path body string true "path"
// @Param content body string true "content"
// @Success 200 {string} string "ok"
// @Router /file/update [put]
func PutFileContent(c *gin.Context) {
	fi := model.FileUpdate{}
	c.ShouldBind(&fi)

	// path := c.PostForm("path")
	// content := c.PostForm("content")
	if !file.Exists(fi.FilePath) {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
		return
	}
	// err := os.Remove(path)
	err := os.RemoveAll(fi.FilePath)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_DELETE_ERROR, Message: common_err.GetMsg(common_err.FILE_DELETE_ERROR), Data: err})
		return
	}
	err = file.CreateFileAndWriteContent(fi.FilePath, fi.FileContent)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary image thumbnail/original image
// @Produce  application/json
// @Accept  application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "path"
// @Param type query string false "original,thumbnail" Enums(original,thumbnail)
// @Success 200 {string} string "ok"
// @Router /file/image [get]
func GetFileImage(c *gin.Context) {
	t := c.Query("type")
	path := c.Query("path")
	if !file.Exists(path) {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
		return
	}
	if t == "thumbnail" {
		f, err := file.GetImage(path, 100, 0)
		if err != nil {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
			return
		}
		c.Writer.WriteString(string(f))
		return
	}
	f, err := os.Open(path)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	c.Writer.WriteString(string(data))
}

func DeleteOperateFileOrDir(c *gin.Context) {
	id := c.Param("id")
	if id == "0" {
		service.FileQueue = sync.Map{}
		service.OpStrArr = []string{}
	} else {

		service.FileQueue.Delete(id)
		tempList := []string{}
		for _, v := range service.OpStrArr {
			if v != id {
				tempList = append(tempList, v)
			}
		}
		service.OpStrArr = tempList

	}

	go service.MyService.Notify().SendFileOperateNotify(true)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
func GetSize(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	path := json["path"]
	size, err := file.GetFileOrDirSize(path)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: size})
}
func Proxy(c *gin.Context) {
	rawPath := c.Query("path")
	filename := filepath.Base(rawPath)
	storage, err := service.MyService.FsService().GetStorage(rawPath)
	if err != nil {
		c.JSON(500, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	if canProxy(storage, filename) {
		downProxyUrl := storage.GetStorage().DownProxyUrl
		if downProxyUrl != "" {
			_, ok := c.GetQuery("d")
			if !ok {
				URL := fmt.Sprintf("%s%s?sign=%s",
					strings.Split(downProxyUrl, "\n")[0],
					utils.EncodePath(rawPath, true),
					sign.Sign(rawPath))
				c.Redirect(302, URL)
				return
			}
		}
		link, file, err := service.MyService.FsService().Link(c, rawPath, model.LinkArgs{
			Header: c.Request.Header,
			Type:   c.Query("type"),
		})
		if err != nil {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})

			return
		}
		err = CommonProxy(c.Writer, c.Request, link, file)
		if err != nil {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
			return
		}
	} else {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: "proxy not allowed"})
		return
	}
}

// TODO need optimize
// when should be proxy?
// 1. config.MustProxy()
// 2. storage.WebProxy
// 3. proxy_types
func shouldProxy(storage driver.Driver, filename string) bool {
	if storage.Config().MustProxy() || storage.GetStorage().WebProxy {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.ProxyTypes], utils.Ext(filename)) {
		return true
	}
	return false
}

// TODO need optimize
// when can be proxy?
// 1. text file
// 2. config.MustProxy()
// 3. storage.WebProxy
// 4. proxy_types
// solution: text_file + shouldProxy()
func canProxy(storage driver.Driver, filename string) bool {
	if storage.Config().MustProxy() || storage.GetStorage().WebProxy || storage.GetStorage().WebdavProxy() {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.ProxyTypes], utils.Ext(filename)) {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.TextTypes], utils.Ext(filename)) {
		return true
	}
	return false
}

var HttpClient = &http.Client{}

func CommonProxy(w http.ResponseWriter, r *http.Request, link *model.Link, file model.Obj) error {
	// read data with native
	var err error
	if link.Data != nil {
		defer func() {
			_ = link.Data.Close()
		}()
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, file.GetName(), url.QueryEscape(file.GetName())))
		w.Header().Set("Content-Length", strconv.FormatInt(file.GetSize(), 10))
		if link.Header != nil {
			// TODO clean header with blacklist or whitelist
			link.Header.Del("set-cookie")
			for h, val := range link.Header {
				w.Header()[h] = val
			}
		}
		if link.Status == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(link.Status)
		}
		_, err = io.Copy(w, link.Data)
		if err != nil {
			return err
		}
		return nil
	}
	// local file
	if link.FilePath != nil && *link.FilePath != "" {
		f, err := os.Open(*link.FilePath)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
		fileStat, err := os.Stat(*link.FilePath)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, file.GetName(), url.QueryEscape(file.GetName())))
		http.ServeContent(w, r, file.GetName(), fileStat.ModTime(), f)
		return nil
	} else {
		req, err := http.NewRequest(link.Method, link.URL, nil)
		if err != nil {
			return err
		}
		for h, val := range r.Header {
			if utils.SliceContains(conf.SlicesMap[conf.ProxyIgnoreHeaders], strings.ToLower(h)) {
				continue
			}
			req.Header[h] = val
		}
		for h, val := range link.Header {
			req.Header[h] = val
		}
		res, err := HttpClient.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			_ = res.Body.Close()
		}()
		logger.Info("proxy status", zap.Any("status", res.StatusCode))
		// TODO clean header with blacklist or whitelist
		res.Header.Del("set-cookie")
		for h, v := range res.Header {
			w.Header()[h] = v
		}
		w.WriteHeader(res.StatusCode)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, file.GetName(), url.QueryEscape(file.GetName())))
		if res.StatusCode >= 400 {
			all, _ := ioutil.ReadAll(res.Body)
			msg := string(all)
			logger.Info("msg", zap.Any("msg", msg))

			return errors.New(msg)
		}
		_, err = io.Copy(w, res.Body)
		if err != nil {
			return err
		}
		return nil
	}
}
