package service

import (
	json2 "encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type DiskService interface {
	GetPlugInDisk() []string
	LSBLK() []model.LSBLKModel
	SmartCTL(path string) model.SmartctlA
	FormatDisk(path, format string) []string
	UmountPointAndRemoveDir(path string) []string
	GetDiskInfo(path string) model.LSBLKModel
	DelPartition(path, num string) string
	AddPartition(path string) string
	GetDiskInfoByPath(path string) *disk.UsageStat
	MountDisk(path, volume string)
	GetSerialAll() []model2.SerialDisk
	SaveMountPoint(m model2.SerialDisk)
	DeleteMountPoint(path, mountPoint string)
	DeleteMount(id string)
	UpdateMountPoint(m model2.SerialDisk)
	RemoveLSBLKCache()
}
type diskService struct {
	log loger2.OLog
	db  *gorm.DB
}

func (d *diskService) RemoveLSBLKCache() {
	key := "system_lsblk"
	Cache.Delete(key)
}
func (d *diskService) SmartCTL(path string) model.SmartctlA {

	key := "system_smart_" + path
	if result, ok := Cache.Get(key); ok {

		res, ok := result.(model.SmartctlA)
		if ok {
			return res
		}
	}
	var m model.SmartctlA
	str := command2.ExecSmartCTLByPath(path)
	if str == nil {
		d.log.Error("smartctl exec error,smartctl")
		return m
	}

	err := json2.Unmarshal([]byte(str), &m)
	if err != nil {
		d.log.Error("json ummarshal error", err)
	}
	if !reflect.DeepEqual(m, model.SmartctlA{}) {
		Cache.Add(key, m, time.Second*10)
	}
	return m
}

//通过脚本获取外挂磁盘
func (d *diskService) GetPlugInDisk() []string {
	return command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetPlugInDisk")
}

//格式化硬盘
func (d *diskService) FormatDisk(path, format string) []string {
	r := command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;FormatDisk " + path + " " + format)
	return r
}

//移除挂载点,删除目录
func (d *diskService) UmountPointAndRemoveDir(path string) []string {
	r := command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;UMountPorintAndRemoveDir " + path)
	return r
}

//删除分区
func (d *diskService) DelPartition(path, num string) string {
	r := command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;DelPartition " + path + " " + num)
	fmt.Println(r)
	return ""
}

//part
func (d *diskService) AddPartition(path string) string {
	command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;AddPartition " + path)
	return ""
}

func (d *diskService) AddAllPartition(path string) {

}

//获取硬盘详情
func (d *diskService) GetDiskInfoByPath(path string) *disk.UsageStat {
	diskInfo, err := disk.Usage(path + "1")

	if err != nil {
		fmt.Println(err)
	}
	diskInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.UsedPercent), 64)
	diskInfo.InodesUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.InodesUsedPercent), 64)
	return diskInfo
}

//get disk details
func (d *diskService) LSBLK() []model.LSBLKModel {
	key := "system_lsblk"
	var n []model.LSBLKModel

	if result, ok := Cache.Get(key); ok {

		res, ok := result.([]model.LSBLKModel)
		if ok {
			return res
		}
	}

	str := command2.ExecLSBLK()
	if str == nil {
		d.log.Error("lsblk exec error,lsblk")
		return nil
	}
	var m []model.LSBLKModel
	err := json2.Unmarshal([]byte(gjson.Get(string(str), "blockdevices").String()), &m)
	if err != nil {
		d.log.Error("json ummarshal error", err)
	}

	var c []model.LSBLKModel

	var fsused uint64

	var health = true
	for _, i := range m {
		if i.Type != "loop" && !i.RO {
			fsused = 0
			for _, child := range i.Children {
				if child.RM {
					child.Health = strings.TrimSpace(command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetDiskHealthState " + child.Path))
					if strings.ToLower(strings.TrimSpace(child.State)) != "ok" {
						health = false
					}
					f, _ := strconv.ParseUint(child.FSUsed, 10, 64)
					fsused += f
				} else {
					health = false
				}
				c = append(c, child)
			}
			i.Format = strings.TrimSpace(command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetDiskType " + i.Path))
			if health {
				i.Health = "OK"
			}
			i.FSUsed = strconv.FormatUint(fsused, 10)
			i.Children = c
			if fsused > 0 {
				i.UsedPercent, err = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(fsused)/float64(i.Size)), 64)
				d.log.Fatal("diskservice_lsblk_fsused", err)
			}
			n = append(n, i)
			health = true
			c = []model.LSBLKModel{}
			fsused = 0
		}
	}
	if len(n) > 0 {
		Cache.Add(key, n, time.Second*100)
	}
	return n
}

func (d *diskService) GetDiskInfo(path string) model.LSBLKModel {
	str := command2.ExecLSBLKByPath(path)
	if str == nil {
		d.log.Error("lsblk exec error,str")
		return model.LSBLKModel{}
	}

	var ml []model.LSBLKModel
	err := json2.Unmarshal([]byte(gjson.Get(string(str), "blockdevices").String()), &ml)
	if err != nil {
		d.log.Info(string(str))
		d.log.Error("json ummarshal error", err)
		return model.LSBLKModel{}
	}

	m := model.LSBLKModel{}
	if len(ml) > 0 {
		m = ml[0]
	}
	return m
	// 下面为计算是否可以继续分区的部分,暂时不需要
	chiArr := make(map[string]string)
	chiList := command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetPartitionSectors " + m.Path)
	if len(chiList) == 0 {
		d.log.Error(m.Path, chiList)
		d.log.Error("chiList length error")
	}
	for i := 0; i < len(chiList); i++ {
		tempArr := strings.Split(chiList[i], ",")
		chiArr[tempArr[0]] = chiList[i]
	}
	var maxSector uint64 = 0
	for i := 0; i < len(m.Children); i++ {
		tempArr := strings.Split(chiArr[m.Children[i].Path], ",")
		m.Children[i].StartSector, _ = strconv.ParseUint(tempArr[1], 10, 64)
		m.Children[i].EndSector, _ = strconv.ParseUint(tempArr[2], 10, 64)
		if m.Children[i].EndSector > maxSector {
			maxSector = m.Children[i].EndSector
		}

	}
	diskEndSector := command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetDiskSizeAndSectors " + m.Path)

	if len(diskEndSector) < 2 {
		d.log.Error("diskEndSector length error")
	}
	diskEndSectorInt, _ := strconv.ParseUint(diskEndSector[len(diskEndSector)-1], 10, 64)
	if (diskEndSectorInt-maxSector)*m.MinIO/1024/1024 > 100 {
		//添加可以分区情况
		p := model.LSBLKModel{}
		p.Path = "可以添加"
		m.Children = append(m.Children, p)
	}
	return m
}

func (d *diskService) MountDisk(path, volume string) {
	r := command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;do_mount " + path + " " + volume)
	fmt.Print(r)
}

func (d *diskService) SaveMountPoint(m model2.SerialDisk) {
	d.db.Where("serial = ?", m.Serial).Delete(&model2.SerialDisk{})
	d.db.Create(&m)
}

func (d *diskService) UpdateMountPoint(m model2.SerialDisk) {
	d.db.Model(&model2.SerialDisk{}).Where("serial = ?", m.Serial).Update("mount_point", m.MountPoint)
}

func (d *diskService) DeleteMount(id string) {

	d.db.Delete(&model2.SerialDisk{}).Where("id = ?", id)
}

func (d *diskService) DeleteMountPoint(path, mountPoint string) {

	d.db.Where("path = ? AND mount_point = ?", path, mountPoint).Delete(&model2.SerialDisk{})

	command2.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;do_umount " + path)
}

func (d *diskService) GetSerialAll() []model2.SerialDisk {
	var m []model2.SerialDisk
	d.db.Find(&m)
	return m
}

func NewDiskService(log loger2.OLog, db *gorm.DB) DiskService {
	return &diskService{log: log, db: db}
}
