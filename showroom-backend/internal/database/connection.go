package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
)

// DB holds the database connection
var DB *sql.DB

// Connect establishes database connection
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
		cfg.Database.Timezone,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// CreateDatabase creates the database if it doesn't exist
func CreateDatabase(cfg *config.Config) error {
	// Connect to postgres database to create our target database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.SSLMode,
		cfg.Database.Timezone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", cfg.Database.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s created successfully", cfg.Database.Name)
	} else {
		log.Printf("Database %s already exists", cfg.Database.Name)
	}

	return nil
}