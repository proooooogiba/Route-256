package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework-3/infrastructure/kafka"
	"strconv"
	"time"
)

type RequestMessage struct {
	Time   time.Time `json:"time"`
	Method string    `json:"method"`
	Body   string    `json:"body"`
}

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

func (s *KafkaSender) sendAsyncMessage(message RequestMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	s.producer.SendAsyncMessage(kafkaMsg)

	fmt.Println("Send async message with key:", kafkaMsg.Key)

	return nil
}

func (s *KafkaSender) sendMessage(message RequestMessage) error {
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

func (s *KafkaSender) sendMessages(messages []RequestMessage) error {
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

func (s *KafkaSender) buildMessage(message RequestMessage) (*sarama.ProducerMessage, error) {
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
		Headers: []sarama.RecordHeader{ // например, в хедер можно записать версию релиза
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}
