package settings

type RabbitMqSettings struct {
	Host                    string `mapstructure:"Host"`
	Port                    int    `mapstructure:"Port"`
	User                    string `mapstructure:"User"`
	Password                string `mapstructure:"Password"`
	VHost                   string `mapstructure:"VHost"`
	HeartbeatSeconds        int    `mapstructure:"HeartbeatSeconds"`
	MaxReconnectAttempts    int    `mapstructure:"MaxReconnectAttempts"`
	ReconnectTimeoutSeconds int    `mapstructure:"ReconnectTimeoutSeconds"`
}
