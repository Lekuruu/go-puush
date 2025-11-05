//go:build cgo
// +build cgo

// Conditional build for database opener - cgo version (used when CGO_ENABLED=1)

package database

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
