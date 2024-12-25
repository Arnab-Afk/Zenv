package database

import (
	"fmt"

	"github.com/arnab-afk/Zenv/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Secret struct {
	gorm.Model
	UserID  uint
	Name    string
	Value   []byte // Encrypted value
	Version int
}

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
	DB.AutoMigrate(&Secret{})
}
