export interface IResponse {
	code: number;
	data: any;
	message: string;
}

export interface IHostInfo {
	bootTime: number;
	cpuCore: number;
	cpuThread: number;
	hostId: string;
	hostname: string;
	kernelArch: string;
	kernelVersion: string;
	maxMem: string;
	os: string;
	platform: string;
	platformFamily: string;
	platformVersion: string;
	uptime: number;
}

export interface ICpuInfo {
	modelName: string;
	arch: string;
	frequency: number;
	cacheSize: number;
	flags: string[];
}

export interface ISingleCpuPercent {
	id: string;
	singleCpuPercent: number;
	timestamp: number;
}
export interface IMultiCpuPercent {
	id: string;
	multiCpuPercent: number[];
	timestamp: number;
}

export interface IMemStat {
	active: number;
	available: number;
	buffers: number;
	cached: number;
	commitLimit: number;
	committedAS: number;
	dirty: number;
	free: number;
	highFree: number;
	highTotal: number;
	hugePageSize: number;
	hugePagesFree: number;
	hugePagesTotal: number;
	inactive: number;
	laundry: number;
	lowFree: number;
	lowTotal: number;
	mapped: number;
	pageTables: number;
	shared: number;
	slab: number;
	sreclaimable: number;
	sunreclaim: number;
	swapCached: number;
	swapFree: number;
	swapTotal: number;
	total: number;
	used: number;
	usedPercent: number;
	vmallocChunk: number;
	vmallocTotal: number;
	vmallocUsed: number;
	wired: number;
	writeBack: number;
	writeBackTmp: number;
}

export interface IMemPercent {
	id: string;
	memStat: IMemStat;
	timestamp: number;
}

export interface IDiskUsage {
	device: string;
	fsType: string;
	mountPoint: string;
	totalBytes: number;
	freeBytes: number;
	usedBytes: number;
	usedPercent: number;
}

export interface IDiskStat {
	device: string;
	ioWait: number;
	readIO: number;
	writeIO: number;
}

export interface IDiskIO {
	id: string;
	diskIO: IDiskStat[];
	timestamp: number;
}

export interface INetworkStat {
	device: string;
	recvBytes: number;
	recvPacket: number;
	sentBytes: number;
	sentPacket: number;
}

export interface INetworkLoad {
	id: string;
	networkLoad: INetworkStat[];
	timestamp: number;
}

export interface IProcessInfo {
	pid: string;
	runUser: string;
	cpuPercent: string;
	memPercent: string;
	command: string;
}
