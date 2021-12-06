package v1

import (
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
)

// @Summary 获取磁盘列表
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/list [get]
func GetPlugInDisk(c *gin.Context) {

	//ls := service.MyService.Disk().GetPlugInDisk()
	//fmt.Println(ls)
	//dd, _ := disk.Partitions(true)
	//fmt.Println(dd)
	//
	//dir, err := ioutil.ReadDir("/sys/block")
	//if err != nil {
	//	panic(err)
	//}
	//
	//files := make([]string, 0)
	//
	////fmt.Println(regexp.MatchString("sd[a-z]*[0-9]", "sda"))
	//
	//for _, f := range dir {
	//	if match, _ := regexp.MatchString("sd[a-z]", f.Name()); match {
	//		files = append(files, f.Name())
	//	}
	//}
	//fmt.Println(files)
	//filess := make([]string, 0)
	//for _, file := range files {
	//	dirs, _ := ioutil.ReadDir("/sys/block/" + file)
	//
	//	for _, f := range dirs {
	//		if match, _ := regexp.MatchString("sd[a-z]*[0-9]", f.Name()); match {
	//			filess = append(filess, f.Name())
	//		}
	//	}
	//}
	//fmt.Println(filess)
	//
	//for _, s := range filess {
	//	fmt.Println(disk.Usage("/dev/" + s))
	//}

	lst := service.MyService.Disk().LSBLK()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: lst})
}

// @Summary get disk list
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/lists [get]
func GetPlugInDisks(c *gin.Context) {

	list := service.MyService.Disk().LSBLK()
	var result []*disk.UsageStat
	for _, item := range list {
		result = append(result, service.MyService.Disk().GetDiskInfoByPath(item.Path))
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: result})
}

// @Summary disk detail
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Param  path query string true "for example /dev/sda"
// @Success 200 {string} string "ok"
// @Router /disk/info [get]
func GetDiskInfo(c *gin.Context) {
	path := c.Query("path")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
	}
	m := service.MyService.Disk().GetDiskInfo(path)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: m})
}

// @Summary format disk
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "磁盘路径 例如/dev/sda1"
// @Success 200 {string} string "ok"
// @Router /disk/format [post]
func FormatDisk(c *gin.Context) {
	path := c.PostForm("path")

	t := c.PostForm("type")

	if len(path) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
	}

	//删除挂载点
	service.MyService.Disk().UmountPointAndRemoveDir(path)

	//格式化磁盘
	service.MyService.Disk().FormatDisk(path, t)

	//重新挂载

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary 获取支持的格式
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/type [get]
func FormatDiskType(c *gin.Context) {
	var strArr = [4]string{"fat32", "ntfs", "ext4", "exfat"}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: strArr})

}

// @Summary 删除分区
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "磁盘路径 例如/dev/sda1"
// @Success 200 {string} string "ok"
// @Router /disk/delpart [delete]
func RemovePartition(c *gin.Context) {
	path := c.PostForm("path")

	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
	}
	var p = path[:len(path)-1]
	var n = path[len(path)-1:]
	service.MyService.Disk().DelPartition(p, n)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary 添加分区
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "磁盘路径 例如/dev/sda"
// @Param  size formData string true "需要分区容量大小(MB)"
// @Param  num formData string true "磁盘符号"
// @Success 200 {string} string "ok"
// @Router /disk/addpart [post]
func AddPartition(c *gin.Context) {
	path := c.PostForm("path")
	size, _ := strconv.Atoi(c.DefaultPostForm("size", "0"))
	num := c.DefaultPostForm("num", "9")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
	}

	//size*1024*1024/512
	service.MyService.Disk().AddPartition(path, num, uint64(size*1024*2))
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}
