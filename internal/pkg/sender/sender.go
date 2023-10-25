//go:generate mockgen -source ./sender.go -destination=./mocks/sender.go -package=mock_sender

package sender

import (
	"time"
)

type RequestMessage struct {
	Time   time.Time `json:"time"`
	Method string    `json:"method"`
	Body   string    `json:"body"`
}

type Sender interface {
	Send(method string, body []byte, sync bool) error
}
