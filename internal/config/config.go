package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	AppPort    string `mapstructure:"APP_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPass     string `mapstructure:"DB_PASS"`
	DBName     string `mapstructure:"DB_NAME"`
	RedisHost  string `mapstructure:"REDIS_HOST"`
	RedisPort  string `mapstructure:"REDIS_PORT"`
	RedisPass  string `mapstructure:"REDIS_PASS"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return &config
}
