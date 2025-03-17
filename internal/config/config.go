package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func LoadEnv() {
	const op = "internal.config.LoadEnv"

	if err := godotenv.Load("F:\\TestEffectiveMobile\\.env"); err != nil {
		slog.Error(op, "Файл .env не найден", slog.String("error", err.Error()))
	}
}

func GetDBConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func GetExternalAPIURL() string {
	return os.Getenv("EXTERNAL_URL")
}
