package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Redis struct {
		Host     string
		Port     string
		DB       int
		Password string
	}
	HashMaxIterations int
	HashZerosCount    int
	HashDuration      time.Duration
}

func Load() (*Config, error) {
	var config Config

	// initialization server
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", "9001")

	// initialization redis
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	redisDb, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("error during init REDIS_DB: %w", err)
	}
	config.Redis.DB = redisDb

	// initialization hash max iterations
	hashMaxIterations, err := strconv.Atoi(getEnv("HASH_MAX_ITERATIONS", "100000"))
	if err != nil {
		return nil, fmt.Errorf("error during init HASH_MAX_ITERATIONS: %w", err)
	}
	config.HashMaxIterations = hashMaxIterations

	// initialization hash zeros count
	HashZerosCount, err := strconv.Atoi(getEnv("HASH_ZEROS_COUNT", "3"))
	if err != nil {
		return nil, fmt.Errorf("error during init HASH_ZEROS_COUNT: %w", err)
	}
	config.HashZerosCount = HashZerosCount

	// initialization hash duration
	hashDuration, err := time.ParseDuration(getEnv("HASH_DURATION", "500s"))
	if err != nil {
		return nil, fmt.Errorf("error during init HASH_DURATION: %w", err)
	}
	config.HashDuration = hashDuration
	return &config, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
