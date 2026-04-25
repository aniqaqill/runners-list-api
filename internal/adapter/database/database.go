package database

import (
	"fmt"
	"log/slog"

	"github.com/aniqaqill/runners-list/internal/config"
	"github.com/aniqaqill/runners-list/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a Postgres connection using the provided Config, runs
// auto-migrations, and returns the *gorm.DB handle.
//
// Key changes from the original:
//   - sslmode=require   → enforces TLS to Supabase (was "disable")
//   - No global var DB  → caller owns the reference and passes it down
//   - Returns error     → caller decides whether to os.Exit
func Connect(cfg *config.Config) (*gorm.DB, error) {
	sslMode := cfg.DBSSLMode
	if sslMode == "" {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Singapore",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	slog.Info("database connected", "host", cfg.DBHost)

	if err := db.AutoMigrate(&domain.Events{}, &domain.Users{}); err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}

	slog.Info("database migrations applied")

	return db, nil
}
