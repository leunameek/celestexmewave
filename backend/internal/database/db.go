package database

import (
	"fmt"
	"log"
	"time"

	"github.com/leunameek/celestexmewave/internal/config"
	"github.com/leunameek/celestexmewave/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize initializes the database connection with Supabase PostgreSQL
func Initialize(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	// Validate DSN is configured
	if dsn == "" {
		return fmt.Errorf("database configuration error: DATABASE_URL environment variable is required for Supabase connection")
	}

	// Configure GORM with connection pool settings optimized for Supabase Session Pooler
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Supabase database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool for Session Pooler
	// Session Pooler: Best for web applications with many concurrent connections
	// - Max open connections: 25 (Supabase Session Pooler default)
	// - Max idle connections: 5
	// - Connection max lifetime: 30 minutes
	// - Connection max idle time: 5 minutes
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	DB = db
	log.Println("✓ Connected to PostgreSQL database (Session Pooler mode)")

	// Run migrations
	if err := Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	return DB.AutoMigrate(
		&models.Store{},
		&models.Product{},
		&models.User{},
		&models.PasswordReset{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
