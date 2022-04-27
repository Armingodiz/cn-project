package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	cpuState "github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp", ":80")
	if err != nil {
		log.Println(err.Error())
		return
	}
	connection, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Println(err.Error())
		return
	}
	info, _ := getSystemInfo()
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = connection.Write(bytes)
	if err != nil {
		log.Println(err.Error())
		return
	}
	buffer := make([]byte, 512)
	_, err = connection.Read(buffer[0:])
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(buffer))
}

type SysInfo struct {
	Hostname     string  `json:"hostname"`
	Platform     string  `json:"platform"`
	CPU          string  `json:"cpu"`
	RAM          uint64  `json:"ram"`
	Disk         uint64  `json:"disk"`
	UsedMemory   uint64  `json:"used_memory"`
	CachedMemory uint64  `json:"cached_memory"`
	TotalCpu     float64 `json:"total_cpu"`
}

func getSystemInfo() (*SysInfo, error) {
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")

	info := new(SysInfo)

	info.Hostname = hostStat.Hostname
	info.Platform = hostStat.Platform
	info.CPU = cpuStat[0].ModelName
	info.RAM = vmStat.Total / 1024 / 1024
	info.Disk = diskStat.Total / 1024 / 1024
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	info.UsedMemory = memory.Used
	info.CachedMemory = memory.Cached

	before, err := cpuState.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpuState.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	total := float64(after.Total - before.Total)
	info.TotalCpu = total
	fmt.Printf("%+v\n", info)
	return info, nil
}
