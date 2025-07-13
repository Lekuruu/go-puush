package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateSession(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var models = []interface{}{
		&User{},
		&Upload{},
		&Pool{},
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}
	return db, nil
}
