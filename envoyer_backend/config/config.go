package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

var Config Conf

type Conf struct {
	AppName                string        `mapstructure:"APP_NAME"`
	AppCommand             string        `mapstructure:"APP_COMMAND"`
	Port                   string        `mapstructure:"SERVER_PORT"`
	DbUrl                  string        `mapstructure:"DB_URL"`
	RabbitMQUrl            string        `mapstructure:"RABBITMQ_URL"`
	MaxDbConnection        int           `mapstructure:"DB_MAX_CONNECTIONS"`
	GinMode                string        `mapstructure:"GIN_MODE"`
	LogLevel               string        `mapstructure:"LOG_LEVEL"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	SuperUserName          string        `mapstructure:"SUPER_USERNAME"`
	SuperPassword          string        `mapstructure:"SUPER_PASSWORD"`
	PublisherConfig
	ConsumerConfig
}

func LoadConfig() {
	var err error

	viper.SetConfigFile("base.env")
	viper.SetConfigType("props")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("WARNING: file .env not found")
	} else {
		viper.SetConfigFile(".env")
		viper.SetConfigType("props")
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Override config parameters from environment variables if specified
	err = viper.Unmarshal(&Config)
	for _, key := range viper.AllKeys() {
		err = viper.BindEnv(key)
		if err != nil {
			fmt.Println(err)
		}
	}

	LoadPublisherConfig()
	LoadConsumerConfig()
}

// LoadPublisherConfig load the default publisher config
func LoadPublisherConfig() {

	exchangeConfig := ExchangeConfig{
		Name:       "my_exchange",
		Type:       "x-delayed-message",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}

	Config.PublisherConfig = PublisherConfig{
		ServerUrl:      Config.RabbitMQUrl,
		ExchangeConfig: exchangeConfig,
		Mandatory:      false,
		Immidiate:      false,
		Exclusive:      false,
		ContentType:    "",
		BackoffPolicy:  []time.Duration{2, 4, 8, 16, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
	}
}

// LoadConsumerConfig load the default consumer config
func LoadConsumerConfig() {

	exchangeConfig := ExchangeConfig{
		Name:       "my_exchange",
		Type:       "x-delayed-message",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}

	queueConfig := QueueConfig{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}

	Config.ConsumerConfig = ConsumerConfig{
		ExchangeConfig: exchangeConfig,
		QueueConfig:    queueConfig,
		ServerUrl:      Config.RabbitMQUrl,
		NoWait:         false,
		AutoAck:        false,
		BackoffPolicy:  []time.Duration{2, 4, 8, 16, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32},
	}
}
