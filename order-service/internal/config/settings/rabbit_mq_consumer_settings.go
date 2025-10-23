package settings

type RabbitMqConsumerSettings struct {
	RabbitMqSettings    RabbitMqSettings `mapstructure:"RabbitMqSettings"`
	Queue               string           `mapstructure:"Queue"`
	BatchSize           int              `mapstructure:"BatchSize"`
	BatchTimeoutSeconds int              `mapstructure:"BatchTimeoutSeconds"`
	Consumer            string           `mapstructure:"Consumer"`
}
