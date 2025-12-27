package consumer

import (
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type batchMessage struct {
	msg kafka.TopicPartition
	val dalconsumer.Message
}
