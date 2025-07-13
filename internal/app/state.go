package app

import (
	"log/slog"

	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/storage"
	"github.com/sytallax/prettylog"
	"gorm.io/gorm"
)

type State struct {
	Config   *Config
	Database *gorm.DB
	Logger   *slog.Logger
	Storage  storage.Storage
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	if config.Storage.Type != "local" {
		slog.Error("Unsupported storage type", "type", config.Storage.Type)
		return nil
	}

	fs := storage.NewFileStorage(config.Storage.Uri)
	err = fs.Setup()
	if err != nil {
		slog.Error("Failed to setup file storage", "error", err)
		return nil
	}

	db, err := database.CreateSession(config.Database.Path)
	if err != nil {
		slog.Error("Failed to create database session", "error", err)
		return nil
	}

	handler := prettylog.NewHandler(nil)
	logger := slog.New(handler)

	return &State{
		Config:   config,
		Logger:   logger,
		Database: db,
		Storage:  fs,
	}
}
