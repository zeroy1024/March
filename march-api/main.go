package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	config := LoadConfig()
	mgo := Connect()
	db := mgo.Database()

	go MonitorRun(db)

	app := fiber.New(fiber.Config{
		AppName: "March",
	})

	if config.LocalWeb {
		app.Static("/", "./web/")
		fmt.Println("已启动本地web")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	v1 := app.Group("/api/v1")
	{
		// 实时数据
		v1.Get("/current", func(c *fiber.Ctx) error {
			currentCpuPercent, _ := cpu.Percent(time.Second, false)
			memStat, _ := mem.VirtualMemory()
			currentMemPercent := memStat.UsedPercent

			return c.JSON(Response{
				Code: 200,
				Data: fiber.Map{
					"cpuPercent": currentCpuPercent,
					"memPercent": currentMemPercent,
				},
				Message: "success",
			})
		})

		// 主机相关
		hostGroup := v1.Group("/host")
		{
			// 系统信息
			hostGroup.Get("/", func(c *fiber.Ctx) error {
				hostInfo, err := host.Info()
				if err != nil {
					panic(err.Error())
				}
				memInfo, _ := mem.VirtualMemory()
				cpuCore, _ := cpu.Counts(false)
				cpuThread, _ := cpu.Counts(true)
				hInfo := HostInfo{
					Hostname:        hostInfo.Hostname,
					Uptime:          hostInfo.Uptime,
					BootTime:        hostInfo.BootTime,
					OS:              hostInfo.OS,
					Platform:        hostInfo.Platform,
					PlatformFamily:  hostInfo.PlatformFamily,
					PlatformVersion: hostInfo.PlatformVersion,
					KernelVersion:   hostInfo.KernelVersion,
					KernelArch:      hostInfo.KernelArch,
					HostID:          hostInfo.HostID,
					MaxMem:          fmt.Sprintf("%.2f", float64(memInfo.Total)/1024.00/1024.00/1024.00) + "G",
					CpuCore:         cpuCore,
					CpuThread:       cpuThread,
				}

				return c.JSON(Response{
					Code:    200,
					Data:    hInfo,
					Message: "success",
				})
			})

			// 温度相关
			hostGroup.Get("/temp", func(c *fiber.Ctx) error {
				return c.JSON(Response{
					Code:    404,
					Message: "未完成",
				})
			})
		}

		// CPU相关
		cpuGroup := v1.Group("/cpu")
		{
			cpuGroup.Get("/", func(c *fiber.Ctx) error {
				cpuInfo, err := cpu.Info()
				if err != nil {
					panic(err.Error())
				}

				return c.JSON(Response{
					Code:    200,
					Data:    cpuInfo,
					Message: "success",
				})
			})

			cpuGroup.Get("/percent", func(c *fiber.Ctx) error {
				var params struct {
					DataStart int64 `query:"dateStart"`
					DateEnd   int64 `query:"dateEnd"`
					TimeSize  int64 `query:"timeSize"`
				}
				_ = c.QueryParser(&params)

				if params.DataStart == 0 {
					params.DataStart = time.Now().Unix() - 3600
				}
				if params.DateEnd == 0 {
					params.DateEnd = time.Now().Unix()
				}
				if params.TimeSize == 0 {
					params.TimeSize = 5
				}

				var result []CpuPercentModule
				err := mgo.Query("cpuPercent", &result, params.DataStart, params.DateEnd)
				if err != nil {
					panic(err.Error())
				}

				var chartData []ChartData[float64]
				for _, cpuPercent := range result {
					if cpuPercent.Timestamp%params.TimeSize == 0 {
						chartData = append(chartData, ChartData[float64]{
							Datetime: time.Unix(cpuPercent.Timestamp, 0).Format("2006-01-02 15:04:05"),
							Name:     "CPU",
							Value:    cpuPercent.SingleCpuPercent,
						})
					}
				}

				return c.JSON(Response{
					Code:    200,
					Data:    chartData,
					Message: "success",
				})
			})
		}

		// 内存相关
		memGroup := v1.Group("/mem")
		{
			memGroup.Get("/percent", func(c *fiber.Ctx) error {
				var params struct {
					DataStart int64 `query:"dateStart"`
					DateEnd   int64 `query:"dateEnd"`
					TimeSize  int64 `query:"timeSize"`
				}
				_ = c.QueryParser(&params)

				if params.DataStart == 0 {
					params.DataStart = time.Now().Unix() - 3600
				}
				if params.DateEnd == 0 {
					params.DateEnd = time.Now().Unix()
				}
				if params.TimeSize == 0 {
					params.TimeSize = 5
				}

				var result []MemUsageModule
				err := mgo.Query("memUsage", &result, params.DataStart, params.DateEnd)
				if err != nil {
					panic(err.Error())
				}

				var chartData []ChartData[float64]
				for _, memUsage := range result {
					if memUsage.Timestamp%params.TimeSize == 0 {
						chartData = append(chartData, ChartData[float64]{
							Datetime: time.Unix(memUsage.Timestamp, 0).Format("2006-01-02 15:04:05"),
							Name:     "Memory",
							Value:    memUsage.MemStat.UsedPercent,
						})

						chartData = append(chartData, ChartData[float64]{
							Datetime: time.Unix(memUsage.Timestamp, 0).Format("2006-01-02 15:04:05"),
							Name:     "SWAP",
							Value:    float64(memUsage.MemStat.SwapFree) / float64(memUsage.MemStat.SwapTotal),
						})
					}
				}

				return c.JSON(Response{
					Code:    200,
					Data:    chartData,
					Message: "success",
				})
			})
		}

		// 磁盘相关
		diskGroup := v1.Group("/disk")
		{
			diskGroup.Get("/io", func(c *fiber.Ctx) error {
				var params struct {
					DataStart int64 `query:"dateStart"`
					DateEnd   int64 `query:"dateEnd"`
					TimeSize  int64 `query:"timeSize"`
				}
				_ = c.QueryParser(&params)

				if params.DataStart == 0 {
					params.DataStart = time.Now().Unix() - 3600
				}
				if params.DateEnd == 0 {
					params.DateEnd = time.Now().Unix()
				}
				if params.TimeSize == 0 {
					params.TimeSize = 5
				}

				var result []DiskIOModule
				err := mgo.Query("diskIO", &result, params.DataStart, params.DateEnd)
				if err != nil {
					panic(err.Error())
				}

				var chartData []ChartData[uint64]
				for _, diskIO := range result {
					if diskIO.Timestamp%params.TimeSize == 0 {
						for _, disk := range diskIO.DiskIO {
							chartData = append(chartData, ChartData[uint64]{
								Datetime: time.Unix(diskIO.Timestamp, 0).Format("2006-01-02 15:04:05"),
								Name:     disk.Device + " read",
								Value:    disk.ReadIO,
							})
							chartData = append(chartData, ChartData[uint64]{
								Datetime: time.Unix(diskIO.Timestamp, 0).Format("2006-01-02 15:04:05"),
								Name:     disk.Device + " write",
								Value:    disk.WriteIO,
							})
						}
					}
				}

				return c.JSON(Response{
					Code:    200,
					Data:    chartData,
					Message: "success",
				})
			})

			diskGroup.Get("/usage", func(c *fiber.Ctx) error {
				diskPart, _ := disk.Partitions(false)

				var diskUsages []DiskUsage
				for i := range diskPart {
					diskUsage, _ := disk.Usage(diskPart[i].Mountpoint)
					diskUsages = append(diskUsages, DiskUsage{
						Device:      diskPart[i].Device,
						FsType:      diskPart[i].Fstype,
						MountPoint:  diskPart[i].Mountpoint,
						TotalBytes:  diskUsage.Total,
						FreeBytes:   diskUsage.Free,
						UsedBytes:   diskUsage.Used,
						UsedPercent: diskUsage.UsedPercent,
					})
				}

				return c.JSON(Response{
					Code:    200,
					Data:    diskUsages,
					Message: "success",
				})
			})
		}

		// 网络相关
		networkGroup := v1.Group("/network")
		{
			networkGroup.Get("/load", func(c *fiber.Ctx) error {
				var params struct {
					DataStart int64 `query:"dateStart"`
					DateEnd   int64 `query:"dateEnd"`
					TimeSize  int64 `query:"timeSize"`
				}
				_ = c.QueryParser(&params)

				if params.DataStart == 0 {
					params.DataStart = time.Now().Unix() - 3600
				}
				if params.DateEnd == 0 {
					params.DateEnd = time.Now().Unix()
				}
				if params.TimeSize == 0 {
					params.TimeSize = 5
				}

				var result []NetworkLoadModule
				err := mgo.Query("networkLoad", &result, params.DataStart, params.DateEnd)
				if err != nil {
					panic(err.Error())
				}

				var chartData []ChartData[float64]
				for _, networkLoad := range result {
					if networkLoad.Timestamp%params.TimeSize == 0 {
						for _, network := range networkLoad.NetworkLoad {
							chartData = append(chartData, ChartData[float64]{
								Datetime: time.Unix(networkLoad.Timestamp, 0).Format("2006-01-02 15:04:05"),
								Name:     network.Device + " sent",
								Value:    network.SentBytes,
							})

							chartData = append(chartData, ChartData[float64]{
								Datetime: time.Unix(networkLoad.Timestamp, 0).Format("2006-01-02 15:04:05"),
								Name:     network.Device + " recv",
								Value:    network.RecvBytes,
							})
						}
					}
				}

				return c.JSON(Response{
					Code:    200,
					Data:    chartData,
					Message: "success",
				})
			})
		}

		// 进程相关
		processGroup := v1.Group("process")
		{
			processGroup.Get("/", func(c *fiber.Ctx) error {
				var params struct {
					SortBy string `json:"sortBy"` // cpu  mem  pid user
					Type   string `json:"type"`   //asc or desc
				}
				c.QueryParser(&params)

				/* var sortBy string
				var sortType string

				switch params.SortBy {
				case "cpu":
					sortBy = "pcpu"
				case "mem":
					sortBy = "pmem"
				case "pid":
					sortBy = "pid"
				case "user":
					sortBy = "user"
				default:
					sortBy = "pcpu"
				}

				switch params.Type {
				case "asc":
					sortType = ""
				case "desc":
					sortType = "-"
				default:
					sortType = "-"
				} */

				processStat := readProcessStat()
				formatProcess := formatProcessStat(processStat)

				return c.JSON(Response{
					Code:    200,
					Data:    formatProcess,
					Message: "success",
				})
				/* process, _ := process.Processes()

				var processInfoList []ProcessInfo
				for i := range process {

					if status, _ := process[i].Status(); len(status) >= 1 && status[0] != "" {
						pid := process[i].Pid
						processName, _ := process[i].Name()
						runUser, _ := process[i].Username()
						cpuPercent, _ := process[i].CPUPercent()
						memPercent, _ := process[i].MemoryPercent()
						createTime, _ := process[i].CreateTime()
						command, _ := process[i].Cmdline()
						memInfo, _ := process[i].MemoryInfo()
						processInfoList = append(processInfoList, ProcessInfo{
							PID:         pid,
							ProcessName: processName,
							RunUser:     runUser,
							CpuPercent:  cpuPercent,
							MemPercent:  memPercent,
							RSS:         memInfo.RSS,
							VMS:         memInfo.VMS,
							Status:      status,
							CreateTime:  createTime,
							Command:     command,
						})
					}
				}

				sort.SliceStable(processInfoList, func(i, j int) bool {
					return processInfoList[i].CpuPercent > processInfoList[j].CpuPercent
				}) */

			})
		}
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%d", config.WebPort)))

}
