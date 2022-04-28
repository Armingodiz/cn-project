package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Ram = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "RAM",
		Help: "ram of agent",
	}, []string{"agent"})
	Disk = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "DISK",
		Help: "disk of agent",
	}, []string{"agent"})
	UsedMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "USED_MEMORY",
		Help: "used memory of agent",
	}, []string{"agent"})
	CachedMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "CACHED_MEMORY",
		Help: "cached memory of agent",
	}, []string{"agent"})
	TotalCpu = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "TOTAL_CPU",
		Help: "total cpu of agent",
	}, []string{"agent"})
)

func init() {
	prometheus.MustRegister(Ram)
	prometheus.MustRegister(Disk)
	prometheus.MustRegister(UsedMemory)
	prometheus.MustRegister(CachedMemory)
	prometheus.MustRegister(TotalCpu)
}

func recordMetrics() {
}

type SysInfo struct {
	Hostname     string  `json:"hostname"`
	RAM          uint64  `json:"ram"`
	Disk         uint64  `json:"disk"`
	UsedMemory   uint64  `json:"used_memory"`
	CachedMemory uint64  `json:"cached_memory"`
	TotalCpu     float64 `json:"total_cpu"`
}

func main() {
	recordMetrics()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("Listening port: %s \n", httpServer.Addr)
		log.Println(httpServer.ListenAndServe())
	}()

	/////////////

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
		_, err = connection.Write([]byte("hi agent"))
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
		Ram.With(prometheus.Labels{"agent": metrics.Hostname}).Add(float64(metrics.RAM))
		Disk.With(prometheus.Labels{"agent": metrics.Hostname}).Add(float64(metrics.Disk))
		UsedMemory.With(prometheus.Labels{"agent": metrics.Hostname}).Add(float64(metrics.UsedMemory))
		CachedMemory.With(prometheus.Labels{"agent": metrics.Hostname}).Add(float64(metrics.CachedMemory))
		TotalCpu.With(prometheus.Labels{"agent": metrics.Hostname}).Add(float64(metrics.TotalCpu))
	}
}
