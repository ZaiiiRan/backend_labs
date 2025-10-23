package config

import (
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-generator/internal/config/settings"
	"github.com/spf13/viper"
)

type Config struct {
	OmsClientGrpcSettings settings.GrpcClientSettings `mapstructure:"OmsGrpcClient"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigType("yaml")
	v.SetConfigName("appsettings")
	v.AddConfigPath("/etc/order-service-consumer")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}