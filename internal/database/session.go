package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateSession(config DatabaseConfig) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Build DSN with configurable parameters
	journalMode := config.JournalMode
	if !config.EnableWAL {
		journalMode = "DELETE"
	}

	dsn := fmt.Sprintf("%s?_journal_mode=%s&_busy_timeout=%d&_synchronous=%s&cache=%s",
		config.Path,
		journalMode,
		config.BusyTimeout,
		config.Synchronous,
		config.CacheMode,
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool with values from config
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	// Apply additional sqlite performance settings
	if err := applyPerformanceSettings(db, config); err != nil {
		return nil, err
	}

	var models = []interface{}{
		&User{},
		&Upload{},
		&Pool{},
		&Session{},
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}

	// Run migration to move shortlink identifiers to uploads
	if err := MigrateShortLinksToUploads(db); err != nil {
		return nil, err
	}

	return db, nil
}

func applyPerformanceSettings(db *gorm.DB, config DatabaseConfig) error {
	if err := db.Exec(fmt.Sprintf("PRAGMA cache_size = %d;", config.CacheSize)).Error; err != nil {
		return err
	}

	if config.EnableWAL {
		if err := db.Exec(fmt.Sprintf("PRAGMA wal_autocheckpoint = %d;", config.WALAutoCheckpoint)).Error; err != nil {
			return err
		}
	}

	if err := db.Exec("PRAGMA mmap_size = 67108864;").Error; err != nil {
		return err
	}

	if err := db.Exec("PRAGMA temp_store = MEMORY;").Error; err != nil {
		return err
	}

	return nil
}

func MigrateShortLinksToUploads(db *gorm.DB) error {
	// Check if there are any shortlinks to migrate
	var count int64
	if err := db.Table("short_links").Count(&count).Error; err != nil {
		return nil
	}

	if count == 0 {
		log.Println("No shortlinks to migrate, dropping short_links table...")
		return db.Exec("DROP TABLE IF EXISTS short_links").Error
	}

	log.Printf("Migrating %d shortlinks to uploads table...\n", count)

	// Update uploads with their shortlink identifiers
	result := db.Exec(`
		UPDATE uploads 
		SET identifier = (
			SELECT identifier 
			FROM short_links 
			WHERE short_links.upload_id = uploads.id
		)
		WHERE id IN (
			SELECT upload_id 
			FROM short_links
		)
	`)

	if result.Error != nil {
		return fmt.Errorf("failed to migrate shortlinks: %w", result.Error)
	}

	log.Printf("Successfully migrated %d shortlinks\n", result.RowsAffected)

	// Drop the short_links table after successful migration
	if err := db.Exec("DROP TABLE IF EXISTS short_links").Error; err != nil {
		return fmt.Errorf("failed to drop short_links table: %w", err)
	}

	log.Println("Successfully dropped short_links table")
	return nil
}
