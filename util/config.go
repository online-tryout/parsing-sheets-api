package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress      string `mapstructure:"BACKEND_SERVER_ADDRESS"`
	RabbitSource       string `mapstructure:"RABBIT_SOURCE"`
	ServerUrl          string `mapstructure:"SERVER_URL"`
	BackendSwaggerHost string `mapstructure:"BACKEND_SWAGGER_HOST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
