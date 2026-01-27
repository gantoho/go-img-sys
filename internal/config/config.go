package config

import "time"

type Config struct {
	Server ServerConfig
	File   FileConfig
	Auth   AuthConfig
}

type ServerConfig struct {
	Port    string
	Env     string
	Timeout int
}

type AuthConfig struct {
	JWTSecret string
	JWTExpire time.Duration
}

type FileConfig struct {
	UploadDir  string
	MaxSize    int64 // MB
	AllowTypes []string
	// DuplicateStrategy controls how to handle filename conflicts: "overwrite", "rename", "reject"
	DuplicateStrategy string
}

var AppConfig *Config

func Init() *Config {
	AppConfig = &Config{
		Server: ServerConfig{
			Port:    ":3128",
			Env:     "development",
			Timeout: 30,
		},
		File: FileConfig{
			UploadDir:         "./files",
			MaxSize:           100, // 100MB
			AllowTypes:        []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			DuplicateStrategy: "rename",
		},
		Auth: AuthConfig{
			JWTSecret: "your-secret-key-change-this-in-production", // Change this in production!
			JWTExpire: 24 * time.Hour,
		},
	}
	return AppConfig
}

func GetConfig() *Config {
	if AppConfig == nil {
		return Init()
	}
	return AppConfig
}
