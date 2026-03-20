package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	ServerPort string

	// MySQL
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string

	// Redis
	RedisHost string
	RedisPort string
	RedisPassword string

	// AI Provider (minimax or self-hosted)
	AIProvider     string
	MinimaxAPIKey   string
	MinimaxBaseURL  string

	// Self-hosted (n8n + ComfyUI + Ollama)
	N8NBaseURL      string
	ComfyUIBaseURL  string
	OllamaBaseURL   string

	// OSS
	OSSEndpoint     string
	OSSAccessKey    string
	OSSSecretKey    string
	OSSBucket       string

	// Video download sources
	VideoSources    []string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "3004"),
		MySQLHost:      getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:      getEnv("MYSQL_PORT", "3306"),
		MySQLUser:      getEnv("MYSQL_USER", "root"),
		MySQLPassword:  getEnv("MYSQL_PASSWORD", "password"),
		MySQLDatabase:  getEnv("MYSQL_DATABASE", "genvideo"),
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		AIProvider:     getEnv("AI_PROVIDER", "minimax"),
		MinimaxAPIKey:  os.Getenv("MINIMAX_API_KEY"),
		MinimaxBaseURL: getEnv("MINIMAX_BASE_URL", "https://api.minimax.chat"),
		N8NBaseURL:     getEnv("N8N_BASE_URL", "http://localhost:5678"),
		ComfyUIBaseURL: getEnv("COMFYUI_BASE_URL", "http://localhost:8188"),
		OllamaBaseURL:  getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		OSSEndpoint:    getEnv("OSS_ENDPOINT", "oss-cn-hangzhou.aliyuncs.com"),
		OSSAccessKey:   os.Getenv("OSS_ACCESS_KEY"),
		OSSSecretKey:  os.Getenv("OSS_SECRET_KEY"),
		OSSBucket:      getEnv("OSS_BUCKET", "genvideo"),
		VideoSources:   []string{"douyin", "kuaishou", "bilibili", "xiaohongshu"},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
