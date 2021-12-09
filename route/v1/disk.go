package v1

import (
	"net/http"

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
	//格式化磁盘
	service.MyService.Disk().FormatDisk(path, t)

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

// @Summary serial number
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "磁盘路径 例如/dev/sda"
// @Param  serial formData string true "serial"
// @Success 200 {string} string "ok"
// @Router /disk/addpart [post]
func AddPartition(c *gin.Context) {
	path := c.PostForm("path")
	serial := c.PostForm("serial")
	if len(path) == 0 || len(serial) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	service.MyService.Disk().AddPartition(path)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

func PostMountDisk(c *gin.Context) {
	// for example: path=/dev/sda1
	path := c.PostForm("path")
	//执行挂载目录
	service.MyService.Disk().MountDisk(path, "volume")
	//添加到数据库

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

func DeleteUmountDisk(c *gin.Context) {

	// for example: path=/dev/sda1
	path := c.PostForm("path")
	service.MyService.Disk().UmountPointAndRemoveDir(path)

	//删除数据库记录

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}
