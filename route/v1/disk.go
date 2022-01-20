package v1

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
)

var diskMap = make(map[string]string)

// @Summary disk list
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/list [get]
func GetDiskList(c *gin.Context) {
	list := service.MyService.Disk().LSBLK()
	newList := []model.LSBLKModel{}
	for i := len(list) - 1; i >= 0; i-- {
		if list[i].Rota {
			list[i].DiskType = "HDD"
		} else {
			list[i].DiskType = "SSD"
		}
		if list[i].Tran == "sata" {
			temp := service.MyService.Disk().SmartCTL(list[i].Path)

			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				continue
			}
			if len(list[i].Children) == 1 && len(list[i].Children[0].MountPoint) > 0 {
				pathArr := strings.Split(list[i].Children[0].MountPoint, "/")
				if len(pathArr) == 3 {
					list[i].Children[0].Name = pathArr[2]
				}
			}

			list[i].Temperature = temp.Temperature.Current
			list[i].Health = strconv.FormatBool(temp.SmartStatus.Passed)

			newList = append(newList, list[i])
		} else if len(list[i].Children) > 0 && list[i].Children[0].MountPoint == "/" {
			//system
			list[i].Children[0].Name = "System"
			list[i].Model = "System"
			list[i].DiskType = "EMMC"
			list[i].Health = "true"
			newList = append(newList, list[i])

		}
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: newList})
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

// @Summary format storage
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "e.g. /dev/sda1"
// @Param  pwd formData string true "user password"
// @Param  volume formData string true "mount point"
// @Success 200 {string} string "ok"
// @Router /disk/format [post]
func FormatDisk(c *gin.Context) {
	path := c.PostForm("path")
	t := "ext4"
	pwd := c.PostForm("pwd")
	volume := c.PostForm("volume")

	if pwd != config.UserInfo.PWD {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.PWD_INVALID, Message: oasis_err.GetMsg(oasis_err.PWD_INVALID)})
		return
	}

	if len(path) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	if _, ok := diskMap[path]; ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.DISK_BUSYING, Message: oasis_err.GetMsg(oasis_err.DISK_BUSYING)})
		return
	}
	diskMap[path] = "busying"
	service.MyService.Disk().UmountPointAndRemoveDir(path)
	format := service.MyService.Disk().FormatDisk(path, t)
	if len(format) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.FORMAT_ERROR, Message: oasis_err.GetMsg(oasis_err.FORMAT_ERROR)})
		delete(diskMap, path)
		return
	}
	service.MyService.Disk().MountDisk(path, volume)
	service.MyService.Disk().RemoveLSBLKCache()
	delete(diskMap, path)
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

// @Summary  add storage
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "disk path  e.g. /dev/sda"
// @Param  serial formData string true "serial"
// @Param  name formData string true "name"
// @Param  format formData bool true "need format(true)"
// @Success 200 {string} string "ok"
// @Router /disk/storage [post]
func AddPartition(c *gin.Context) {
	name := c.PostForm("name")
	path := c.PostForm("path")
	serial := c.PostForm("serial")
	format, _ := strconv.ParseBool(c.PostForm("format"))
	if len(name) == 0 || len(path) == 0 || len(serial) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	if _, ok := diskMap[serial]; ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.DISK_BUSYING, Message: oasis_err.GetMsg(oasis_err.DISK_BUSYING)})
		return
	}
	if !file.CheckNotExist("/mnt/" + name) {
		// /mnt/name exist
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.NAME_NOT_AVAILABLE, Message: oasis_err.GetMsg(oasis_err.NAME_NOT_AVAILABLE)})
		return
	}
	diskMap[serial] = "busying"
	currentDisk := service.MyService.Disk().GetDiskInfo(path)
	if !format {
		if len(currentDisk.Children) != 1 || !(len(currentDisk.Children) > 0 && currentDisk.Children[0].FsType == "ext4") {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err.DISK_NEEDS_FORMAT, Message: oasis_err.GetMsg(oasis_err.DISK_NEEDS_FORMAT)})
			delete(diskMap, serial)
			return
		}
	} else {
		service.MyService.Disk().AddPartition(path)
	}

	mountPath := "/mnt/" + name

	service.MyService.Disk().MountDisk(path, mountPath)

	m := model2.SerialDisk{}
	m.MountPoint = mountPath
	m.Path = path + "1"
	m.Serial = serial
	m.State = 0
	service.MyService.Disk().SaveMountPoint(m)

	//mount dir
	service.MyService.Disk().MountDisk(path+"1", mountPath)

	service.MyService.Disk().RemoveLSBLKCache()

	delete(diskMap, serial)
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

	m := model2.SerialDisk{}
	m.MountPoint = mountPath
	m.Path = path
	m.Serial = serial
	m.State = 0
	//service.MyService.Disk().SaveMountPoint(m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary remove mount point
// @Produce  application/json
// @Accept multipart/form-data
// @Tags disk
// @Security ApiKeyAuth
// @Param  path formData string true "e.g. /dev/sda1"
// @Param  mount_point formData string true "e.g. /mnt/volume1"
// @Param  pwd formData string true "user password"
// @Success 200 {string} string "ok"
// @Router /disk/umount [post]
func PostDiskUmount(c *gin.Context) {

	path := c.PostForm("path")
	mountPoint := c.PostForm("volume")
	pwd := c.PostForm("pwd")

	if len(path) == 0 || len(mountPoint) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	if pwd != config.UserInfo.PWD {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.PWD_INVALID, Message: oasis_err.GetMsg(oasis_err.PWD_INVALID)})
		return
	}

	if _, ok := diskMap[path]; ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.DISK_BUSYING, Message: oasis_err.GetMsg(oasis_err.DISK_BUSYING)})
		return
	}

	service.MyService.Disk().UmountPointAndRemoveDir(path)
	//delete data
	service.MyService.Disk().DeleteMountPoint(path, mountPoint)
	service.MyService.Disk().RemoveLSBLKCache()
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
