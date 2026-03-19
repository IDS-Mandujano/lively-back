package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // El guión bajo es importante, carga el driver en memoria
	"github.com/joho/godotenv"
)

// DB es la variable global que usaremos en nuestros repositorios para consultar cosas
var DB *sql.DB

func Connect() {
	// 1. Cargamos las variables del archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env o no se pudo cargar.")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 2. Armamos la cadena de conexión (DSN)
	// parseTime=true es crucial para que los TIMESTAMP de MySQL se conviertan a time.Time en Go
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	// 3. Abrimos la conexión
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}

	// 4. Verificamos que la base de datos realmente responda
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error haciendo ping a MySQL. Revisa tus credenciales en el .env: %v", err)
	}

	log.Println("¡Conexión exitosa a MySQL (lively_db)!")
}
