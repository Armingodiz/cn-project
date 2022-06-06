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

type SysInfo struct {
	Hostname     string  `json:"hostname"`
	RAM          uint64  `json:"ram"`
	Disk         uint64  `json:"disk"`
	UsedMemory   uint64  `json:"used_memory"`
	CachedMemory uint64  `json:"cached_memory"`
	TotalCpu     float64 `json:"total_cpu"`
}

func main() {
	// starting server to expose metrics
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

	listener, err := createListener("tcp", ":80")
	if err != nil {
		log.Println(err)
		return
	}
	for {
		connection, err := listenningForConnection(listener, ":80")
		if err != nil {
			log.Println(err)
			return
		}
		defer connection.Close()
		metrics, err := readMetrics(connection)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(metrics)
			convertInfoToPrometheusMetrics(metrics)
		}
	}
}

func createListener(network, address string) (listener *net.TCPListener, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	listener, err = net.ListenTCP(network, tcpAddress)
	return
}

func listenningForConnection(listener *net.TCPListener, port string) (conn net.Conn, err error) {
	fmt.Println("listenning on port " + port + " ...")
	conn, err = listener.Accept()
	if err != nil {
		return
	}
	return
}

func readMetrics(connection net.Conn) (metrics SysInfo, err error) {
	buffer := make([]byte, 4096)
	n, err := connection.Read(buffer[0:])
	if err != nil {
		return
	}
	err = json.Unmarshal(buffer[:n], &metrics)
	return
}

func convertInfoToPrometheusMetrics(info SysInfo) {
	Ram.With(prometheus.Labels{"agent": info.Hostname}).Add(float64(info.RAM))
	Disk.With(prometheus.Labels{"agent": info.Hostname}).Add(float64(info.Disk))
	UsedMemory.With(prometheus.Labels{"agent": info.Hostname}).Add(float64(info.UsedMemory))
	CachedMemory.With(prometheus.Labels{"agent": info.Hostname}).Add(float64(info.CachedMemory))
	TotalCpu.With(prometheus.Labels{"agent": info.Hostname}).Add(float64(info.TotalCpu))
}
