package kafka_test

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"homework-3/infrastructure/kafka"
	"log"
	"strings"
	"sync"
	"testing"
)

type TProducer struct {
	Producer *kafka.Producer
	brokers  []string
	sync.Mutex
}

func NewProducerFromEnv() *TProducer {
	brokersValues := viper.GetString("BROKERS")
	brokers := strings.Split(brokersValues, ",")

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatal(err)
	}

	return &TProducer{
		Producer: kafkaProducer,
		brokers:  brokers,
	}
}

func (p *TProducer) SetUp(t *testing.T, args ...interface{}) {
	t.Helper()
	p.Lock()
	p.Truncate()
}

func (p *TProducer) TearDown() {
	defer p.Unlock()
	p.Producer.Close()
	p.Truncate()
}

func (p *TProducer) Truncate() {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(p.brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer admin.Close()

	topic := viper.GetString("TOPIC")
	partitions, err := client.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	offsets := make(map[int32]int64)
	for _, partition := range partitions {
		offsets[partition] = sarama.OffsetNewest
	}

	err = admin.DeleteRecords(topic, offsets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("All messages from kafka topic - %s is deleted\n", topic)
}
