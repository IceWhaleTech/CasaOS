package v1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
)

func downloadReadFile(c *gin.Context) {
	//http下载地址 csv
	csvFileUrl := c.PostForm("file_name")
	res, err := http.Get(csvFileUrl)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer res.Body.Close()
	//读取csv
	reader := csv.NewReader(bufio.NewReader(res.Body))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			c.String(400, err.Error())
			return
		}
		//line 就是每一行的内容
		fmt.Println(line)
		//line[0] 就是第几列
		fmt.Println(line[0])
	}
}

func downloadWriteFile(c *gin.Context) {
	//写文件
	var filename = "./output1.csv"

	file, err := os.Create(filename) //创建文件
	if err != nil {
		c.String(400, err.Error())
		return
	}
	buf := bufio.NewWriter(file) //创建新的 Writer 对象
	buf.WriteString("test")
	buf.Flush()
	defer file.Close()

	//返回文件流
	c.File(filename)
}

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
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	//文件读取任务是将文件内容读取到内存中。
	info, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_READ_ERROR,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_READ_ERROR),
			Data:    err.Error(),
		})
		return
	}
	result := string(info)

	//返回结果
	c.JSON(http.StatusOK, model.Result{
		Success: oasis_err2.SUCCESS,
		Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
		Data:    result,
	})
}

func GetLocalFile(c *gin.Context) {
	path := c.Query("path")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	c.File(path)
	return
}

// @Summary download
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "path of file"
// @Success 200 {string} string "ok"
// @Router /file/download [get]
func GetDownloadFile(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	//打开文件
	fileTmp, _ := os.Open(filePath)
	defer fileTmp.Close()
	//获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url2.PathEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(filePath)
}

// @Summary download
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "path of file"
// @Success 200 {string} string "ok"
// @Router /file/new/download [get]
func GetFileDownloadNew(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	//打开文件
	fileStat, _ := os.Stat(filePath)
	var AppFs = afero.NewOsFs()
	fileT, _ := AppFs.Open(filePath)
	//fileTmp, _ := os.Open(filePath)
	//defer fileTmp.Close()
	//获取文件的名称
	//fileName := path.Base(filePath)

	//c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url2.PathEscape(fileName))
	//在线
	//c.Header("Content-Disposition", "inline")
	// extraHeaders := map[string]string{
	// 	"Content-Disposition": `attachment; filename="` + url2.PathEscape(fileName) + `"`,
	// }

	//c.Header("Cache-Control", "private")
	//c.Header("Content-Type", "application/octet-stream")

	http.ServeContent(c.Writer, c.Request, fileStat.Name(), fileStat.ModTime(), fileT)
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
	path := c.DefaultQuery("path", "")
	info := service.MyService.ZiMa().GetDirPath(path)
	if path == "/DATA/AppData" {
		list := service.MyService.Docker().DockerContainerList()
		apps := make(map[string]string, len(list))
		for _, v := range list {
			apps[strings.ReplaceAll(v.Names[0], "/", "")] = strings.ReplaceAll(v.Names[0], "/", "")
		}
		for i := 0; i < len(info); i++ {
			if v, ok := apps[info[i].Name]; ok {
				info[i].Label = v
				info[i].Type = "application"
			}
		}
	} else if path == "/DATA" {
		disk := make(map[string]string)
		lsblk := service.MyService.Disk().LSBLK(true)
		for _, v := range lsblk {
			if len(v.Children) > 0 {
				t := v.Tran
				for _, c := range v.Children {
					if len(c.Children) > 0 {
						for _, gc := range c.Children {
							if len(gc.MountPoint) > 0 {
								disk[gc.MountPoint] = t
							}
						}
					}
					if len(c.MountPoint) > 0 {
						disk[c.MountPoint] = t
					}
				}

			}
		}
		for i := 0; i < len(info); i++ {
			if v, ok := disk[info[i].Path]; ok {
				info[i].Type = v
			}
		}
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary rename file or dir
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param oldpath formData string true "path of old"
// @Param newpath formData string true "path of new"
// @Success 200 {string} string "ok"
// @Router /file/rename [put]
func RenamePath(c *gin.Context) {
	op := c.PostForm("oldpath")
	np := c.PostForm("newpath")
	if len(op) == 0 || len(np) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	success, err := service.MyService.ZiMa().RenameFile(op, np)
	c.JSON(http.StatusOK, model.Result{Success: success, Message: oasis_err2.GetMsg(success), Data: err})
}

// @Summary create folder
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string true "path of folder"
// @Success 200 {string} string "ok"
// @Router /file/mkdir [post]
func MkdirAll(c *gin.Context) {
	path := c.PostForm("path")
	var code int
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	code, _ = service.MyService.ZiMa().MkdirAll(path)
	c.JSON(http.StatusOK, model.Result{Success: code, Message: oasis_err2.GetMsg(code)})
}

// @Summary create file
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/create [post]
func PostCreateFile(c *gin.Context) {
	path := c.PostForm("path")
	var code int
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	code, _ = service.MyService.ZiMa().CreateFile(path)
	c.JSON(http.StatusOK, model.Result{Success: code, Message: oasis_err2.GetMsg(code)})
}

// @Summary upload file
// @Produce  application/json
// @Accept  multipart/form-data
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
	tempDir := config.AppInfo.RootPath + "/temp/" + hash + strconv.Itoa(totalChunks) + "/"
	if fileName != relative {
		dirPath = strings.TrimSuffix(relative, fileName)
		tempDir += dirPath
		file.MkDir(path + "/" + dirPath)
	}
	tempDir += chunkNumber
	if !file.CheckNotExist(tempDir) {
		c.JSON(200, model.Result{Success: 200, Message: oasis_err2.GetMsg(oasis_err2.FILE_ALREADY_EXISTS)})
		return
	}

	c.JSON(204, model.Result{Success: 204, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
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
		c.JSON(oasis_err2.INVALID_PARAMS, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	tempDir := config.AppInfo.RootPath + "/temp/" + hash + strconv.Itoa(totalChunks) + "/"

	if fileName != relative {
		dirPath = strings.TrimSuffix(relative, fileName)
		tempDir += dirPath
		file.MkDir(path + "/" + dirPath)
	}

	path += "/" + relative

	if !file.CheckNotExist(tempDir + chunkNumber) {
		file.RMDir(tempDir + chunkNumber)
	}

	if totalChunks > 1 {
		file.IsNotExistMkDir(tempDir)

		out, _ := os.OpenFile(tempDir+chunkNumber, os.O_WRONLY|os.O_CREATE, 0644)
		defer out.Close()
		_, err := io.Copy(out, f)
		if err != nil {
			c.JSON(oasis_err2.ERROR, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
	} else {
		out, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		defer out.Close()
		_, err := io.Copy(out, f)
		if err != nil {
			c.JSON(oasis_err2.ERROR, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
		return
	}
	fileNum, err := ioutil.ReadDir(tempDir)
	if err != nil {
		c.JSON(oasis_err2.ERROR, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	if totalChunks == len(fileNum) {
		file.SpliceFiles(tempDir, path, totalChunks, 1)
		file.RMDir(tempDir)
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary copy or move file
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param from formData string true "from path"
// @Param to formData string true "to path"
// @Param type formData string true "action" Enums(move,copy)
// @Success 200 {string} string "ok"
// @Router /file/operate [post]
func PostOperateFileOrDir(c *gin.Context) {
	from := c.PostForm("from")
	to := c.PostForm("to")
	t := c.PostForm("type")
	if len(from) == 0 || len(t) == 0 || len(to) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if t == "move" {
		lastPath := from[strings.LastIndex(from, "/")+1:]
		if !file.CheckNotExist(to + "/" + lastPath) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.FILE_OR_DIR_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.FILE_ALREADY_EXISTS)})
			return
		}
		err := os.Rename(from, to+"/"+lastPath)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
	} else if t == "copy" {
		err := file.CopyDir(from, to)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary delete file
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "path"
// @Success 200 {string} string "ok"
// @Router /file/delete [delete]
func DeleteFile(c *gin.Context) {
	path := c.Query("path")
	//err := os.Remove(path)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.FILE_DELETE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.FILE_DELETE_ERROR), Data: err})
		return
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary update file
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string true "path"
// @Param content formData string true "content"
// @Success 200 {string} string "ok"
// @Router /file/update [put]
func PutFileContent(c *gin.Context) {
	path := c.PostForm("path")
	content := c.PostForm("content")
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.FILE_ALREADY_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.FILE_ALREADY_EXISTS)})
		return
	}
	//err := os.Remove(path)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.FILE_DELETE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.FILE_DELETE_ERROR), Data: err})
		return
	}
	err = file.CreateFileAndWriteContent(path, content)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err})
		return
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
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
		c.JSON(http.StatusInternalServerError, model.Result{Success: oasis_err2.FILE_ALREADY_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.FILE_ALREADY_EXISTS)})
		return
	}
	if t == "thumbnail" {
		f, err := file.GetImage(path, 100, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
		c.Writer.WriteString(string(f))
		return
	}
	f, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	c.Writer.WriteString(string(data))
}
