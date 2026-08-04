//go:build !cgo
// +build !cgo

// Conditional build for database opener - cgo-less version (used when CGO_ENABLED=0)

package database

import (
	"fmt"
	"net/url"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

func buildDatabaseDSN(config DatabaseConfig, journalMode string) string {
	options := url.Values{}
	options.Set("cache", config.CacheMode)
	options.Set("_time_format", "sqlite")
	options.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", config.BusyTimeout))
	options.Add("_pragma", fmt.Sprintf("journal_mode(%s)", journalMode))
	options.Add("_pragma", fmt.Sprintf("synchronous(%s)", config.Synchronous))
	options.Add("_pragma", fmt.Sprintf("cache_size(%d)", config.CacheSize))
	if config.EnableWAL {
		options.Add("_pragma", fmt.Sprintf("wal_autocheckpoint(%d)", config.WALAutoCheckpoint))
	}
	options.Add("_pragma", "mmap_size(67108864)")
	options.Add("_pragma", "temp_store(MEMORY)")
	return appendDatabaseOptions(config.Path, options)
}
