//go:build !cgo_sqlite
// +build !cgo_sqlite

// Conditional build for database opener - cgo-less version (default)

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
