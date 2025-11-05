//go:build cgo_sqlite
// +build cgo_sqlite

// Conditional build for database opener - cgo version (for architectures unsupported by native driver, for example MIPS)

package database

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
