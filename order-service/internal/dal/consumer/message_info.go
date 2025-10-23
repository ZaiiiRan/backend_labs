package consumer

import "time"

type MessageInfo struct {
	DeliveryTag uint64
	Body        []byte
	ReceivedAt  time.Time
}
