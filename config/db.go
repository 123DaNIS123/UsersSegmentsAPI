package config

import (
	"fmt"
	"os"

	"github.com/123DaNIS123/UsersSegments/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		print(db)
		print("!!!!")
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	DB = db
}