package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// DB Config
	DB_Driver string `mapstructure:"DB_DRIVER"`
	DB_Source string `mapstructure:"DB_SOURCE"`
	// Web Configs
	WEB_Port      string `mapstructure:"WEB_PORT"`
	CORS_Frontend string `mapstructure:"CORS_FRONTEND"`
	CORS_Enable   bool   `mapstructure:"CORS_ENABLE"`
	// JWT Configs
	JWT_SECRET   string        `mapstructure:"JWT_SECRET"`
	JWT_DURATION time.Duration `mapstructure:"JWT_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return
}
