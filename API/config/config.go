package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DBURL() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error al cargar el archivo .env")
	}

	DBHost := os.Getenv("host")
	DBUser := os.Getenv("user")
	DBPassword := os.Getenv("password")
	DBPort := os.Getenv("port")
	DBName := os.Getenv("dbname")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&statement_cache_mode=describe", DBUser, DBPassword, DBHost, DBPort, DBName)
}
