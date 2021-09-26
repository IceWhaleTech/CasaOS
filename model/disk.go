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
	FSUsed      string       `json:"fsused"`
	Tran        string       `json:"tran"`
	MinIO       uint64       `json:"min-io"`
	UsedPercent float64      `json:"used_percent"`
	Children    []LSBLKModel `json:"children"`
	//详情特有
	StartSector uint64 `json:"start_sector,omitempty"`
	EndSector   uint64 `json:"end_sector,omitempty"`
}
