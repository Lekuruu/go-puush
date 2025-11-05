//go:build !cgo
// +build !cgo

// Conditional build for database opener - cgo-less version (used when CGO_ENABLED=0)

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func openDatabase(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
