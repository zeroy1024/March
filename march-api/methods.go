package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func formatBytes(bytes float64) string {
	if bytes > 1024 && bytes <= 1024*1024 {
		return fmt.Sprintf("%.2fKB", bytes/1024)
		// return strconv.Itoa(int(bytes/1024)) + "KB"
	} else if bytes > 1024*1024 && bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", bytes/1024/1024)
		// return strconv.Itoa(int(bytes/1024/1024)) + "MB"
	} else if bytes > 1024*1024*1024 && bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", bytes/1024/1024/1024)
		// return strconv.Itoa(int(bytes/1024/1024/1024)) + "GB"
	} else {
		return fmt.Sprintf("%.0fB", bytes)
		// return strconv.Itoa(int(bytes)) + "B"
	}
}

func compressString(str string) string {
	return strings.TrimLeft(regexp.MustCompile("\\s+").ReplaceAllString(str, " "), " ")
}

func readProcessStat() string {
	execResult := exec.Command("ps", "-aux")
	output, _ := execResult.Output()
	return string(output)
}

func formatProcessStat(processStat string) []ProcessInfo {
	lines := strings.Split(processStat, "\n")
	lines = lines[1:]
	var formatProcessInfo []ProcessInfo
	for i := range lines {
		if len(lines[i]) >= 10 {
			precessArr := strings.Split(compressString(lines[i]), " ")

			pid, _ := strconv.Atoi(precessArr[1])
			cpuPercent, _ := strconv.ParseFloat(precessArr[2], 64)
			memPercent, _ := strconv.ParseFloat(precessArr[3], 64)
			formatProcessInfo = append(formatProcessInfo, ProcessInfo{
				PID:        uint64(pid),
				RunUser:    precessArr[0],
				CpuPercent: cpuPercent,
				MemPercent: memPercent,
				Command:    strings.Join(precessArr[10:], " "),
			})
		}
	}
	return formatProcessInfo
}

func readDiskStats() string {
	execResult := exec.Command("cat", "/proc/diskstats")
	output, _ := execResult.Output()
	return string(output)
}

func formatDiskStats(diskStatsStr string) [][]string {
	diskStatsStr = strings.TrimRight(diskStatsStr, "\n")
	diskStatsArr := strings.Split(diskStatsStr, "\n")

	var formatDiskArr [][]string
	for i := range diskStatsArr {
		formatDiskArr = append(formatDiskArr, strings.Split(compressString(diskStatsArr[i]), " "))
	}

	return formatDiskArr
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
