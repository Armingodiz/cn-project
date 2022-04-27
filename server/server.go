package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

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

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp", ":80")
	if err != nil {
		log.Println(err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for {
		fmt.Println("listenning on port 80 ...")
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer connection.Close()
		_, err = connection.Write([]byte("hi armin"))
		if err != nil {
			log.Println(err.Error())
			return
		}
		buffer := make([]byte, 4096)
		n, err := connection.Read(buffer[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		var metrics SysInfo
		err = json.Unmarshal(buffer[:n], &metrics)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(metrics)
	}
}
