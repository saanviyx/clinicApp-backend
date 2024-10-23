package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the application
type Config struct {
	DBConnStr      string
	DBMaxIdleConns int // Maximum number of idle connections
	DBMaxOpenConns int // Maximum number of open connections
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	// Get database connection string from environment variables
	dbConnStr := getRequiredEnv("DB_CONN_STR")
	return &Config{
		DBConnStr:      dbConnStr,
		DBMaxIdleConns: 10,  // Adjust the default value as needed
		DBMaxOpenConns: 100, // Adjust the default value as needed
	}
}

// getRequiredEnv retrieves an environment variable and panics if it is not set
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
