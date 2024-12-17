package config

import (
	"fmt"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct{
	DatabaseURL		string
	DatabaseCA		string
	ServerAddress	string
	JWTSecret		string
	DatabaseDriver	string
	MaxConnections	int
}

func LoadConfig() (*Config, error){
	err := godotenv.Load()
	if err != nil{
		fmt.Println(err)
	}

	config := &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/mobile_banking"),
		ServerAddress:   getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		DatabaseDriver:  getEnv("DATABASE_DRIVER", "postgres"),
		MaxConnections:  getEnvAsInt("MAX_DB_CONNECTIONS", 10),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string{
	value := os.Getenv(key)
	if value == ""{
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}


