package sender

import "errors"

var (
	ErrSendSyncMessage  = errors.New("send sync message error")
	ErrSendASyncMessage = errors.New("send async message error")
)
