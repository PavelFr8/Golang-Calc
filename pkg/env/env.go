package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
    // Загружаем значения из .env в систему
    if err := godotenv.Load(); err != nil {
        fmt.Println("No .env file found")
    }
}

// Read logging lever from environment variable
func GetLoggingLevelEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		switch value {
		case "DEBUG", "INFO", "WARNING", "ERROR":
			return value
		}
	}
	return defaultVal
}

// Get code bolow from Habr tutorial: https://habr.com/ru/articles/446468/

// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
	return value
    }

    return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func GetEnvAsInt(name string, defaultVal int) int {
    valueStr := GetEnv(name, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
	return value
    }

    return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func GetEnvAsBool(name string, defaultVal bool) bool {
    valStr := GetEnv(name, "")
    if val, err := strconv.ParseBool(valStr); err == nil {
	return val
    }

    return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func GetEnvAsSlice(name string, defaultVal []string, sep string) []string {
    valStr := GetEnv(name, "")

    if valStr == "" {
	return defaultVal
    }

    val := strings.Split(valStr, sep)

    return val
}