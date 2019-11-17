package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type SaramaAsyncProducer struct {
	rawProducer sarama.AsyncProducer
}

func NewSaramaAsyncProducer(rawProducer sarama.AsyncProducer) *SaramaAsyncProducer {
	producer := &SaramaAsyncProducer{
		rawProducer: rawProducer,
	}
	producer.consumeResponses()
	return producer
}

func NewDefaultSaramaAsyncProducer(hosts []string, config *sarama.Config) (producer *SaramaAsyncProducer, err error) {
	rawProducer, err := sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		log.Printf("cannot connect to kafka (%s)", err)
		return nil, err
	}
	log.Printf("kafka async producer connected (hosts:%s timeout:%s retries:%d)", hosts, config.Net.DialTimeout, config.Producer.Retry.Max)

	return NewSaramaAsyncProducer(rawProducer), nil
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
}

func (producer *SaramaAsyncProducer) consumeResponses() {
	if producer.rawProducer == nil {
		return
	}
}
