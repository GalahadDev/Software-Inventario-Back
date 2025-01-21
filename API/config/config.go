package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func DBURL() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env o hubo un error al cargarlo. Se usarán las variables de entorno del sistema.")
	}

	DBUser := strings.TrimSpace(os.Getenv("user"))
	DBPassword := strings.TrimSpace(os.Getenv("password"))
	DBHost := strings.TrimSpace(os.Getenv("host"))
	DBPort := strings.TrimSpace(os.Getenv("port"))
	DBName := strings.TrimSpace(os.Getenv("dbname"))

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&statement_cache_mode=describe", DBUser, DBPassword, DBHost, DBPort, DBName)
}
