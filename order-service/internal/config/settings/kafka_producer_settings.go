package settings

type KafkaProducerSettings struct {
	KafkaSettings           KafkaSettings `mapstructure:"KafkaSettings"`
	OrderCreatedTopic       string        `mapstructure:"OrderCreatedTopic"`
	OrderStatusChangedTopic string        `mapstructure:"OrderStatusChangedTopic"`
}
