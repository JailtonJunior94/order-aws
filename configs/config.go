package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTPConfig     HTTPConfig     `mapstructure:",squash"`
		AWSConfig      AWSConfig      `mapstructure:",squash"`
		SQSConfig      SQSConfig      `mapstructure:",squash"`
		S3Config       S3Config       `mapstructure:",squash"`
		DynamoDBConfig DynamoDBConfig `mapstructure:",squash"`
	}

	HTTPConfig struct {
		Port string `mapstructure:"HTTP_PORT"`
	}

	AWSConfig struct {
		Region     string `mapstructure:"AWS_REGION"`
		Endpoint   string `mapstructure:"AWS_ENDPOINT"`
		S3Endpoint string `mapstructure:"AWS_S3_ENDPOINT"`
	}

	SQSConfig struct {
		QueueName           string `mapstructure:"QUEUE_NAME"`
		MaxNumberOfMessages int32  `mapstructure:"MAX_NUMBER_OF_MESSAGES"`
		WaitTimeSeconds     int32  `mapstructure:"WAIT_TIME_SECONDS"`
		VisibilityTimeout   int32  `mapstructure:"VISIBILITY_TIMEOUT"`
	}

	S3Config struct {
		BucketName string `mapstructure:"BUCKET_NAME"`
	}

	DynamoDBConfig struct {
		TableName string `mapstructure:"TABLE_NAME"`
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
