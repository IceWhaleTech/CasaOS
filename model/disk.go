package model

type LSBLKModel struct {
	Name        string       `json:"name"`
	FsType      string       `json:"fstype"`
	Size        uint64       `json:"size"`
	FSSize      string       `json:"fssize"`
	Path        string       `json:"path"`
	Model       string       `json:"model"` //设备标识符
	RM          bool         `json:"rm"`    //是否为可移动设备
	RO          bool         `json:"ro"`    //是否为只读设备
	State       string       `json:"state"`
	PhySec      int          `json:"phy-sec"` //物理扇区大小
	Type        string       `json:"type"`
	Vendor      string       `json:"vendor"`  //供应商
	Rev         string       `json:"rev"`     //修订版本
	FSAvail     string       `json:"fsavail"` //可用空间
	FSUse       string       `json:"fsuse%"`  //已用百分比
	MountPoint  string       `json:"mountpoint"`
	Format      string       `json:"format"`
	Health      string       `json:"health"`
	HotPlug     bool         `json:"hotplug"`
	UUID        string       `json:"uuid"`
	FSUsed      string       `json:"fsused"`
	Temperature int          `json:"temperature"`
	Tran        string       `json:"tran"`
	MinIO       uint64       `json:"min-io"`
	UsedPercent float64      `json:"used_percent"`
	Serial      string       `json:"serial"`
	Children    []LSBLKModel `json:"children"`
	SubSystems  string       `json:"subsystems"`
	//详情特有
	StartSector uint64 `json:"start_sector,omitempty"`
	Rota        bool   `json:"rota"` //true(hhd) false(ssd)
	DiskType    string `json:"disk_type"`
	EndSector   uint64 `json:"end_sector,omitempty"`
}

type Drive struct {
	Name        string `json:"name"`
	Size        uint64 `json:"size"`
	Model       string `json:"model"`
	Health      string `json:"health"`
	Temperature int    `json:"temperature"`
	DiskType    string `json:"disk_type"`
	NeedFormat  bool   `json:"need_format"`
	Serial      string `json:"serial"`
	Path        string `json:"path"`
}

type DriveUSB struct {
	Name  string `json:"name"`
	Size  uint64 `json:"size"`
	Used  uint64 `json:"use"` // @tiger - 改成 used_space
	Model string `json:"model"`
	Mount bool   `json:"mount"` //是否完全挂载
	Avail uint64 `json:"avail"` //可用空间 // @tiger - 改成 available_space
}

type Storage struct {
	Name       string `json:"name"`
	MountPoint string `json:"mountpoint"`
	Size       string `json:"size"`
	Avail      string `json:"avail"` //可用空间 // @tiger - 改成 available_space
	Type       string `json:"type"`
	CreatedAt  int64  `json:"create_at"`
	Path       string `json:"path"`
	DriveName  string `json:"drive_name"`
}

type Summary struct {
	Size   uint64 `json:"size"`
	Avail  uint64 `json:"avail"` //可用空间 	// @tiger - 改成 available_space
	Health bool   `json:"health"`
	Used   uint64 `json:"used"` // @tiger - 改成 used_space
}
