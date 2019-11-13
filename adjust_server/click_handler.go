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
	}
)

func NewClickHandler() *ClickHandler {
	return &ClickHandler{}
}

func (handler *ClickHandler) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	httpRequestMetrics.With(clickhandlerSuccessLabels).Inc()
	fmt.Fprintf(writer, "click tracked")
}
