package config

import (
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/spf13/viper"
)

type ConsumerConfig struct {
	RabbitMqSettings      settings.RabbitMqSettings   `mapstructure:"RabbitMqSettings"`
	OmsClientHttpSettings settings.HttpClientSettings `mapstructure:"OmsHttpClient"`
}

func LoadConsumerConfig() (*ConsumerConfig, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	env := v.GetString("APP_ENVIRONMENT")
	if env == "" {
		env = "Development"
	}

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
