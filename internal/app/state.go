package app

import (
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/storage"
	"gorm.io/gorm"
)

type State struct {
	Config   *Config
	Database *gorm.DB
	Logger   *Logger
	Storage  storage.Storage
	Email    email.Email
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	logger := NewLogger("puush")

	if config.Storage.Type != "local" {
		logger.Logf("Unsupported storage type: %s", config.Storage.Type)
		return nil
	}

	fs := storage.NewFileStorage(config.Storage.Uri)
	err = fs.Setup()
	if err != nil {
		logger.Logf("Failed to setup file storage: %v", err)
		return nil
	}

	db, err := database.CreateSession(config.Database)
	if err != nil {
		logger.Logf("Failed to create database session: %v", err)
		return nil
	}

	mailer, err := email.NewEmailFromConfig(config.Email.Type, config.Email.From)
	if err != nil {
		logger.Logf("Failed to create email service: %v", err)
		return nil
	}

	if err := mailer.Setup(); err != nil {
		logger.Logf("Failed to setup email service: %v", err)
		return nil
	}

	return &State{
		Logger:   logger,
		Config:   config,
		Database: db,
		Email:    mailer,
		Storage:  fs,
	}
}
