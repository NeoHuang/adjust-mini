package metrics

import "github.com/prometheus/client_golang/prometheus"

func CopyLabels(labels prometheus.Labels) prometheus.Labels {
	if labels == nil {
		return nil
	}

	newLabels := prometheus.Labels{}
	for k, v := range labels {
		newLabels[k] = v
	}

	return newLabels
}
