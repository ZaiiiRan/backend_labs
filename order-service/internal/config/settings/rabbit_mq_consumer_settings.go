package settings

type RabbitMqConsumerSettings struct {
	RabbitMqSettings    RabbitMqSettings `mapstructure:"RabbitMqSettings"`
	Consumer            string           `mapstructure:"Consumer"`
	Queue               string           `mapstructure:"Queue"`
	BatchSize           int              `mapstructure:"BatchSize"`
	BatchTimeoutSeconds int              `mapstructure:"BatchTimeoutSeconds"`
}
