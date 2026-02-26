package repository

import (
	"log"

	"lively-backend/internal/models" 

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Room{},
	)
	if err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}

	log.Println("Base de datos conectada y tablas migradas con éxito")
	return db
}