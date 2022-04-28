package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "armin",
		Help: "name of user",
	}, []string{"agent"})
)

func init() {
	prometheus.MustRegister(tss)
	//	prometheus.MustRegister(opsProcessed)
}

func recordMetrics() {
	tss.With(prometheus.Labels{"system_name": "mySystem"}).Add(118)
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
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

func main() {

	recordMetrics()
	/*Push(&PushConfig{
		Instance: "server",
		URL:      "http://prom-pushgateway:9091",
		Job:      "metrics",
	})*/
	///////////////////////
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

type PushConfig struct {
	Instance string
	URL      string
	Job      string
}
