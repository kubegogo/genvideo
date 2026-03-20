package config

import "os"

type Config struct {
	ServerPort     string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	AIProvider    string
	MinimaxAPIKey string
	MinimaxBaseURL string
	N8NBaseURL    string
	ComfyUIBaseURL string
	OllamaBaseURL  string
	OSSEndpoint  string
	OSSAccessKey string
	OSSSecretKey string
	OSSBucket    string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "3004"),
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", "password"),
		MySQLDatabase: getEnv("MYSQL_DATABASE", "genvideo"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		AIProvider:    getEnv("AI_PROVIDER", "minimax"),
		MinimaxAPIKey: os.Getenv("MINIMAX_API_KEY"),
		MinimaxBaseURL: getEnv("MINIMAX_BASE_URL", "https://api.minimax.chat"),
		N8NBaseURL:    getEnv("N8N_BASE_URL", "http://localhost:5678"),
		ComfyUIBaseURL: getEnv("COMFYUI_BASE_URL", "http://localhost:8188"),
		OllamaBaseURL: getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		OSSEndpoint:  getEnv("OSS_ENDPOINT", "oss-cn-hangzhou.aliyuncs.com"),
		OSSAccessKey: os.Getenv("OSS_ACCESS_KEY"),
		OSSSecretKey: os.Getenv("OSS_SECRET_KEY"),
		OSSBucket:    getEnv("OSS_BUCKET", "genvideo"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
