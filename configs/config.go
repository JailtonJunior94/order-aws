package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTPConfig HTTPConfig `mapstructure:",squash"`
		SQSConfig  SQSConfig  `mapstructure:",squash"`
	}

	HTTPConfig struct {
		Port string `mapstructure:"HTTP_PORT"`
	}

	SQSConfig struct {
		QueueName string `mapstructure:"QUEUE_NAME"`
	}
)

func LoadConfig(path string) (*Config, error) {
	var config *Config

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
