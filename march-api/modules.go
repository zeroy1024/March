package main

import (
	"github.com/shirou/gopsutil/v3/mem"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CpuPercentModule struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp        int64              `json:"timestamp" bson:"timestamp"`
	SingleCpuPercent float64            `json:"singleCpuPercent" bson:"singleCpuPercent"`
}

type SingleCpuPercentModule struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp        int64              `json:"timestamp" bson:"timestamp"`
	SingleCpuPercent float64            `json:"singleCpuPercent" bson:"singleCpuPercent"`
}

type MultiCpuPercentModule struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp       int64              `json:"timestamp" bson:"timestamp"`
	MultiCpuPercent []float64          `json:"multiCpuPercent" bson:"multiCpuPercent"`
}

type MemUsageModule struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id"`
	MemStat   *mem.VirtualMemoryStat `json:"memStat" bson:"memStat"`
	Timestamp int64                  `json:"timestamp" bson:"timestamp"`
}

type DiskIOModule struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
	DiskIO    []DiskIO           `json:"diskIO" bson:"diskIO"`
}

type NetworkLoadModule struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp   int64              `json:"timestamp" bson:"timestamp"`
	NetworkLoad []NetworkLoad      `json:"networkLoad" bson:"networkLoad"`
}
