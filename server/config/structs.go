package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppEnv        string
	Port          string
	DbURL         string
	Host          string
	DbName        string
	JWTSecret     string
	ServerAddr    string
	ClientAddr    string
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

// Helper function to get environment variable as int with default
func getEnvInt(key string, defaultValue int) int {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, val, defaultValue)
		return defaultValue
	}

	return intVal
}

// Helper function to get environment variable as bool with default
func getEnvBool(key string, defaultValue bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s: %s, using default: %t", key, val, defaultValue)
		return defaultValue
	}

	return boolVal
}
