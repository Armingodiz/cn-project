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
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	var conn *net.TCPConn
	var err error
	for {
		conn, err = connectToServer("tcp", "server:80")
		if err != nil {
			log.Println(err)
		} else {
			err = sendMetrics(conn)
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Second * 10)
	}

}

// try to make a tcp connection and returns the connection or error
func connectToServer(network, address string) (connection *net.TCPConn, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	connection, err = net.DialTCP(network, nil, tcpAddress)
	return
}

// send system info as metrics through tcp connection every timeInterval second
func sendMetrics(conn *net.TCPConn) error {
	info, _ := getSystemInfo()
	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

type SysInfo struct {
	Hostname     string  `json:"hostname"`
	RAM          uint64  `json:"ram"`
	Disk         uint64  `json:"disk"`
	UsedMemory   uint64  `json:"used_memory"`
	CachedMemory uint64  `json:"cached_memory"`
	TotalCpu     float64 `json:"total_cpu"`
}

func getSystemInfo() (*SysInfo, error) {
	hostStat, _ := host.Info()
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")
	info := new(SysInfo)
	info.Hostname = hostStat.Hostname
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
