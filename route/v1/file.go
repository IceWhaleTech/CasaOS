package v1

import (
	"encoding/json"
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
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/h2non/filetype"
)

type ListReq struct {
	model.PageReq
	Path string `json:"path" form:"path"`
	//Refresh bool   `json:"refresh"`
}
type ObjResp struct {
	Name       string                 `json:"name"`
	Size       int64                  `json:"size"`
	IsDir      bool                   `json:"is_dir"`
	Modified   time.Time              `json:"modified"`
	Sign       string                 `json:"sign"`
	Thumb      string                 `json:"thumb"`
	Type       int                    `json:"type"`
	Path       string                 `json:"path"`
	Date       time.Time              `json:"date"`
	Extensions map[string]interface{} `json:"extensions"`
}
type FsListResp struct {
	Content  []ObjResp `json:"content"`
	Total    int64     `json:"total"`
	Readme   string    `json:"readme,omitempty"`
	Write    bool      `json:"write,omitempty"`
	Provider string    `json:"provider,omitempty"`
	Index    int       `json:"index"`
	Size     int       `json:"size"`
}

var (
	// 升级成 WebSocket 协议
	upgraderFile = websocket.Upgrader{
		// 允许CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn *websocket.Conn
	err  error
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

	fi, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// We only have to pass the file header = first 261 bytes
	buffer := make([]byte, 261)

	_, _ = fi.Read(buffer)

	kind, _ := filetype.Match(buffer)
	if kind != filetype.Unknown {
		c.Header("Content-Type", kind.MIME.Value)
	}
	node, err := os.Stat(filePath)
	// Set the Last-Modified header to the timestamp
	c.Header("Last-Modified", node.ModTime().UTC().Format(http.TimeFormat))

	knownSize := node.Size() >= 0
	if knownSize {
		c.Header("Content-Length", strconv.FormatInt(node.Size(), 10))
	}
	http.ServeContent(c.Writer, c.Request, fileName, node.ModTime(), fi)
	//http.ServeFile(c.Writer, c.Request, filePath)
	defer fi.Close()
	return
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
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
		return
	}
	req.Validate()
	info, err := service.MyService.System().GetDirPath(req.Path)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
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
			ex["mounted"] = false
			info[i].Extensions = ex
		}
	}
	if strings.HasPrefix(req.Path, "/mnt") || strings.HasPrefix(req.Path, "/media") {
		for i := (req.Index - 1) * req.Size; i < forEnd; i++ {
			ex := info[i].Extensions
			if ex == nil {
				ex = make(map[string]interface{})
			}
			mounted := service.IsMounted(info[i].Path)
			ex["mounted"] = mounted
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
			t.Date = info[i].Date
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
	mounted := service.IsMounted(op)
	if mounted {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.MOUNTED_DIRECTIORIES, Message: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES), Data: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES)})
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
	if file.Exists(path + "/" + relative) {
		c.JSON(http.StatusConflict, model.Result{Success: http.StatusConflict, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
		return
	}
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
			go func() {
				time.Sleep(11 * time.Second)
				if err := file.RMDir(tempDir); err != nil {
					logger.Error("error when trying to remove `"+tempDir+"`", zap.Error(err))
				}
			}()
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

func PostFileOctet(c *gin.Context) {

	content_length := c.Request.ContentLength
	if content_length <= 0 || content_length > 1024*1024*1024*2*1024 {
		log.Printf("content_length error\n")
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "content_length error"})
		return
	}
	content_type_, has_key := c.Request.Header["Content-Type"]
	if !has_key {
		log.Printf("Content-Type error\n")
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Content-Type error"})
		return
	}
	if len(content_type_) != 1 {
		log.Printf("Content-Type count error\n")
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Content-Type count error"})
		return
	}
	content_type := content_type_[0]
	const BOUNDARY string = "; boundary="
	loc := strings.Index(content_type, BOUNDARY)
	if loc == -1 {
		log.Printf("Content-Type error, no boundary\n")
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Content-Type error, no boundary"})
		return
	}
	boundary := []byte(content_type[(loc + len(BOUNDARY)):])
	log.Printf("[%s]\n\n", boundary)
	read_data := make([]byte, 1024*24)
	var read_total int = 0
	for {
		file_header, file_data, err := file.ParseFromHead(read_data, read_total, append(boundary, []byte("\r\n")...), c.Request.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		log.Printf("file :%s\n", file_header)
		//
		//os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o644)
		f, err := os.OpenFile(file_header["path"]+"/"+file_header["filename"], os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			log.Printf("create file fail:%v\n", err)
			return
		}
		f.Write(file_data)
		file_data = nil

		temp_data, reach_end, err := file.ReadToBoundary(boundary, c.Request.Body, f)
		f.Close()
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		if reach_end {
			break
		} else {
			copy(read_data[0:], temp_data)
			read_total = len(temp_data)
			continue
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
		if list.Type == "move" {
			mounted := service.IsMounted(list.Item[i].From)
			if mounted {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.MOUNTED_DIRECTIORIES, Message: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES), Data: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES)})
				return
			}
		}
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
		mounted := service.IsMounted(v)
		if mounted {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.MOUNTED_DIRECTIORIES, Message: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES), Data: common_err.GetMsg(common_err.MOUNTED_DIRECTIORIES)})
			return
		}
	}

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
	f, err := os.Stat(fi.FilePath)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
		return
	}
	fm := f.Mode()
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.FILE_DELETE_ERROR, Message: common_err.GetMsg(common_err.FILE_DELETE_ERROR), Data: err})
		return
	}
	os.OpenFile(fi.FilePath, os.O_CREATE, fm)
	err = file.WriteToFullPath([]byte(fi.FileContent), fi.FilePath, fm)
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

func GetFileCount(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	path := json["path"]
	list, err := ioutil.ReadDir(path)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: len(list)})
}

type CenterHandler struct {
	// 广播通道，有数据则循环每个用户广播出去
	broadcast chan []byte
	// 注册通道，有用户进来 则推到用户集合map中
	register chan *Client
	// 注销通道，有用户关闭连接 则将该用户剔出集合map中
	unregister chan *Client
	// 用户集合，每个用户本身也在跑两个协程，监听用户的读、写的状态
	clients map[string]*Client
}

type Client struct {
	handler *CenterHandler
	conn    *websocket.Conn
	// 每个用户自己的循环跑起来的状态监控
	send         chan []byte
	ID           string       `json:"id"`
	IP           string       `json:"ip"`
	Name         service.Name `json:"name"`
	RtcSupported bool         `json:"rtcSupported"`
	TimerId      int          `json:"timerId"`
	LastBeat     time.Time    `json:"lastBeat"`
}

type PeerModel struct {
	ID           string       `json:"id"`
	Name         service.Name `json:"name"`
	RtcSupported bool         `json:"rtcSupported"`
}

func ConnectWebSocket(c *gin.Context) {
	peerId := c.Query("peer")
	writer := c.Writer
	request := c.Request
	key := uuid.NewV4().String()
	//peerModel := service.MyService.Peer().GetPeerByUserAgent(c.Request.UserAgent())
	peerModel := model2.PeerDriveDBModel{}
	name := service.GetName(request)
	if conn, err = upgraderFile.Upgrade(writer, request, writer.Header()); err != nil {
		log.Println(err)
		return
	}
	client := &Client{handler: &handler, conn: conn, send: make(chan []byte, 256), ID: service.GetPeerId(request, key), IP: service.GetIP(request), Name: name, RtcSupported: true, TimerId: 0, LastBeat: time.Now()}
	if peerId != "" || len(peerModel.ID) > 0 {
		if len(peerModel.ID) == 0 {
			peerModel = service.MyService.Peer().GetPeerByID(peerId)
		}
		if len(peerModel.ID) > 0 {
			key = peerId
			client.ID = peerModel.ID
			client.Name = service.GetNameByDB(peerModel)
		}
	}
	var list = service.MyService.Peer().GetPeers()
	if len(peerModel.ID) == 0 {
		peerModel.ID = key
		peerModel.DisplayName = name.DisplayName
		peerModel.DeviceName = name.DeviceName
		peerModel.Model = name.Model
		peerModel.OS = name.OS
		peerModel.Browser = name.Browser
		peerModel.UserAgent = c.Request.UserAgent()
		peerModel.IP = client.IP
		service.MyService.Peer().CreatePeer(&peerModel)
		list = append(list, peerModel)
	}

	cookie := http.Cookie{
		Name:  "peerid",
		Value: key,
		Path:  "/",
	}
	http.SetCookie(writer, &cookie)
	if len(list) > 10 {
		kickoutList := []Client{}
		count := len(list) - 10
		for i := len(list) - 1; count > 0 && i > -1; i-- {
			if _, ok := handler.clients[list[i].ID]; !ok {
				count--
				kickoutList = append(kickoutList, Client{ID: list[i].ID, Name: service.GetNameByDB(list[i]), IP: list[i].IP})
				service.MyService.Peer().DeletePeer(list[i].ID)
			}
		}
		// if len(kickoutList) > 0 {
		// 	other := make(map[string]interface{})
		// 	other["type"] = "kickout"
		// 	other["peers"] = kickoutList
		// 	otherBy, err := json.Marshal(other)
		// 	fmt.Println(err)
		// 	client.handler.broadcast <- otherBy
		// }
	}
	list = service.MyService.Peer().GetPeers()
	if len(list) > 10 {
		fmt.Println("解决完后依然有溢出", list)
	}
	currentPeer := PeerModel{ID: client.ID, Name: client.Name, RtcSupported: client.RtcSupported}
	pmsg := make(map[string]interface{})
	pmsg["type"] = "peer-joined"
	pmsg["peer"] = currentPeer
	pby, err := json.Marshal(pmsg)
	fmt.Println(err)
	for _, v := range handler.clients {
		v.send <- pby
	}
	//client.handler.broadcast <- pby
	clients := []PeerModel{}
	for _, v := range client.handler.clients {
		if _, ok := handler.clients[v.ID]; ok {
			clients = append(clients, PeerModel{ID: v.ID, Name: v.Name, RtcSupported: v.RtcSupported})
		}
	}

	other := make(map[string]interface{})
	other["type"] = "peers"
	other["peers"] = clients
	otherBy, err := json.Marshal(other)
	fmt.Println(err)
	client.send <- otherBy

	// 推给监控中心注册到用户集合中
	handler.register <- client

	client.send <- []byte(`{"type":"ping"}`)

	data := make(map[string]string)
	data["displayName"] = client.Name.DisplayName
	data["deviceName"] = client.Name.DeviceName
	data["id"] = client.ID
	msg := make(map[string]interface{})
	msg["type"] = "display-name"
	msg["message"] = data
	by, _ := json.Marshal(msg)
	client.send <- by

	// 每个 client 都挂起 2 个新的协程，监控读、写状态
	go client.writePump()
	go client.readPump()
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

var handler = CenterHandler{broadcast: make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[string]*Client)}

func init() {
	// 起个协程跑起来，监听注册、注销、消息 3 个 channel
	go handler.monitoring()

	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定义定时器调用的任务函数

	task := func() {
		handler.broadcast <- []byte(`{"type":"ping"}`)
	}
	//定时任务
	spec := "*/30 * * * * ?" //cron表达式，每五秒一次
	// 添加定时任务,
	crontab.AddFunc(spec, task)
	// 启动定时器
	crontab.Start()
}
func (c *Client) writePump() {
	defer func() {
		c.handler.unregister <- c

		c.conn.Close()
	}()
	for {
		// 广播推过来的新消息，马上通过websocket推给自己
		message, _ := <-c.send
		fmt.Println("推送消息", string(message), "1")
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 读，监听客户端是否有推送内容过来服务端
func (c *Client) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 循环监听是否该用户是否要发言
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 异常关闭的处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			c.handler.broadcast <- []byte(`{"type":"peer-left","peerId":"` + c.ID + `"}`)
			break
		}
		// 要的话，推给广播中心，广播中心再推给每个用户

		t := gjson.GetBytes(message, "type")
		if t.String() == "disconnect" {
			c.handler.unregister <- c
			c.conn.Close()
			// clients := []Client{}
			// list := service.MyService.Peer().GetPeers()
			// for _, v := range list {
			// 	if _, ok := handler.clients[v.ID]; ok {
			// 		clients = append(clients, *handler.clients[v.ID])
			// 	} else {
			// 		clients = append(clients, Client{ID: v.ID, Name: service.GetNameByDB(v), IP: v.IP, Offline: true})
			// 	}
			// }
			// other := make(map[string]interface{})
			// other["type"] = "peers"
			// other["peers"] = clients
			// otherBy, err := json.Marshal(other)
			// fmt.Println(err)
			c.handler.broadcast <- []byte(`{"type":"peer-left","peerId":"` + c.ID + `"}`)
			//c.handler.broadcast <- otherBy
			break
		} else if t.String() == "pong" {
			c.LastBeat = time.Now()
			continue
		}
		to := gjson.GetBytes(message, "to")

		if len(to.String()) > 0 {
			toC := c.handler.clients[to.String()]
			if toC == nil {
				continue
			}
			data := map[string]interface{}{}
			json.Unmarshal(message, &data)
			data["sender"] = c.ID
			delete(data, "to")
			message, err = json.Marshal(data)
			toC.send <- message
			continue
		}

		c.handler.broadcast <- message
	}
}
func (ch *CenterHandler) monitoring() {
	for {
		select {
		// 注册，新用户连接过来会推进注册通道，这里接收推进来的用户指针
		case client := <-ch.register:
			ch.clients[client.ID] = client
			// 注销，关闭连接或连接异常会将用户推出群聊
		case client := <-ch.unregister:
			delete(ch.clients, client.ID)
			// 消息，监听到有新消息到来
		case message := <-ch.broadcast:
			println("消息来了，message：" + string(message))
			// 推送给每个用户的通道，每个用户都有跑协程起了writePump的监听
			for _, client := range ch.clients {
				client.send <- message
			}
		}
	}
}
func GetPeers(c *gin.Context) {
	peers := service.MyService.Peer().GetPeers()
	for i := 0; i < len(peers); i++ {
		if _, ok := handler.clients[peers[i].ID]; ok {
			peers[i].Online = true
		}
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: peers})
}
