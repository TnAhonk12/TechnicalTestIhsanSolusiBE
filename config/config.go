// ===================== config/config.go =====================
package config

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	AppPort    string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		AppPort:    getArg("-port", "8080"),
	}
}

func getArg(flag string, defaultVal string) string {
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return defaultVal
}
