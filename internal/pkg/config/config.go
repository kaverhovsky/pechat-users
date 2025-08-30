package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"postgres"`

	Logger struct {
		Level string `mapstructure:"level"`
		Mode  string `mapstructure:"mode"`
	} `mapstructure:"logger"`

	HTTP struct {
		Listen string `mapstructure:"listen"`
	} `mapstructure:"http"`
}

// LoadConfig TODO добавить opts для тестов
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	// сделать viper.SetEnvPrefix через opts

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
