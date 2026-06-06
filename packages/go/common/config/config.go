package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName                  string
	AppEnv                   string
	AppPort                  string
	DatabaseHost             string
	DatabasePort             string
	DatabaseUser             string
	DatabasePassword         string
	DatabaseName             string
	JWTSecret                string
	JWTAccessTokenTTLMinutes int
	JWTRefreshTokenTTLDays   int
}

func Load() Config {
	viper.SetDefault("APP_NAME", "coffee-service")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("JWT_SECRET", "dev-secret")
	viper.SetDefault("JWT_ACCESS_TOKEN_TTL_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_TOKEN_TTL_DAYS", 7)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()

	return Config{
		AppName:                  viper.GetString("APP_NAME"),
		AppEnv:                   viper.GetString("APP_ENV"),
		AppPort:                  viper.GetString("APP_PORT"),
		DatabaseHost:             viper.GetString("DATABASE_HOST"),
		DatabasePort:             viper.GetString("DATABASE_PORT"),
		DatabaseUser:             viper.GetString("DATABASE_USER"),
		DatabasePassword:         viper.GetString("DATABASE_PASSWORD"),
		DatabaseName:             viper.GetString("DATABASE_NAME"),
		JWTSecret:                viper.GetString("JWT_SECRET"),
		JWTAccessTokenTTLMinutes: viper.GetInt("JWT_ACCESS_TOKEN_TTL_MINUTES"),
		JWTRefreshTokenTTLDays:   viper.GetInt("JWT_REFRESH_TOKEN_TTL_DAYS"),
	}
}
