package csv_collector

import (
	"context"
	"os"
	"sync"

	"github.com/Shopify/sarama"
)

type Collector struct {
	f     *os.File
	mutex sync.Mutex
}

func NewCollector(file *os.File) *Collector {
	return &Collector{
		f: file,
	}
}
func (collector *Collector) Process(ctx context.Context, msg *sarama.ConsumerMessage) bool {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()
	collector.f.Write(msg.Value)
	collector.f.Write([]byte("\n"))
	return true
}
