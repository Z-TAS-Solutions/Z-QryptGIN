package database

import (
	"time"

	"gorm.io/driver/postgres" 
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

	// Ensure all defined models are migrated to the DB schema.
	// If the tables (or new fields) are missing, GORM will create/alter them.
	if err := db.AutoMigrate(
		&User{},
		&Notification{},
		&CrypticRecord{},
		&MfaChallenge{},
		&ActivityLog{},
		&Session{},
		&WebAuthnCredential{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
