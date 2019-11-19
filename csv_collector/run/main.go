package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NeoHuang/adjust-mini/core/kafka"
	"github.com/NeoHuang/adjust-mini/csv_collector"
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
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
	prometheus.MustRegister(kafka.KafkaMetrics)

	collector := csv_collector.NewCollector(f)

	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V0_10_2_0
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup([]string{"192.168.2.1:9092"}, "csv-collector-group", consumerConfig)

	consumer := kafka.NewClusterConsumer(consumerGroup, "adjust-mini-clicks", "csv-collector")
	consumer.StartConsuming(collector, true)
	defer f.Close()
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
