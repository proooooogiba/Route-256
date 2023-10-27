package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers       []string
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

func newAsyncProducer(brokers []string) (sarama.AsyncProducer, error) {
	asyncProducerConfig := sarama.NewConfig()

	asyncProducerConfig.Producer.Partitioner = sarama.NewHashPartitioner
	asyncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	asyncProducerConfig.Producer.Return.Successes = true
	asyncProducerConfig.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(brokers, asyncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with async producer-producer")
	}

	go func() {
		// Error и Retry топики можно использовать при получении ошибки
		for e := range asyncProducer.Errors() {
			fmt.Println(e.Error())
		}
	}()

	go func() {
		for m := range asyncProducer.Successes() {
			fmt.Println("Async success with key", m.Key)
		}
	}()

	return asyncProducer, nil
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync producer-producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync producer-producer")
	}

	asyncProducer, err := newAsyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with async producer-producer")
	}

	producer := &Producer{
		brokers:       brokers,
		syncProducer:  syncProducer,
		asyncProducer: asyncProducer,
	}

	return producer, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	err := k.syncProducer.SendMessages(messages)
	if err != nil {
		fmt.Println("producer.Connector.SendMessages error", err)
	}

	return err
}

func (k *Producer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.asyncProducer.Input() <- message
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "producer.Connector.Close")
	}

	err = k.asyncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "producer.Connector.Close")
	}

	return nil
}
