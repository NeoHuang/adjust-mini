package kafka

import (
	"log"

	"github.com/NeoHuang/adjust-mini/core/metrics"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	KafkaMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "backend",
		Subsystem: "core",
		Name:      "kafka",
		Help:      "kafka package metrics",
	}, []string{
		"service",
		"action",
		"topic",
		"result",
	})

	kafkaProducerCalledLabels = prometheus.Labels{
		"action": "called",
		"result": "successful",
	}

	kafkaProducerSuccessLabels = prometheus.Labels{
		"action": "produce",
		"result": "successful",
	}
	kafkaProducerFailedLabels = prometheus.Labels{
		"action": "produce",
		"result": "failed",
	}
)

type SaramaAsyncProducer struct {
	service     string
	rawProducer sarama.AsyncProducer
}

func NewSaramaAsyncProducer(rawProducer sarama.AsyncProducer, service string) *SaramaAsyncProducer {
	producer := &SaramaAsyncProducer{
		service:     service,
		rawProducer: rawProducer,
	}
	producer.consumeResponses()
	return producer
}

func NewDefaultSaramaAsyncProducer(hosts []string, config *sarama.Config, service string) (producer *SaramaAsyncProducer, err error) {
	rawProducer, err := sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		log.Printf("cannot connect to kafka (%s)", err)
		return nil, err
	}
	log.Printf("kafka async producer connected (hosts:%s timeout:%s retries:%d)", hosts, config.Net.DialTimeout, config.Producer.Retry.Max)

	return NewSaramaAsyncProducer(rawProducer, service), nil
}

func (producer *SaramaAsyncProducer) Produce(topic string, key []byte, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}

	producer.rawProducer.Input() <- msg
	successLabel := metrics.CopyLabels(kafkaProducerCalledLabels)
	successLabel["service"] = producer.service
	successLabel["topic"] = topic
	KafkaMetrics.With(successLabel).Inc()
	log.Printf("msg produce to topic %s", topic)
}

func (producer *SaramaAsyncProducer) consumeResponses() {
	if producer.rawProducer == nil {
		return
	}

	go func() {
		for success := range producer.rawProducer.Successes() {
			successLabel := metrics.CopyLabels(kafkaProducerSuccessLabels)
			successLabel["service"] = producer.service
			successLabel["topic"] = success.Topic
			KafkaMetrics.With(successLabel).Inc()
		}
	}()

	go func() {
		for err := range producer.rawProducer.Errors() {
			failedLabel := metrics.CopyLabels(kafkaProducerFailedLabels)
			failedLabel["service"] = producer.service
			failedLabel["topic"] = err.Msg.Topic
			KafkaMetrics.With(failedLabel).Inc()
			log.Printf("failed to produce to kafka:%s", err)
		}
	}()
}
