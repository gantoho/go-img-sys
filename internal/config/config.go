package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	File   FileConfig
	Auth   AuthConfig
	Rate   RateConfig
	Log    LogConfig
}

type ServerConfig struct {
	Port        string
	Env         string
	Timeout     int
	ExternalURL string
}

type AuthConfig struct {
	JWTSecret string
	JWTExpire time.Duration
}

type FileConfig struct {
	UploadDir         string
	MaxSize           int64
	AllowTypes        []string
	DuplicateStrategy string
}

type RateConfig struct {
	RequestsPerSec  int
	ConcurrentLimit int
}

type LogConfig struct {
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var appConfig *Config

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func Init() *Config {
	appConfig = &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", ":3128"),
			Env:         getEnv("SERVER_ENV", "development"),
			Timeout:     getEnvInt("SERVER_TIMEOUT", 30),
			ExternalURL: getEnv("EXTERNAL_URL", ""),
		},
		File: FileConfig{
			UploadDir:         getEnv("UPLOAD_DIR", "./files"),
			MaxSize:           int64(getEnvInt("MAX_FILE_SIZE", 100)),
			AllowTypes:        []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			DuplicateStrategy: getEnv("DUPLICATE_STRATEGY", "rename"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpire: 24 * time.Hour,
		},
		Rate: RateConfig{
			RequestsPerSec:  getEnvInt("RATE_LIMIT_REQUESTS", 100),
			ConcurrentLimit: getEnvInt("RATE_LIMIT_CONCURRENT", 10),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 7),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
		},
	}
	return appConfig
}

func GetConfig() *Config {
	if appConfig == nil {
		return Init()
	}
	return appConfig
}

func (c *Config) Reload() {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", ":3128"),
			Env:         getEnv("SERVER_ENV", "development"),
			Timeout:     getEnvInt("SERVER_TIMEOUT", 30),
			ExternalURL: getEnv("EXTERNAL_URL", ""),
		},
		File: FileConfig{
			UploadDir:         getEnv("UPLOAD_DIR", "./files"),
			MaxSize:           int64(getEnvInt("MAX_FILE_SIZE", 100)),
			AllowTypes:        []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			DuplicateStrategy: getEnv("DUPLICATE_STRATEGY", "rename"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpire: 24 * time.Hour,
		},
		Rate: RateConfig{
			RequestsPerSec:  getEnvInt("RATE_LIMIT_REQUESTS", 100),
			ConcurrentLimit: getEnvInt("RATE_LIMIT_CONCURRENT", 10),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 7),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
		},
	}
	*c = *cfg
}
