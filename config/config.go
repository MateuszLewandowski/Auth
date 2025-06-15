package config

import (
	"os"
	"strconv"
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

func LoadConfig() *Config {
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	jwtExp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))

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
