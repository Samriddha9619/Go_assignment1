package config
import (
    "os"
    "strconv"
)
type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    ServerPort string
}
func Load() *Config {
    port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     port,
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "hotel_db"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
    }
}
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
