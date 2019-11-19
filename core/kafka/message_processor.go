package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type MessageProcessor interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) bool
}
