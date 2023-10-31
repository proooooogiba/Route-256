package consumer

import (
	"errors"
	"github.com/IBM/sarama"
	"homework-3/infrastructure/kafka"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.Consumer
	handlers map[string]HandleFunc
}

func NewReceiver(consumer *kafka.Consumer, handlers map[string]HandleFunc) *KafkaReceiver {
	return &KafkaReceiver{
		consumer: consumer,
		handlers: handlers,
	}
}

func (r *KafkaReceiver) Subscribe(topic string, wantMessages int, messagesChan chan bool) error {
	handler, ok := r.handlers[topic]

	if !ok {
		return errors.New("can not find handler")
	}

	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)

	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		if wantMessages == -1 {
			go func(pc sarama.PartitionConsumer, partition int32) {
				for message := range pc.Messages() {
					handler(message)
				}
			}(pc, partition)
		} else {
			go func(pc sarama.PartitionConsumer, partition int32) {
				var count int
				for message := range pc.Messages() {
					if count == wantMessages {
						break
					}
					handler(message)
					messagesChan <- true
					count++
				}
			}(pc, partition)
		}
	}

	return nil
}

func (r *KafkaReceiver) Close() error {
	err := r.consumer.SingleConsumer.Close()
	if err != nil {
		return err
	}
	return nil
}
