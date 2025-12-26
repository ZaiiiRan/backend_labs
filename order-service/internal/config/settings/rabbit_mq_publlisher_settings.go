package settings

type RabbitMqPublisherSettings struct {
	RabbitMqSettings RabbitMqSettings          `mapstructure:"RabbitMqSettings"`
	Exchange         string                    `mapstructure:"Exchange"`
	ExchangeMappings []RabbitMqExchangeMapping `mapstructure:"ExchangeMappings"`
}

type RabbitMqExchangeMapping struct {
	Queue             string `mapstructure:"Queue"`
	RoutingKeyPattern string `mapstructure:"RoutingKeyPattern"`
}
