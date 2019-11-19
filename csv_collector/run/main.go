package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct{}

func main() {
	podName := os.Getenv("MY_POD_NAME")
	log.Println("MY_POD_NAME:", podName)
	fileName := "/collected/" + podName + ".json"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicf("Failed to open file %q, err:%s", fileName, err)
	}
	defer f.Close()
	now := time.Now().Format("2006-01-02T15:04:05.000Z") + "\n"
	f.Write([]byte("Started at " + now))
	server := &MetricsServer{}
	server.Start()
}

func (server *MetricsServer) Start() {
	mux := server.ServeMux()

	log.Printf("csv_collector server started....")

	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Panicf("CSV COLLECTOR SHUTTING DOWN (%s)\n\n", err)
	}
}

func (server *MetricsServer) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
