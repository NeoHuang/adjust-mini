package adjust_server

import (
	"log"
	"net/http"

	"github.com/NeoHuang/adjust-mini/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AdjustServer struct {
}

func New() *AdjustServer {
	server := &AdjustServer{}
	return server
}

func (server *AdjustServer) Start() {
	mux := server.ServeMux()
	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Panicf("ADJUST SERVER SHUTTING DOWN (%s)\n\n", err)
	}
}

func (server *AdjustServer) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/heartbeat", handlers.NewHeartbeatHandler())

	// mux.Handle("/heartbeat.gif", heartbeatHandler)

	return mux
}
