package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	config *Config
)

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Jika .env tidak ditemukan, gunakan variabel lingkungan
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using system environment variables: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Server.AllowedOrigins = strings.Split(viper.GetString("ALLOWED_ORIGINS"), ",")

	return config, nil
}
