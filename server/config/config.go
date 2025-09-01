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
			Host:          getEnvStr("HOST", "http://localhost"),
			DbURL:         mustGetEnv("DATABASE_URL"),
			DbName:        getEnvStr("DATABASE_NAME", "postgres"),
			JWTSecret:     mustGetEnv("JWT_SECRET"),
			ClientAddr:    getEnvStrArray("CLIENT_ADDR", []string{"*"}),
			CloudinaryUrl: getEnvStr("CLOUDINARY_URL", ""),
		}
	})

	return cfg
}
