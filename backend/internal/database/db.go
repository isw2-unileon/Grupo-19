package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is variable we will use to access to the DB in our project
var DB *gorm.DB

// Initialize the conection to the DB
func Connect(databaseURL string) {
	var err error

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	log.Println("✅ PostgreSQL conexion established correctly!")
}
