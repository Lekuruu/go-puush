//go:build cgo
// +build cgo

// Conditional build for database opener - cgo version (used when CGO_ENABLED=1)

package database

import (
	"net/url"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

func buildDatabaseDSN(config DatabaseConfig, journalMode string) string {
	options := url.Values{}
	options.Set("_journal_mode", journalMode)
	options.Set("_busy_timeout", strconv.Itoa(config.BusyTimeout))
	options.Set("_synchronous", config.Synchronous)
	options.Set("_cache_size", strconv.Itoa(config.CacheSize))
	options.Set("cache", config.CacheMode)
	options.Set("_time_format", "sqlite")
	options.Set("_loc", "auto")
	return appendDatabaseOptions(config.Path, options)
}
