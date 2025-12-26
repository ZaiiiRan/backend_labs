package settings

type KafkaSettings struct {
	BootstrapServers string `mapstructure:"BootstrapServers"`
	ClientId         string `mapstructure:"ClientId"`
}
