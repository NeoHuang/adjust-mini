package kafka

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "backend",
		Subsystem: "core",
		Name:      "kafka",
		Help:      "kafka package metrics",
	}, []string{
		"service",
		"component",
		"action",
		"topic",
		"result",
	})

	kafkaProducerCalledLabels = prometheus.Labels{
		"component": "producer",
		"action":    "called",
		"result":    "successful",
	}

	kafkaProducerSuccessLabels = prometheus.Labels{
		"component": "producer",
		"action":    "produce",
		"result":    "successful",
	}
	kafkaProducerFailedLabels = prometheus.Labels{
		"component": "producer",
		"action":    "produce",
		"result":    "failed",
	}

	kafkaConsumerSetupLabels = prometheus.Labels{
		"component": "consumer",
		"action":    "setup",
		"result":    "successful",
	}
	kafkaConsumerConsumedLabels = prometheus.Labels{
		"component": "consumer",
		"action":    "consumed",
		"result":    "successful",
	}
	kafkaConsumerConsumedFailedLabels = prometheus.Labels{
		"component": "consumer",
		"action":    "consumed",
		"result":    "failed",
	}
)
