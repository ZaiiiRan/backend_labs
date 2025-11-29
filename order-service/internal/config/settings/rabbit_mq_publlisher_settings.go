package settings

type RabbitMqPublisherSettings struct {
	RabbitMqSettings RabbitMqSettings `mapstructure:"RabbitMqSettings"`
	Queue            string           `mapstructure:"Queue"`
}
