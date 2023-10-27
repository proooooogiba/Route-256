package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework-3/infrastructure/kafka"
	"homework-3/internal/pkg/sender"
	"strconv"
	"time"
)

type KafkaSender struct {
	producer *kafka.Producer
	topic    string
	msgID    int64
}

func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
		0,
	}
}
func (s *KafkaSender) Send(method string, body []byte, sync bool) error {
	reqMsg := sender.RequestMessage{
		Time:   time.Now(),
		Method: method,
		Body:   string(body),
	}
	switch sync {
	case true:
		err := s.sendMessage(reqMsg)
		if err != nil {
			return sender.ErrSendSyncMessage
		}
	case false:
		err := s.sendAsyncMessage(reqMsg)
		if err != nil {
			return sender.ErrSendASyncMessage
		}
	}
	return nil
}

func (s *KafkaSender) sendAsyncMessage(message sender.RequestMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	s.producer.SendAsyncMessage(kafkaMsg)

	fmt.Println("Send async message with key:", kafkaMsg.Key)

	return nil
}

func (s *KafkaSender) sendMessage(message sender.RequestMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	_, _, err = s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		fmt.Println("Send message connector error", err)
		return err
	}
	return nil
}

func (s *KafkaSender) sendMessages(messages []sender.RequestMessage) error {
	var kafkaMsg []*sarama.ProducerMessage
	var message *sarama.ProducerMessage
	var err error

	for _, m := range messages {
		message, err = s.buildMessage(m)
		kafkaMsg = append(kafkaMsg, message)

		if err != nil {
			fmt.Println("Send message marshal error", err)
			return err
		}
	}

	err = s.producer.SendSyncMessages(kafkaMsg)

	if err != nil {
		fmt.Println("Send message connector error", err)
		return err
	}

	fmt.Println("Send messages count:", len(messages))
	return nil
}

func (s *KafkaSender) buildMessage(message sender.RequestMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		fmt.Println("Send message marshal error", err)
		return nil, err
	}

	s.msgID++

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(strconv.FormatInt(s.msgID, 10)),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}
