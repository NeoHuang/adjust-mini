package adjust_server

import (
	"log"
	"net/http"

	"github.com/NeoHuang/adjust-mini/core/kafka"
	"github.com/NeoHuang/adjust-mini/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "backend",
		Subsystem: "adjust_server",
		Name:      "http",
		Help:      "http request metrics",
	}, []string{
		"activity",
		"result",
		"error",
	})
)

type AdjustServer struct {
	Version   string
	panicChan chan struct{}
}

func New(version string) *AdjustServer {
	server := &AdjustServer{
		Version:   version,
		panicChan: make(chan struct{}),
	}
	return server
}

func (server *AdjustServer) Start() {
	mux := server.ServeMux()
	go func() {
		<-server.panicChan
		log.Panicf("don't panic.....")
	}()

	prometheus.MustRegister(httpRequestMetrics)
	prometheus.MustRegister(kafka.KafkaMetrics)

	log.Printf("Adjust server started....")

	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Panicf("ADJUST SERVER SHUTTING DOWN (%s)\n\n", err)
	}
}

func (server *AdjustServer) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/heartbeat", handlers.NewHeartbeatHandler())
	mux.Handle("/version", handlers.NewVersionHandler(server.Version))
	mux.HandleFunc("/panic", func(http.ResponseWriter, *http.Request) {
		close(server.panicChan)
	})
	mux.Handle("/", NewClickHandler())

	return mux
}
