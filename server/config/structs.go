package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	AppEnv        string
	Port          string
	Host          string
	DbURL         string
	DbName        string
	JWTSecret     string
	ClientAddr    []string
	CloudinaryUrl string
}

// Helper function to get required environment variable as string
func mustGetEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		log.Fatalf("Environment variable %s is required but not set or empty", key)
	}
	return val
}

// Helper function to get environment variable as string with default
func getEnvStr(key, defaultVal string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		log.Printf("Warning: No value found for '%s', using default value", key)
		return defaultVal
	}
	return val
}

// Helper function to get environment variable as []string (JSON array) with default
func getEnvStrArray(key string, defaultVal []string) []string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		log.Printf("Warning: No value found for '%s', using default value", key)
		return defaultVal
	}

	var arr []string
	if err := json.Unmarshal([]byte(val), &arr); err != nil {
		log.Printf("Error parsing '%s' as JSON array: %v, using default value", key, err)
		return defaultVal
	}

	return arr
}
