package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables de sistema")
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
		user, pass, host, port, name)
}