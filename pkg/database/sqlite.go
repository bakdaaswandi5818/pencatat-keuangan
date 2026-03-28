package database

import (
	"fmt"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens a SQLite database at dsn and runs auto-migrations.
// WAL mode and a busy timeout are set for better concurrency on a shared VPS.
func New(dsn string) (*gorm.DB, error) {
	// Append SQLite pragmas to the DSN.
	fullDSN := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", dsn)

	db, err := gorm.Open(sqlite.Open(fullDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(&domain.Transaction{}); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	return db, nil
}
