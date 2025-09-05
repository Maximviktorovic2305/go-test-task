package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Структура конфигурации
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
	ServerPort int
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	// Загрузить файл .env, если он существует
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Предупреждение: файл .env не найден, используются переменные окружения")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "subscriptions"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
	}

	return config
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как целое число или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
