package database

import (
	"time"

	"gorm.io/driver/postgres" // Or mysql/sqlite depending on your driver
	"gorm.io/gorm"
)

func NewDatabaseConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure the connection pool for maximum concurrency
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)           // Keep idle connections ready
	sqlDB.SetMaxOpenConns(100)          // Limit max open connections to prevent DB overload
	sqlDB.SetConnMaxLifetime(time.Hour) // Recycle connections every hour

	// Auto-migrate your models here (optional but helpful in dev)
	// db.AutoMigrate(&User{}, &Notification{}, ...)

	return db, nil
}
