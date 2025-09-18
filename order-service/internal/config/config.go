package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type DbSettings struct {
	MigrationConnectionString string `mapstructure:"MigrationConnectionString"`
	ConnectionString          string `mapstructure:"ConnectionString"`
}

type ServerSettings struct {
	Port int `mapstructure:"Port"`
}

type Config struct {
	Db     DbSettings     `mapstructure:"DbSettings"`
	Server ServerSettings `mapstructure:"ServerSettings"`
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "Development"
	}

	v := viper.New()
	v.SetConfigName(fmt.Sprintf("appsettings.%s", env))
	v.SetConfigType("json")
	v.AddConfigPath("/etc/order-service")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config read error: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}
	return &cfg, nil
}
