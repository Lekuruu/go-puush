package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/Lekuruu/go-puush/internal/config"
	"gorm.io/gorm"
)

func CreateSession(config config.DatabaseConfig) (*gorm.DB, error) {
	// Build DSN with configurable parameters
	journalMode := config.JournalMode
	if !config.EnableWAL {
		journalMode = "DELETE"
	}
	dsn := buildDatabaseDSN(config, journalMode)

	db, err := gorm.Open(openDatabase(dsn), &gorm.Config{
		Logger: NewGormLogger(),
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

func applyPerformanceSettings(db *gorm.DB, config config.DatabaseConfig) error {
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

// CloseSession applies final SQLite maintenance and closes the connection.
func CloseSession(db *gorm.DB, config config.DatabaseConfig) error {
	if db == nil {
		return nil
	}

	var maintenanceErr error
	if err := db.Exec("PRAGMA optimize;").Error; err != nil {
		maintenanceErr = fmt.Errorf("database: failed to optimize: %w", err)
	}
	if err := CheckpointWAL(db, config); err != nil {
		maintenanceErr = errors.Join(maintenanceErr, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.Join(maintenanceErr, fmt.Errorf("database: failed to resolve SQL connection: %w", err))
	}
	if err := sqlDB.Close(); err != nil {
		return errors.Join(maintenanceErr, fmt.Errorf("database: failed to close: %w", err))
	}
	return maintenanceErr
}

// CheckpointWAL flushes the current write-ahead log when WAL mode is enabled.
func CheckpointWAL(db *gorm.DB, config config.DatabaseConfig) error {
	if db == nil || !config.EnableWAL {
		return nil
	}
	if err := db.Exec("PRAGMA wal_checkpoint(FULL);").Error; err != nil {
		return fmt.Errorf("database: failed to checkpoint WAL: %w", err)
	}
	return nil
}
