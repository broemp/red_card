package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// DB Config
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBAdress   string `mapstructure:"DB_ADDRESS"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
	// Web Configs
	WebPort string `mapstructure:"WEB_PORT"`
	// JWT Configs
	JWT_SECRET   string        `mapstructure:"JWT_SECRET"`
	JWT_DURATION time.Duration `mapstructure:"JWT_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
