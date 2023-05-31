package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Environment      string
	LogLevel         string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PosgresDatabase  string
	SigninKey        string
	PGXPoolMax       int
	RedisHost        string
	RedisPort        string
	HttpPort         string
	AuthConfigPath   string
	CsvFilePath      string
}

func Load() Config {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("error connect .env loading ", err)
	}
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "developer"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "2002"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PosgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "climatedb"))
	c.SigninKey = cast.ToString(getOrReturnDefault("SINGINKEY", "murtazoxonSinginkey"))
	c.PGXPoolMax = cast.ToInt(getOrReturnDefault("PGX_POOL_MAX", 2))
	c.HttpPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":7070"))
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_FILE_PATH", "./config/auth.conf"))
	c.CsvFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/roles.csv"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
