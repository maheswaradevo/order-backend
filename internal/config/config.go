package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	AppEnvironment        string
	ProductionEnvironment string
	Database              Database
	RabbitMqConfig        RabbitMQConfig
}

type Database struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Address  string
	Port     string
	SSL      bool
}

var config Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR .env Not found")
	}
	// setup environment
	config.AppEnvironment = os.Getenv("APP_ENV")
	config.Port = os.Getenv("PORT")
	config.ProductionEnvironment = os.Getenv("PRODUCTION_ENV")

	// setup database
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Address = os.Getenv("DB_ADDRESS")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")

	// setup rabbitmq
	config.RabbitMqConfig.Username = os.Getenv("RABBIT_USERNAME")
	config.RabbitMqConfig.Password = os.Getenv("RABBIT_PASSWORD")
	config.RabbitMqConfig.Address = os.Getenv("RABBIT_HOST")
	config.RabbitMqConfig.Port = os.Getenv("RABBIT_PORT")
	if os.Getenv("RABBIT_SSL") == "true" {
		config.RabbitMqConfig.SSL = true
	} else {
		config.RabbitMqConfig.SSL = false
	}
}

func GetConfig() *Config {
	return &config
}
