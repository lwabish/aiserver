package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabaseURL string `mapstructure:"database_url"`
	LogLevel    string `mapstructure:"log_level"`
	ServerPort  string `mapstructure:"server_port"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return
}
