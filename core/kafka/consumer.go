package kafka

import (
	"context"
	"log"

	"github.com/NeoHuang/adjust-mini/core/metrics"
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

type ClusterConsumer struct {
	service       string
	consumerGroup sarama.ConsumerGroup
	msgProcessor  MessageProcessor

	topic         string
	commitOffsets bool

	ctx         context.Context
	cancelCtxFn context.CancelFunc
}

func NewClusterConsumer(consumerGroup sarama.ConsumerGroup, topic string, service string) *ClusterConsumer {
	ctx, cancelCtxFn := context.WithCancel(context.Background())
	return &ClusterConsumer{
		service:       service,
		consumerGroup: consumerGroup,
		topic:         topic,

		ctx:         ctx,
		cancelCtxFn: cancelCtxFn,
	}
}

func (consumer *ClusterConsumer) StartConsuming(msgProcessor MessageProcessor, commitOffsets bool) {
	consumer.msgProcessor = msgProcessor
	consumer.commitOffsets = commitOffsets
	go func() {
		for err := range consumer.consumerGroup.Errors() {
			consumer.IncMetrics(kafkaConsumerConsumedFailedLabels)
			log.Printf("cluster consumer error:%s", err)
		}
	}()

	// gracefully close a session, preventing consumer group from deadlock if the session didn't
	// receive any claims to consume
	// See https://github.com/Shopify/sarama/issues/1351
	ctx, _ := context.WithCancel(context.Background())

	go func() {
		for {
			//consumer.countM("consumer.session.begin")
			// consume blocks until rebalance happens
			log.Println("kafka consumer rebalanced")
			err := consumer.consumerGroup.Consume(ctx, []string{consumer.topic}, consumer)
			if err != nil {
				log.Printf("consumer group consume error:%s", err)
				//	consumer.countM("consumer.consume-error")
			}
			//consumer.countM("consumer.session.end")
		}

		// consumer.countM("consumer.closed")
	}()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ClusterConsumer) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("kafka consumer setup")
	consumer.IncMetrics(kafkaConsumerSetupLabels)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ClusterConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// ConsumeClaim is called within a goroutine. each partition has its own goroutine
func (consumer *ClusterConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("kafka consumer claim")
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Println("kafka consumer consumed")
		processed := consumer.msgProcessor.Process(consumer.ctx, message)
		consumer.IncMetrics(kafkaConsumerConsumedLabels)
		if processed {
			if consumer.commitOffsets {
				session.MarkMessage(message, "")
			}
		}
	}

	return nil
}

func (consumer *ClusterConsumer) IncMetrics(labels prometheus.Labels) {
	newLabels := metrics.CopyLabels(labels)
	newLabels["service"] = consumer.service
	log.Println("XXXXX ", consumer.topic)
	newLabels["topic"] = consumer.topic
	KafkaMetrics.With(newLabels).Inc()
}
