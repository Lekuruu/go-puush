package app

import (
	"log/slog"

	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/logging"
	"github.com/Lekuruu/go-puush/internal/storage"
	"gorm.io/gorm"
)

type State struct {
	Config   *Config
	Database *gorm.DB
	Logger   *slog.Logger
	Storage  storage.Storage
	Email    email.Email
}

func (state *State) ExecuteWalCheckpoint() {
	if state.Config == nil {
		return
	}
	if !state.Config.Database.EnableWAL {
		return
	}

	if state.Database == nil {
		return
	}
	state.Database.Exec("PRAGMA wal_checkpoint(FULL);")
}

func (state *State) Close() {
	if state.Database == nil {
		return
	}
	if err := state.Database.Exec("PRAGMA optimize;").Error; err != nil && state.Logger != nil {
		state.Logger.Error("Failed to optimize database", "error", err)
	}
	state.ExecuteWalCheckpoint()
	db, err := state.Database.DB()
	if err != nil {
		return
	}
	db.Close()
}

func NewState() *State {
	logging.SetDefault("puush", slog.LevelInfo)
	logger := slog.Default()

	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	if config.Storage.Type != "local" {
		logger.Error("Unsupported storage type", "storage_type", config.Storage.Type)
		return nil
	}

	fs := storage.NewFileStorage(config.Storage.Uri)
	err = fs.Setup()
	if err != nil {
		logger.Error("Failed to set up file storage", "error", err)
		return nil
	}

	db, err := database.CreateSession(config.Database)
	if err != nil {
		logger.Error("Failed to create database session", "error", err)
		return nil
	}

	mailer, err := email.NewEmailFromConfig(config.Email.Type, config.Email.From)
	if err != nil {
		logger.Error("Failed to create email service", "error", err)
		return nil
	}

	if err := mailer.Setup(); err != nil {
		logger.Error("Failed to set up email service", "error", err)
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
