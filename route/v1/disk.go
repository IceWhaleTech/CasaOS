package v1

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
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
	list := service.MyService.Disk().LSBLK(false)
	dbList := service.MyService.Disk().GetSerialAll()
	part := make(map[string]int64, len(dbList))
	for _, v := range dbList {
		part[v.MountPoint] = v.CreatedAt
	}
	findSystem := 0

	disks := []model.Drive{}
	storage := []model.Storage{}
	avail := []model.Drive{}
	for i := 0; i < len(list); i++ {
		disk := model.Drive{}
		if list[i].Rota {
			disk.DiskType = "HDD"
		} else {
			disk.DiskType = "SSD"
		}
		disk.Serial = list[i].Serial
		disk.Name = list[i].Name
		disk.Size = list[i].Size
		disk.Path = list[i].Path
		disk.Model = list[i].Model

		if len(list[i].Children) > 0 && findSystem == 0 {
			for j := 0; j < len(list[i].Children); j++ {
				if len(list[i].Children[j].Children) > 0 {
					for _, v := range list[i].Children[j].Children {
						if v.MountPoint == "/" {
							stor := model.Storage{}
							stor.Name = "System"
							stor.MountPoint = v.MountPoint
							stor.Size = v.FSSize
							stor.Avail = v.FSAvail
							stor.Path = v.Path
							stor.Type = v.FsType
							stor.DriveName = "System"
							disk.Model = "System"
							if strings.Contains(v.SubSystems, "mmc") {
								disk.DiskType = "MMC"
							} else if strings.Contains(v.SubSystems, "usb") {
								disk.DiskType = "USB"
							}
							disk.Health = "true"

							disks = append(disks, disk)
							storage = append(storage, stor)
							findSystem = 1
							break
						}
					}
				} else {
					if list[i].Children[j].MountPoint == "/" {
						stor := model.Storage{}
						stor.Name = "System"
						stor.MountPoint = list[i].Children[j].MountPoint
						stor.Size = list[i].Children[j].FSSize
						stor.Avail = list[i].Children[j].FSAvail
						stor.Path = list[i].Children[j].Path
						stor.Type = list[i].Children[j].FsType
						stor.DriveName = "System"
						disk.Model = "System"
						if strings.Contains(list[i].Children[j].SubSystems, "mmc") {
							disk.DiskType = "MMC"
						} else if strings.Contains(list[i].Children[j].SubSystems, "usb") {
							disk.DiskType = "USB"
						}
						disk.Health = "true"

						disks = append(disks, disk)
						storage = append(storage, stor)
						findSystem = 1
						break
					}
				}

			}
		}
		if findSystem == 1 {
			findSystem += 1
			continue
		}

		if list[i].Tran == "sata" || list[i].Tran == "nvme" || list[i].Tran == "spi" || list[i].Tran == "sas" {
			temp := service.MyService.Disk().SmartCTL(list[i].Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				temp.SmartStatus.Passed = true
			}
			if len(list[i].Children) == 1 && len(list[i].Children[0].MountPoint) > 0 {
				stor := model.Storage{}
				stor.MountPoint = list[i].Children[0].MountPoint
				stor.Size = list[i].Children[0].FSSize
				stor.Avail = list[i].Children[0].FSAvail
				stor.Path = list[i].Children[0].Path
				stor.Type = list[i].Children[0].FsType
				stor.DriveName = list[i].Name
				pathArr := strings.Split(list[i].Children[0].MountPoint, "/")
				if len(pathArr) == 3 {
					stor.Name = pathArr[2]
				}
				if t, ok := part[list[i].Children[0].MountPoint]; ok {
					stor.CreatedAt = t
				}
				storage = append(storage, stor)
			} else {
				//todo   长度有问题
				if len(list[i].Children) == 1 && list[i].Children[0].FsType == "ext4" {
					disk.NeedFormat = false
					avail = append(avail, disk)
				} else {
					disk.NeedFormat = true
					avail = append(avail, disk)
				}
			}

			disk.Temperature = temp.Temperature.Current
			disk.Health = strconv.FormatBool(temp.SmartStatus.Passed)

			disks = append(disks, disk)
		}
	}
	data := make(map[string]interface{}, 3)
	data["drive"] = disks
	data["storage"] = storage
	data["avail"] = avail

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

// @Summary get disk list
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/lists [get]
func GetPlugInDisks(c *gin.Context) {

	list := service.MyService.Disk().LSBLK(true)
	var result []*disk.UsageStat
	for _, item := range list {
		result = append(result, service.MyService.Disk().GetDiskInfoByPath(item.Path))
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: result})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	}
	m := service.MyService.Disk().GetDiskInfo(path)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: m})
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
func PostDiskFormat(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)
	path := json["path"]
	t := "ext4"
	pwd := json["pwd"]
	volume := json["volume"]
	id := c.GetHeader("user_id")
	user := service.MyService.User().GetUserAllInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	if encryption.GetMD5ByStr(pwd) != user.Password {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PWD_INVALID, Message: common_err.GetMsg(common_err.PWD_INVALID)})
		return
	}

	if len(path) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if _, ok := diskMap[path]; ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DISK_BUSYING, Message: common_err.GetMsg(common_err.DISK_BUSYING)})
		return
	}
	diskMap[path] = "busying"
	service.MyService.Disk().UmountPointAndRemoveDir(path)
	format := service.MyService.Disk().FormatDisk(path, t)
	if len(format) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.FORMAT_ERROR, Message: common_err.GetMsg(common_err.FORMAT_ERROR)})
		delete(diskMap, path)
		return
	}
	service.MyService.Disk().MountDisk(path, volume)
	service.MyService.Disk().RemoveLSBLKCache()
	delete(diskMap, path)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: strArr})

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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	}
	var p = path[:len(path)-1]
	var n = path[len(path)-1:]
	service.MyService.Disk().DelPartition(p, n)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
func PostDiskAddPartition(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)

	name := json["name"]
	path := json["path"]
	format, _ := strconv.ParseBool(json["format"])
	if len(name) == 0 || len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if _, ok := diskMap[path]; ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DISK_BUSYING, Message: common_err.GetMsg(common_err.DISK_BUSYING)})
		return
	}
	if !file.CheckNotExist("/DATA/" + name) {
		// /mnt/name exist
		c.JSON(http.StatusOK, model.Result{Success: common_err.NAME_NOT_AVAILABLE, Message: common_err.GetMsg(common_err.NAME_NOT_AVAILABLE)})
		return
	}
	diskMap[path] = "busying"
	currentDisk := service.MyService.Disk().GetDiskInfo(path)
	if !format {
		if len(currentDisk.Children) != 1 || !(len(currentDisk.Children) > 0 && currentDisk.Children[0].FsType == "ext4") {
			c.JSON(http.StatusOK, model.Result{Success: common_err.DISK_NEEDS_FORMAT, Message: common_err.GetMsg(common_err.DISK_NEEDS_FORMAT)})
			delete(diskMap, path)
			return
		}
	} else {
		service.MyService.Disk().AddPartition(path)
	}

	formatBool := true
	for formatBool {
		currentDisk = service.MyService.Disk().GetDiskInfo(path)
		fmt.Println(currentDisk.Children)
		if len(currentDisk.Children) > 0 {
			formatBool = false
			break
		}
		time.Sleep(time.Second)
	}
	currentDisk = service.MyService.Disk().GetDiskInfo(path)
	if len(currentDisk.Children) != 1 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DISK_NEEDS_FORMAT, Message: common_err.GetMsg(common_err.DISK_NEEDS_FORMAT)})
		return
	}

	mountPath := "/DATA/" + name
	m := model2.SerialDisk{}
	m.MountPoint = mountPath
	m.Path = currentDisk.Children[0].Path
	m.UUID = currentDisk.Children[0].UUID
	m.State = 0
	m.CreatedAt = time.Now().Unix()
	service.MyService.Disk().SaveMountPoint(m)

	//mount dir
	service.MyService.Disk().MountDisk(currentDisk.Children[0].Path, mountPath)

	service.MyService.Disk().RemoveLSBLKCache()

	delete(diskMap, path)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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

	mountPath := "/DATA/volume"
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
	m.UUID = serial
	m.State = 0
	//service.MyService.Disk().SaveMountPoint(m)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if pwd != config.UserInfo.PWD {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PWD_INVALID, Message: common_err.GetMsg(common_err.PWD_INVALID)})
		return
	}

	if _, ok := diskMap[path]; ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DISK_BUSYING, Message: common_err.GetMsg(common_err.DISK_BUSYING)})
		return
	}

	service.MyService.Disk().UmountPointAndRemoveDir(path)
	//delete data
	service.MyService.Disk().DeleteMountPoint(path, mountPoint)
	service.MyService.Disk().RemoveLSBLKCache()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
	list := service.MyService.Disk().LSBLK(true)

	mapList := make(map[string]string)

	for _, v := range list {
		mapList[v.Serial] = "1"
	}

	for _, v := range dbList {
		if _, ok := mapList[v.UUID]; !ok {
			//disk undefind
			c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: "disk undefind"})
			return
		}
	}

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary check mount point
// @Produce  application/json
// @Accept application/json
// @Tags disk
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /disk/usb [get]
func GetUSBList(c *gin.Context) {
	list := service.MyService.Disk().LSBLK(false)
	data := []model.DriveUSB{}
	for _, v := range list {
		if v.Tran == "usb" {
			temp := model.DriveUSB{}
			temp.Model = v.Model
			temp.Name = v.Name
			temp.Size = v.Size
			mountTemp := true
			if len(v.Children) == 0 {
				mountTemp = false
			}
			for _, child := range v.Children {
				if len(child.MountPoint) > 0 {
					avail, _ := strconv.ParseUint(child.FSAvail, 10, 64)
					temp.Avail += avail
				} else {
					mountTemp = false
				}
			}
			temp.Mount = mountTemp
			data = append(data, temp)
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}
