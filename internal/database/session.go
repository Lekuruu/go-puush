package database

import (
	"fmt"
	"log"
	"os"
	"time"

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
	dsn := buildDatabaseDSN(config, journalMode)

	db, err := gorm.Open(openDatabase(dsn), &gorm.Config{
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

	success := false
	defer func() {
		if !success {
			sqlDB.Close()
		}
	}()

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	// Apply additional sqlite performance settings
	if err := applyPerformanceSettings(db, config); err != nil {
		return nil, err
	}

	var models = []any{
		&User{},
		&Upload{},
		&Pool{},
		&Session{},
		&InvitationKey{},
		&EmailVerification{},
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}
	if err := db.Exec("PRAGMA optimize=0x10002;").Error; err != nil {
		return nil, err
	}

	success = true
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
