package settings

type RabbitMqDeadLetterSettings struct {
	Dlx        string `mapstructure:"Dlx"`
	Dlq        string `mapstructure:"Dlq"`
	RoutingKey string `mapstructure:"RoutingKey"`
}
