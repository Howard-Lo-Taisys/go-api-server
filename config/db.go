package config

import (
	"fmt"
	"go-api-server/models"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var once sync.Once

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Version{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Service{})
	if err != nil {
		panic(err)
	}

	initDB := func() {
		DB = db
	}

	once.Do(initDB)
}
