package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework-3/infrastructure/kafka"
	"homework-3/internal/pkg/sender"
)

type Receiver interface {
	Subscribe(topic string, wantMessages int, messagesChan chan bool) error
	Close() error
}

type Service struct {
	receiver Receiver
}

func (s *Service) Close() {
	s.receiver.Close()
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

func (s *Service) StartConsume(topic string, wantMessages int, messagesChan chan bool) {
	err := s.receiver.Subscribe(topic, wantMessages, messagesChan)

	if err != nil {
		fmt.Println("Subscribe error ", err)
	}
}
