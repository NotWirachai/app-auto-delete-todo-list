package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/niphawanphoopha/go-web-api/config"
)

var (
	// DB is the global database connection
	DB *gorm.DB
)

// Init initializes the database connection
func Init(cfg *config.Config) error {
	var err error
	
	// Connect to the database
	DB, err = gorm.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	
	// Set connection pool settings
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Hour)
	
	// Enable logging if debug mode is on
	DB.LogMode(cfg.Debug)
	
	log.Println("Database connection established")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// AutoMigrate runs auto migration for given models
func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...).Error
} 