package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DbSettings struct {
	ConnectionString          string `mapstructure:"ConnectionString"`
	MigrationConnectionString string `mapstructure:"MigrationConnectionString"`
}

type HttpSettings struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}

type Config struct {
	DbSettings DbSettings   `mapstructure:"DbSettings"`
	Http       HttpSettings `mapstructure:"Http"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	env := v.GetString("APP_ENVIRONMENT")
	if env == "" {
		env = "Development"
	}

	v.SetConfigType("yaml")
	v.SetConfigName("appsettings." + env)
	v.AddConfigPath("/etc/order-service")
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
