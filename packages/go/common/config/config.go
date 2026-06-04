package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string
}

func Load() Config {
	viper.SetDefault("APP_NAME", "coffee-service")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()

	return Config{
		AppName: viper.GetString("APP_NAME"),
		AppEnv:  viper.GetString("APP_ENV"),
		AppPort: viper.GetString("APP_PORT"),
	}
}
