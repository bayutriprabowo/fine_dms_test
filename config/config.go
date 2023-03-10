package config

import "enigmacamp.com/fine_dms/utils"

type DbConfig struct {
	Name           string
	Host, Port     string
	User, Password string
	SslMode        string
}

type ApiConfig struct {
	Host, Port string
}

type AppConfig struct {
	ApiConfig ApiConfig
	DbConfig  DbConfig
}

func NewAppConfig() AppConfig {
	return AppConfig{
		DbConfig: DbConfig{
			Name:     utils.GetEnv("DB_NAME"),
			Host:     utils.GetEnv("DB_HOST"),
			Port:     utils.GetEnv("DB_PORT"),
			User:     utils.GetEnv("DB_UNAME"),
			Password: utils.GetEnv("DB_PASSW"),
			SslMode:  utils.GetEnv("DB_SSL_MODE"),
		},
		ApiConfig: ApiConfig{
			Host: utils.GetEnv("HTTP_SERVER_HOST"),
			Port: utils.GetEnv("HTTP_SERVER_PORT"),
		},
	}
}
