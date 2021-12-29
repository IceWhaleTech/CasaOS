package v1

import (
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
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

	list := service.MyService.Disk().LSBLK()

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: list})
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
// @Param  path formData string true "for example  /dev/sda1"
// @Success 200 {string} string "ok"
// @Router /disk/format [post]
func FormatDisk(c *gin.Context) {
	path := c.PostForm("path")

	t := c.PostForm("type")

	if len(path) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
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

// @Summary add mount point
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "for example: /dev/sda1"
// @Param  serial formData string true "disk id"
// @Success 200 {string} string "ok"
// @Router /disk/mount [post]
func PostMountDisk(c *gin.Context) {
	// for example: path=/dev/sda1
	path := c.PostForm("path")
	serial := c.PostForm("serial")

	mountPath := "/mnt/volume"
	var list = service.MyService.Disk().GetSerialAll()
	var pathMapList = make(map[string]string, len(list))
	for _, v := range list {
		pathMapList[v.MountPoint] = "1"
	}

	for i := 0; i < len(list)+1; i++ {
		if _, ok := pathMapList[mountPath+strconv.Itoa(i)]; !ok {
			mountPath = mountPath + strconv.Itoa(i)
			break
		}
	}

	//mount dir
	service.MyService.Disk().MountDisk(path, mountPath)
	//save to data
	m := model2.SerialDisk{}
	m.MountPoint = mountPath
	m.Path = path
	m.Serial = serial
	m.State = 0
	service.MyService.Disk().SaveMountPoint(m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary remove mount point
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "for example: /dev/sda1"
// @Param  mount_point formData string true "for example: /mnt/volume1"
// @Success 200 {string} string "ok"
// @Router /disk/umount [post]
func PostDiskUmount(c *gin.Context) {

	//
	path := c.PostForm("path")
	mountPoint := c.PostForm("mount_point")
	service.MyService.Disk().UmountPointAndRemoveDir(path)

	//delete data
	service.MyService.Disk().DeleteMountPoint(path, mountPoint)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary confirm delete disk
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Param  id path string true "id"
// @Success 200 {string} string "ok"
// @Router /disk/remove/{id} [delete]
func DeleteDisk(c *gin.Context) {
	id := c.Param("id")
	service.MyService.Disk().DeleteMount(id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary check mount point
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/init [get]
func GetDiskCheck(c *gin.Context) {

	dbList := service.MyService.Disk().GetSerialAll()
	list := service.MyService.Disk().LSBLK()

	mapList := make(map[string]string)

	for _, v := range list {
		mapList[v.Serial] = "1"
	}

	for _, v := range dbList {
		if _, ok := mapList[v.Serial]; !ok {
			//disk undefind
			c.JSON(http.StatusOK, model.Result{Success: oasis_err.ERROR, Message: oasis_err.GetMsg(oasis_err.ERROR), Data: "disk undefind"})
			return
		}
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}
