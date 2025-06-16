package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	Env      string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret            string
	ExpirationMinutes int
}

type ServerConfig struct {
	Port string
}

func LoadConfig(envFile string) *Config {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading environment file %s: %v", envFile, err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	jwtExp, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRATION_MINUTES value: %v", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:            os.Getenv("JWT_SECRET"),
			ExpirationMinutes: jwtExp,
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Env: os.Getenv("ENV"),
	}
}
