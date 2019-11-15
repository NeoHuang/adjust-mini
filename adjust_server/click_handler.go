package adjust_server

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type ClickHandler struct {
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

func NewClickHandler() *ClickHandler {
	return &ClickHandler{}
}

func (handler *ClickHandler) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	tracker := httpRequest.URL.Path[1:]
	if len(tracker) != 7 {
		labels := copyLabels(clickhandlerFailedLabels)
		labels["error"] = "invalid_tracker"
		httpRequestMetrics.With(labels).Inc()
		http.Error(writer, "error: invalid tracker", http.StatusInternalServerError)
		return
	}

	httpRequestMetrics.With(clickhandlerSuccessLabels).Inc()
	fmt.Fprintf(writer, "click tracked")
}

func copyLabels(labels prometheus.Labels) prometheus.Labels {
	if labels == nil {
		return nil
	}

	newLabels := prometheus.Labels{}
	for k, v := range labels {
		newLabels[k] = v
	}

	return newLabels
}
