package main

import (
	"context"
	"log"
	"math"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MonitorRun(db *mongo.Database) {
	go WriteCpuPercent(db)
	// go WriteMultiCpuPercent(db)
	go WriteMemUsage(db)
	go WriteDiskIO(db)
	go WriteNetworkLoad(db)
}

// CPU使用率
func WriteCpuPercent(db *mongo.Database) {
	collection := db.Collection("cpuPercent")
	for {
		if time.Now().Unix()%5 == 0 {
			singleCpuPercent, err := cpu.Percent(time.Second*5, false)
			if err != nil {
				log.Fatal(err.Error())
			}
			collection.InsertOne(context.TODO(), CpuPercentModule{
				ID:               primitive.NewObjectID(),
				Timestamp:        time.Now().Unix(),
				SingleCpuPercent: singleCpuPercent[0],
			})
		}
	}
}

// CPU核心使用率 已弃用
func WriteMultiCpuPercent(db *mongo.Database) {
	collection := db.Collection("multiCpuPercent")
	for {
		if time.Now().Unix()%5 == 0 {
			multiCpuPercent, err := cpu.Percent(time.Second*5, true)
			if err != nil {
				log.Fatal(err.Error())
			}
			collection.InsertOne(context.TODO(), MultiCpuPercentModule{
				ID:              primitive.NewObjectID(),
				Timestamp:       time.Now().Unix(),
				MultiCpuPercent: multiCpuPercent,
			})
		}
	}
}

// 内存使用率
func WriteMemUsage(db *mongo.Database) {
	collection := db.Collection("memUsage")
	for {
		if time.Now().Unix()%5 == 0 {
			memUsage, err := mem.VirtualMemory()
			if err != nil {
				log.Fatal(err.Error())
			}

			collection.InsertOne(context.TODO(), MemUsageModule{
				ID:        primitive.NewObjectID(),
				MemStat:   memUsage,
				Timestamp: time.Now().Unix(),
			})
			time.Sleep(time.Second * 5)
		}

	}
}

// 磁盘IO
func WriteDiskIO(db *mongo.Database) {
	collection := db.Collection("diskIO")
	for {
		if time.Now().Unix()%5 == 0 {
			devices, _ := disk.Partitions(false)

			startIOCounters, _ := disk.IOCounters()
			time.Sleep(time.Second * 5)
			endIOCounters, _ := disk.IOCounters()

			var diskIO []DiskIO
			for i := range devices {
				arr := strings.Split(devices[i].Device, "/")
				device := strings.Split(devices[i].Device, "/")[len(arr)-1]

				readIO := (endIOCounters[device].ReadBytes - startIOCounters[device].ReadBytes) / 5    // bytes
				writeIO := (endIOCounters[device].WriteBytes - startIOCounters[device].WriteBytes) / 5 // bytes
				ioWait := float64(endIOCounters[device].WriteTime-startIOCounters[device].WriteTime) / float64(endIOCounters[device].WriteCount-startIOCounters[device].WriteCount)

				if math.IsNaN(ioWait) {
					ioWait = 0
				}

				diskIO = append(diskIO, DiskIO{
					Device:  device,
					ReadIO:  readIO,
					WriteIO: writeIO,
					IOWait:  ioWait,
				})
			}

			collection.InsertOne(context.TODO(), DiskIOModule{
				ID:        primitive.NewObjectID(),
				Timestamp: time.Now().Unix(),
				DiskIO:    diskIO,
			})
		}

		/*
				浪费一晚上写的磁盘io算法  基于/proc/stats 能用！！！
			formatDiskArr1 := formatDiskStats(readDiskStats())
			startIOCounters, _ := disk.IOCounters()
			time.Sleep(time.Second * 5)
			formatDiskArr2 := formatDiskStats(readDiskStats())
			endIOCounters, _ := disk.IOCounters()

			var diskStatsDetail []DiskStats
			for i := range formatDiskArr1 {
				if formatDiskArr1[i][1] == "0" {
					// read
					rdSectors1, _ := strconv.ParseFloat(formatDiskArr1[i][5], 64)
					rdSectors2, _ := strconv.ParseFloat(formatDiskArr2[i][5], 64)
					totalRead := (rdSectors2 - rdSectors1) / 5.00 / 2.00

					// write
					wrSectors1, _ := strconv.ParseFloat(formatDiskArr1[i][9], 64)
					wrSectors2, _ := strconv.ParseFloat(formatDiskArr2[i][9], 64)
					totalWrite := (wrSectors2 - wrSectors1) / 5.00 / 2.00

					// await
					rdTicks1, _ := strconv.ParseFloat(formatDiskArr1[i][6], 64)
					rdTicks2, _ := strconv.ParseFloat(formatDiskArr2[i][6], 64)

					wrTicks1, _ := strconv.ParseFloat(formatDiskArr1[i][10], 64)
					wrTicks2, _ := strconv.ParseFloat(formatDiskArr2[i][10], 64)

					rdIos1, _ := strconv.ParseFloat(formatDiskArr1[i][3], 64)
					rdIos2, _ := strconv.ParseFloat(formatDiskArr2[i][3], 64)

					wrIos1, _ := strconv.ParseFloat(formatDiskArr1[i][7], 64)
					wrIos2, _ := strconv.ParseFloat(formatDiskArr2[i][7], 64)

					totalAwait := ((rdTicks2 - rdTicks1) + (wrTicks2 - wrTicks1)) / ((rdIos2 - rdIos1) + (wrIos2 - wrIos1))

					diskStatsDetail = append(diskStatsDetail, DiskStats{
						Device:  formatDiskArr1[i][2],
						ReadIO:  totalRead,
						WriteIO: totalWrite,
						Await:   totalAwait,
					})
				}
			}
			fmt.Println(diskStatsDetail)
			fmt.Println(float64(endIOCounters[device].WriteTime-startIOCounters[device].WriteTime) / float64(endIOCounters[device].WriteCount-startIOCounters[device].WriteCount))
			// fmt.Println(endIOCounters[device].WriteCount - startIOCounters[device].WriteCount)
		*/
	}
}

// 网络负载
func WriteNetworkLoad(db *mongo.Database) {
	collection := db.Collection("networkLoad")
	for {
		if time.Now().Unix()%5 == 0 {
			startNetIOCounters, _ := net.IOCounters(true)
			time.Sleep(time.Second * 5)
			endNetIOCounters, _ := net.IOCounters(true)

			var networkLoadArr []NetworkLoad
			for i := range startNetIOCounters {
				sentBytes := float64(endNetIOCounters[i].BytesSent-startNetIOCounters[i].BytesSent) / 5
				recvBytes := float64(endNetIOCounters[i].BytesRecv-startNetIOCounters[i].BytesRecv) / 5

				sentPacket := float64(endNetIOCounters[i].PacketsSent-startNetIOCounters[i].PacketsSent) / 5
				recvPacket := float64(endNetIOCounters[i].PacketsRecv-startNetIOCounters[i].PacketsRecv) / 5

				networkLoadArr = append(networkLoadArr, NetworkLoad{
					Device:     startNetIOCounters[i].Name,
					SentBytes:  sentBytes,
					RecvBytes:  recvBytes,
					SentPacket: sentPacket,
					RecvPacket: recvPacket,
				})
			}
			collection.InsertOne(context.TODO(), NetworkLoadModule{
				ID:          primitive.NewObjectID(),
				Timestamp:   time.Now().Unix(),
				NetworkLoad: networkLoadArr,
			})
		}

	}
}
