package adjust_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NeoHuang/adjust-mini/core"
	"github.com/NeoHuang/adjust-mini/core/kafka"
	"github.com/NeoHuang/adjust-mini/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type ClickHandler struct {
	kafkaProducer *kafka.SaramaAsyncProducer
}

var (
	clickhandlerSuccessLabels = prometheus.Labels{
		"activity": "click",
		"result":   "successful",
		"error":    "",
	}

	clickhandlerFailedLabels = prometheus.Labels{
		"activity": "click",
		"result":   "failed",
		"error":    "unknown",
	}
)

func NewClickHandler(producer *kafka.SaramaAsyncProducer) *ClickHandler {
	return &ClickHandler{
		kafkaProducer: producer,
	}
}

func (handler *ClickHandler) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	tracker := httpRequest.URL.Path[1:]
	if len(tracker) != 7 {
		labels := metrics.CopyLabels(clickhandlerFailedLabels)
		labels["error"] = "invalid_tracker"
		httpRequestMetrics.With(labels).Inc()
		http.Error(writer, "error: invalid tracker", http.StatusBadRequest)
		return
	}

	click := core.NewClick(tracker)
	clickJson, err := json.Marshal(click)
	if err != nil {
		labels := metrics.CopyLabels(clickhandlerFailedLabels)
		labels["error"] = "unmarshal"
		httpRequestMetrics.With(labels).Inc()
		http.Error(writer, "error: unmarshal error", http.StatusInternalServerError)
	}

	handler.kafkaProducer.Produce("adjust-mini-clicks", []byte(""), clickJson)
	httpRequestMetrics.With(clickhandlerSuccessLabels).Inc()
	fmt.Fprintf(writer, "click tracked")
}
