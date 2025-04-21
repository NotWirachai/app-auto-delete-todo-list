package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the API
type Config struct {
	Port         int
	Host         string
	ReadTimeout  int
	WriteTimeout int
	Debug        bool
	
	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	// JWT configuration
	JWTSecret string
	JWTExpiry int // in minutes
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Port:         getEnvAsInt("PORT", 8080),
		Host:         getEnv("HOST", "0.0.0.0"),
		ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
		WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
		Debug:        getEnvAsBool("DEBUG", false),
		
		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "go_web_api"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		
		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiry: getEnvAsInt("JWT_EXPIRY", 60), // 60 minutes default
	}
}

// GetDBConnString returns the PostgreSQL connection string
func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
} 