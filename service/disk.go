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
	LSBLK(isUseCache bool) []model.LSBLKModel
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
func (d *diskService) LSBLK(isUseCache bool) []model.LSBLKModel {
	key := "system_lsblk"
	var n []model.LSBLKModel

	if result, ok := Cache.Get(key); ok && isUseCache {

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
	// strStr := `{
	// 	"blockdevices": [
	// 	   {"name":"loop0", "kname":"loop0", "path":"/dev/loop0", "maj:min":"7:0", "fsavail":"0", "fssize":"62M", "fstype":"squashfs", "fsused":"62M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/core20/1405", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":619, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop1", "kname":"loop1", "path":"/dev/loop1", "maj:min":"7:1", "fsavail":"0", "fssize":"55.6M", "fstype":"squashfs", "fsused":"55.6M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/core18/2344", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":55, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop2", "kname":"loop2", "path":"/dev/loop2", "maj:min":"7:2", "fsavail":"0", "fssize":"44.8M", "fstype":"squashfs", "fsused":"44.8M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/snapd/15314", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":446, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop3", "kname":"loop3", "path":"/dev/loop3", "maj:min":"7:3", "fsavail":"0", "fssize":"78.9M", "fstype":"squashfs", "fsused":"78.9M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/lxd/22754", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":788, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop4", "kname":"loop4", "path":"/dev/loop4", "maj:min":"7:4", "fsavail":"0", "fssize":"43.8M", "fstype":"squashfs", "fsused":"43.8M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/snapd/15177", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":436, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop5", "kname":"loop5", "path":"/dev/loop5", "maj:min":"7:5", "fsavail":"0", "fssize":"55.5M", "fstype":"squashfs", "fsused":"55.5M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/core18/1997", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":554, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"loop6", "kname":"loop6", "path":"/dev/loop6", "maj:min":"7:6", "fsavail":"0", "fssize":"80M", "fstype":"squashfs", "fsused":"80M", "fsuse%":"100%", "fsver":"4.0", "mountpoint":"/snap/lxd/22826", "label":null, "uuid":null, "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":true, "rm":false, "hotplug":false, "model":null, "serial":null, "size":799, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"mq-deadline", "rq-size":256, "type":"loop", "disc-aln":0, "disc-gran":"4K", "disc-max":"4G", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":null, "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 	   {"name":"sda", "kname":"sda", "path":"/dev/sda", "maj:min":"8:0", "fsavail":null, "fssize":null, "fstype":null, "fsused":null, "fsuse%":null, "fsver":null, "mountpoint":null, "label":null, "uuid":null, "ptuuid":"1596101a-e20d-4296-96e2-0870efce554a", "pttype":"gpt", "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":"ST1000DM003-1ER1", "serial":"Z4YCS1B6", "size":9315, "state":"running", "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":4096, "opt-io":0, "phy-sec":4096, "log-sec":512, "rota":true, "sched":"mq-deadline", "rq-size":64, "type":"disk", "disc-aln":0, "disc-gran":"0B", "disc-max":"0B", "disc-zero":false, "wsame":"0B", "wwn":"0x5000c50090db103a", "rand":true, "pkname":null, "hctl":"0:0:0:0", "tran":"sata", "subsystems":"block:scsi:pci", "rev":"CC61", "vendor":"ATA     ", "zoned":"none", "dax":false,
	// 		  "children": [
	// 			 {"name":"sda1", "kname":"sda1", "path":"/dev/sda1", "maj:min":"8:1", "fsavail":null, "fssize":null, "fstype":null, "fsused":null, "fsuse%":null, "fsver":null, "mountpoint":null, "label":null, "uuid":null, "ptuuid":"1596101a-e20d-4296-96e2-0870efce554a", "pttype":"gpt", "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":9315, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":3072, "min-io":4096, "opt-io":0, "phy-sec":4096, "log-sec":512, "rota":true, "sched":"mq-deadline", "rq-size":64, "type":"part", "disc-aln":0, "disc-gran":"0B", "disc-max":"0B", "disc-zero":false, "wsame":"0B", "wwn":"0x5000c50090db103a", "rand":true, "pkname":"sda", "hctl":null, "tran":null, "subsystems":"block:scsi:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false}
	// 		  ]
	// 	   },
	// 	   {"name":"sdb", "kname":"sdb", "path":"/dev/sdb", "maj:min":"8:16", "fsavail":null, "fssize":null, "fstype":null, "fsused":null, "fsuse%":null, "fsver":null, "mountpoint":null, "label":null, "uuid":null, "ptuuid":"baed02d0-e92d-4a00-9609-f94f31271a0e", "pttype":"gpt", "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":"ST1000DM003-1ER1", "serial":"W4Y51MFH", "size":9315, "state":"running", "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":4096, "opt-io":0, "phy-sec":4096, "log-sec":512, "rota":true, "sched":"mq-deadline", "rq-size":64, "type":"disk", "disc-aln":0, "disc-gran":"0B", "disc-max":"0B", "disc-zero":false, "wsame":"0B", "wwn":"0x5000c5008acd2f00", "rand":true, "pkname":null, "hctl":"1:0:0:0", "tran":"sata", "subsystems":"block:scsi:pci", "rev":"CC46", "vendor":"ATA     ", "zoned":"none", "dax":false,
	// 		  "children": [
	// 			 {"name":"sdb1", "kname":"sdb1", "path":"/dev/sdb1", "maj:min":"8:17", "fsavail":null, "fssize":null, "fstype":"zfs_member", "fsused":null, "fsuse%":null, "fsver":"5000", "mountpoint":null, "label":null, "uuid":null, "ptuuid":"baed02d0-e92d-4a00-9609-f94f31271a0e", "pttype":"gpt", "parttype":"0fc63daf-8483-4772-8e79-3d69d8477de4", "parttypename":"Linux filesystem", "partlabel":"primary", "partuuid":"57880cc0-2695-41c3-bf14-7161693e5bff", "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":9315, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":3072, "min-io":4096, "opt-io":0, "phy-sec":4096, "log-sec":512, "rota":true, "sched":"mq-deadline", "rq-size":64, "type":"part", "disc-aln":0, "disc-gran":"0B", "disc-max":"0B", "disc-zero":false, "wsame":"0B", "wwn":"0x5000c5008acd2f00", "rand":true, "pkname":"sdb", "hctl":null, "tran":null, "subsystems":"block:scsi:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false}
	// 		  ]
	// 	   },
	// 	   {"name":"nvme0n1", "kname":"nvme0n1", "path":"/dev/nvme0n1", "maj:min":"259:0", "fsavail":null, "fssize":null, "fstype":null, "fsused":null, "fsuse%":null, "fsver":null, "mountpoint":null, "label":null, "uuid":null, "ptuuid":"338abc31-a3d4-4af2-9342-b53268d9e5ac", "pttype":"gpt", "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":"LITEON CL1-8D128-HP", "serial":"UJDJA01PJDH3UI", "size":1192, "state":"live", "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"none", "rq-size":255, "type":"disk", "disc-aln":0, "disc-gran":"512B", "disc-max":"2T", "disc-zero":false, "wsame":"0B", "wwn":"eui.0023035630392fe7", "rand":false, "pkname":null, "hctl":null, "tran":"nvme", "subsystems":"block:nvme:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false,
	// 		  "children": [
	// 			 {"name":"nvme0n1p1", "kname":"nvme0n1p1", "path":"/dev/nvme0n1p1", "maj:min":"259:1", "fsavail":null, "fssize":null, "fstype":null, "fsused":null, "fsuse%":null, "fsver":null, "mountpoint":null, "label":null, "uuid":null, "ptuuid":"338abc31-a3d4-4af2-9342-b53268d9e5ac", "pttype":"gpt", "parttype":"21686148-6449-6e6f-744e-656564454649", "parttypename":"BIOS boot", "partlabel":null, "partuuid":"b2bac638-9468-449f-9669-79be44e3c80d", "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":1, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"none", "rq-size":255, "type":"part", "disc-aln":0, "disc-gran":"512B", "disc-max":"2T", "disc-zero":false, "wsame":"0B", "wwn":"eui.0023035630392fe7", "rand":false, "pkname":"nvme0n1", "hctl":null, "tran":"nvme", "subsystems":"block:nvme:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 			 {"name":"nvme0n1p2", "kname":"nvme0n1p2", "path":"/dev/nvme0n1p2", "maj:min":"259:2", "fsavail":"1.2G", "fssize":"1.4G", "fstype":"ext4", "fsused":"131.5M", "fsuse%":"9%", "fsver":"1.0", "mountpoint":"/boot", "label":null, "uuid":"cd8164e9-bf7a-4684-8a3b-1d9f209b4930", "ptuuid":"338abc31-a3d4-4af2-9342-b53268d9e5ac", "pttype":"gpt", "parttype":"0fc63daf-8483-4772-8e79-3d69d8477de4", "parttypename":"Linux filesystem", "partlabel":null, "partuuid":"42ed9ed6-1221-4bea-901a-bc2f7b7cb9e1", "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":15, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"none", "rq-size":255, "type":"part", "disc-aln":0, "disc-gran":"512B", "disc-max":"2T", "disc-zero":false, "wsame":"0B", "wwn":"eui.0023035630392fe7", "rand":false, "pkname":"nvme0n1", "hctl":null, "tran":"nvme", "subsystems":"block:nvme:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false},
	// 			 {"name":"nvme0n1p3", "kname":"nvme0n1p3", "path":"/dev/nvme0n1p3", "maj:min":"259:3", "fsavail":null, "fssize":null, "fstype":"LVM2_member", "fsused":null, "fsuse%":null, "fsver":"LVM2 001", "mountpoint":null, "label":null, "uuid":"0G7ryL-p2Ks-i9HS-wvEO-lXHs-oyZX-1KIlZO", "ptuuid":"338abc31-a3d4-4af2-9342-b53268d9e5ac", "pttype":"gpt", "parttype":"0fc63daf-8483-4772-8e79-3d69d8477de4", "parttypename":"Linux filesystem", "partlabel":null, "partuuid":"f6ae2e8c-14ae-4d94-89f1-2c154e909843", "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":1177, "state":null, "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":"none", "rq-size":255, "type":"part", "disc-aln":0, "disc-gran":"512B", "disc-max":"2T", "disc-zero":false, "wsame":"0B", "wwn":"eui.0023035630392fe7", "rand":false, "pkname":"nvme0n1", "hctl":null, "tran":"nvme", "subsystems":"block:nvme:pci", "rev":null, "vendor":null, "zoned":"none", "dax":false,
	// 				"children": [
	// 				   {"name":"ubuntu--vg-ubuntu--lv", "kname":"dm-0", "path":"/dev/mapper/ubuntu--vg-ubuntu--lv", "maj:min":"253:0", "fsavail":"78.6G", "fssize":"115.4G", "fstype":"ext4", "fsused":"30.9G", "fsuse%":"27%", "fsver":"1.0", "mountpoint":"/", "label":null, "uuid":"e8a9082f-3643-4820-a5e5-05817d7738c6", "ptuuid":null, "pttype":null, "parttype":null, "parttypename":null, "partlabel":null, "partuuid":null, "partflags":null, "ra":128, "ro":false, "rm":false, "hotplug":false, "model":null, "serial":null, "size":1177, "state":"running", "owner":"root", "group":"disk", "mode":"brw-rw----", "alignment":0, "min-io":512, "opt-io":0, "phy-sec":512, "log-sec":512, "rota":false, "sched":null, "rq-size":128, "type":"lvm", "disc-aln":0, "disc-gran":"512B", "disc-max":"2T", "disc-zero":false, "wsame":"0B", "wwn":null, "rand":false, "pkname":"nvme0n1p3", "hctl":null, "tran":null, "subsystems":"block", "rev":null, "vendor":null, "zoned":"none", "dax":false}
	// 				]
	// 			 }
	// 		  ]
	// 	   }
	// 	]
	//  }`
	// fmt.Println(gjson.Get(strStr, "blockdevices").String())
	err := json2.Unmarshal([]byte(gjson.Get(string(str), "blockdevices").String()), &m)
	if err != nil {
		fmt.Println(err)
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
				if err != nil {
					d.log.Fatal("diskservice_lsblk_fsused", err)
				}
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
	d.db.Where("uuid = ?", m.UUID).Delete(&model2.SerialDisk{})
	d.db.Create(&m)
}

func (d *diskService) UpdateMountPoint(m model2.SerialDisk) {
	d.db.Model(&model2.SerialDisk{}).Where("uui = ?", m.UUID).Update("mount_point", m.MountPoint)
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
