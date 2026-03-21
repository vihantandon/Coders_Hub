package boot

import (
	"github.com/vihantandon/Coders_Hub/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=avnivihan0 dbname=coders_hub port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	DB = db

	DB.AutoMigrate(&models.User{})
}
