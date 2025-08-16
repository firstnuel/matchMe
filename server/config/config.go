package config

import (
	"sync"

	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		cfg = &Config{
			AppEnv:        mustGetEnv("APP_ENV"),
			Port:          getEnvStr("PORT", "8080"),
			Host:          mustGetEnv("HOST"),
			DbURL:         mustGetEnv("DATABASE_URL"),
			DbName:        mustGetEnv("DATABASE_NAME"),
			JWTSecret:     mustGetEnv("JWT_SECRET"),
			ServerAddr:    mustGetEnv("SERVER_ADDR"),
			ClientAddr:    mustGetEnv("CLIENT_ADDR"),
			CloudinaryUrl: getEnvStr("CLOUDINARY_URL", ""),
		}
	})

	return cfg
}
