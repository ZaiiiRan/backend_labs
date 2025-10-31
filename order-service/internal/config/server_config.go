package config

import (
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	DbSettings                            settings.DbSettings                `mapstructure:"DbSettings"`
	OrderCreatedRabbitMqPublisherSettings settings.RabbitMqPublisherSettings `mapstructure:"OrderCreatedPublisherSettings"`
	Http                                  settings.HttpServerSettings        `mapstructure:"HttpServerSettings"`
	Grpc                                  settings.GrpcServerSettings        `mapstructure:"GrpcServerSettings"`
}

func LoadServerConfig() (*ServerConfig, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	env := v.GetString("APP_ENVIRONMENT")
	if env == "" {
		env = "Development"
	}

	setDefaultServerConfigValues(v)

	v.SetConfigType("yaml")
	v.SetConfigName("appsettings." + env)
	v.AddConfigPath("/etc/order-service")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func setDefaultServerConfigValues(v *viper.Viper) {
	v.SetDefault("OrderCreatedRabbitMqPublisherSettings.RabbitMqSettings.HeartbeatSeconds", 30)
	v.SetDefault("OrderCreatedRabbitMqPublisherSettings.RabbitMqSettings.MaxReconnectAttempts", 3)
	v.SetDefault("OrderCreatedRabbitMqPublisherSettings.RabbitMqSettings.ReconnectTimeoutSeconds", 5)
	v.SetDefault("HttpServerSettings.Port", 5000)
	v.SetDefault("GrpcServerSettings.Port", 50051)
}
