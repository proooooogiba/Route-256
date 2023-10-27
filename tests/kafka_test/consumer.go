package kafka_test

import (
	"github.com/spf13/viper"
	"homework-3/internal/consumer"
	"log"
	"strings"
	"sync"
	"testing"
)

type TConsumer struct {
	Consumer *consumer.Service
	sync.Mutex
}

func NewConsumerFromEnv() *TConsumer {
	brokersValues := viper.GetString("BROKERS")
	brokers := strings.Split(brokersValues, ",")

	kafkaConsumer, err := consumer.NewConsumerService(brokers)
	if err != nil {
		log.Fatal(err)
	}

	return &TConsumer{
		Consumer: kafkaConsumer,
	}
}

func (c *TConsumer) SetUp(t *testing.T, args ...interface{}) {
	t.Helper()
	c.Lock()
}

func (c *TConsumer) TearDown() {
	defer c.Unlock()
	c.Consumer.Close()
}
