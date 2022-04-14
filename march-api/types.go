package main

type Response struct {
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type HostInfo struct {
	Hostname        string `json:"hostname"`
	Uptime          uint64 `json:"uptime"`
	BootTime        uint64 `json:"bootTime"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platformFamily"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	KernelArch      string `json:"kernelArch"`
	HostID          string `json:"hostId"`
	MaxMem          string `json:"maxMem"`
	CpuCore         int    `json:"cpuCore"`
	CpuThread       int    `json:"cpuThread"`
}

type DiskUsage struct {
	Device      string  `json:"device"`
	FsType      string  `json:"fsType"`
	MountPoint  string  `json:"mountPoint"`
	TotalBytes  uint64  `json:"totalBytes"`
	FreeBytes   uint64  `json:"freeBytes"`
	UsedBytes   uint64  `json:"usedBytes"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskIO struct {
	Device  string  `json:"device"`
	ReadIO  uint64  `json:"readIO"`
	WriteIO uint64  `json:"writeIO"`
	IOWait  float64 `json:"ioWait"`
}

type NetworkLoad struct {
	Device     string  `json:"device"`
	SentBytes  float64 `json:"sentBytes"`
	RecvBytes  float64 `json:"recvBytes"`
	SentPacket float64 `json:"sentPacket"`
	RecvPacket float64 `json:"recvPacket"`
}

type ProcessInfo struct {
	PID        uint64  `json:"pid"`
	RunUser    string  `json:"runUser"`
	CpuPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
	Command    string  `json:"command"`
}

type ChartData[T uint64 | float64] struct {
	Datetime string `json:"datetime"`
	Name     string `json:"name"`
	Value    T      `json:"value"`
}

/* type ProcessInfo struct {
	PID         int32    `json:"pid"`
	ProcessName string   `json:"processName"`
	RunUser     string   `json:"runUser"`
	CpuPercent  float64  `json:"cpuPercent"`
	MemPercent  float32  `json:"memPercent"`
	RSS         uint64   `json:"rss"`
	VMS         uint64   `json:"vms"`
	Status      []string `json:"status"`
	CreateTime  int64    `json:"createTime"`
	Command     string   `json:"command"`
}
*/
