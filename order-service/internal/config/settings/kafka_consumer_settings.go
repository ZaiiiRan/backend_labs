package settings

type KafkaConsumerSettings struct {
	KafkaSettings           KafkaSettings `mapstructure:"KafkaSettings"`
	GroupId                 string        `mapstructure:"GroupId"`
	OrderCreatedTopic       string        `mapstructure:"OrderCreatedTopic"`
	OrderStatusChangedTopic string        `mapstructure:"OrderStatusChangedTopic"`
}
