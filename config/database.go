package config

import (
	"agnos-assignment/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=postgres user=postgres password=postgres dbname=his port=5432 sslmode=disable"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Hospital{},
		&models.Staff{},
		&models.Patient{},
	)
}

func ConnectDatabaseAndMigrate() *gorm.DB {
	db, err := ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("database connected")

	err = db.AutoMigrate(
		&models.Hospital{},
		&models.Patient{},
		&models.Staff{},
	)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
