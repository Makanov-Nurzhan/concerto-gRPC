package db

import (
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error with connect to database: %v", err)
	}
	log.Println("Connected to database")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error gettin sql.DB from gorm: %v", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	return db
}
