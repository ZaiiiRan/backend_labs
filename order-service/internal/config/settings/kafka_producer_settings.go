package settings

type KafkaProducerSettings struct {
	KafkaSettings           KafkaSettings `mapstructure:"KafkaSettings"`
	ClientId                string        `mapstructure:"ClientId"`
	OrderCreatedTopic       string        `mapstructure:"OrderCreatedTopic"`
	OrderStatusChangedTopic string        `mapstructure:"OrderStatusChangedTopic"`
}
