package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework-3/infrastructure/kafka"
	"homework-3/internal/pkg/sender"
)

type Receiver interface {
	Subscribe(topic string) error
}

type Service struct {
	receiver Receiver
}

func (s *Service) Close() {
	<-context.TODO().Done()
}

func NewConsumerService(brokers []string) (*Service, error) {
	kafkaConsumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		return nil, err
	}

	messageHandler := map[string]HandleFunc{
		"requests": func(message *sarama.ConsumerMessage) {
			rm := sender.RequestMessage{}
			err = json.Unmarshal(message.Value, &rm)
			if err != nil {
				fmt.Println("Consumer error", err)
			}
			fmt.Println("Received Key: ", string(message.Key), " Value: ", rm)
		},
	}

	requests := NewService(
		NewReceiver(kafkaConsumer, messageHandler),
	)
	return requests, nil
}

func NewService(receiver Receiver) *Service {
	return &Service{
		receiver: receiver,
	}
}

func (s *Service) StartConsume(topic string) {
	err := s.receiver.Subscribe(topic)

	if err != nil {
		fmt.Println("Subscribe error ", err)
	}
}
