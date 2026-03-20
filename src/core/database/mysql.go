package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env o no se pudo cargar.")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error haciendo ping a MySQL. Revisa tus credenciales en el .env: %v", err)
	}

	log.Println("¡Conexión exitosa a MySQL (lively_db)!")
}
