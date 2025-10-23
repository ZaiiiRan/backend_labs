package settings

type RabbitMqSettings struct {
	Host              string `mapstructure:"Host"`
	Port              int    `mapstructure:"Port"`
	User              string `mapstructure:"User"`
	Password          string `mapstructure:"Password"`
	VHost             string `mapstructure:"VHost"`
}
