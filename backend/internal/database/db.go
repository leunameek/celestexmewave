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

// Initialize monta la conexion con Postgres en Supabase, todo fresh
func Initialize(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	// Validamos que si venga el DSN
	if dsn == "" {
		return fmt.Errorf("database configuration error: DATABASE_URL environment variable is required for Supabase connection")
	}

	// Configuramos GORM con el pool ideal pa Supabase Session Pooler
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Supabase database: %w", err)
	}

	// Agarramos la SQL DB para tunear el pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Pool para Session Pooler:
	// - Max open: 25
	// - Max idle: 5
	// - Lifetime: 30 mins
	// - Idle time: 5 mins
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	DB = db
	log.Println("✓ Connected to PostgreSQL database (Session Pooler mode)")

	// Corremos migraciones
	if err := Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Migrate corre todas las migras
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

// GetDB devuelve la instancia
func GetDB() *gorm.DB {
	return DB
}

// Close cierra la conexion
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
