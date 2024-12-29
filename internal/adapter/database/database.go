package database

import (
	"fmt"
	"log"
	"os"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

// connecting to postgres db using gorm
func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Singapore",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&domain.RunningEvents{}, &domain.User{})

	DB = Dbinstance{
		Db: db,
	}
}
