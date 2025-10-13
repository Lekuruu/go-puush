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
