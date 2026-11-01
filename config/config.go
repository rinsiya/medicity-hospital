package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	Port    string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	LogLevel  string
}

func LoadConfig() *Config {

	err := godotenv.Load("../.env")

	if err != nil {

		log.Println(".env file not found")
	}

	return &Config{
		AppName: os.Getenv("APP_NAME"),
		Port:    os.Getenv("PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret: os.Getenv("JWT_SECRET"),
		LogLevel:  os.Getenv("LOG_LEVEL"),
	}
}
