package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName   string
	Port      string
	JWTSecret string
	Debug     bool
}

var appConfig *Config

func init() {
	_ = godotenv.Load()

	appConfig = &Config{
		AppName:   getEnv("APP_NAME", "GO AUTH API"),
		Port:      getEnv("PORT", "8080"),
		Debug:     getAsBoolEnv("DEBUG", false),
		JWTSecret: mustGetEnv("JWT_SECRET"),
	}
}

func Get() *Config {
	return appConfig
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Missing required env var: %s", key)
	}
	return value
}

func getAsBoolEnv(key string, defaultVal bool) bool {
	valStr := getEnv(key, "")

	if val, err := strconv.ParseBool(valStr); err != nil {
		return val
	}

	return defaultVal
}
