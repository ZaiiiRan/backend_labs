package config

import (
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/spf13/viper"
)

type ConsumerConfig struct {
	OrderCreatedRabbitMqConsumerSettings settings.RabbitMqConsumerSettings `mapstructure:"OrderCreatedConsumerSettings"`
	OmsClientGrpcSettings                settings.GrpcClientSettings       `mapstructure:"OmsGrpcClient"`
}

func LoadConsumerConfig() (*ConsumerConfig, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	env := v.GetString("APP_ENVIRONMENT")
	if env == "" {
		env = "Development"
	}

	setDefaultConsumerConfigValues(v)

	v.SetConfigType("yaml")
	v.SetConfigName("appsettings." + env)
	v.AddConfigPath("/etc/order-service-consumer")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg ConsumerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func setDefaultConsumerConfigValues(v *viper.Viper) {
	v.SetDefault("OrderCreatedConsumerSettings.RabbitMqSettings.HeartbeatSeconds", 30)
	v.SetDefault("OrderCreatedConsumerSettings.BatchSize", 100)
	v.SetDefault("OrderCreatedConsumerSettings.BatchTimeoutSeconds", 1)
	v.SetDefault("OrderCreatedConsumerSettings.RabbitMqSettings.MaxReconnectAttempts", 3)
	v.SetDefault("OrderCreatedConsumerSettings.RabbitMqSettings.ReconnectTimeoutSeconds", 5)
}
