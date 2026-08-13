package state

import (
	"fmt"
	"log/slog"

	"github.com/Lekuruu/go-puush/internal/config"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/logging"
	"github.com/Lekuruu/go-puush/internal/storage"
	"gorm.io/gorm"
)

// State holds the shared application dependencies.
type State struct {
	Config   *config.Config
	Database *gorm.DB
	Logger   *slog.Logger
	Storage  storage.Storage
	Email    email.Email
}

func NewState(environmentFiles ...string) (*State, error) {
	cfg, err := config.LoadConfig(environmentFiles...)
	if err != nil {
		return nil, fmt.Errorf("state: failed to load config: %w", err)
	}

	logging.SetDefault("puush", slog.LevelInfo)
	logger := slog.Default()

	if cfg.Storage.Type != "local" {
		return nil, fmt.Errorf("state: unsupported storage type %q", cfg.Storage.Type)
	}

	storageProvider := storage.NewFileStorage(cfg.Storage.Uri)
	if err := storageProvider.Setup(); err != nil {
		return nil, fmt.Errorf("state: failed to set up storage: %w", err)
	}

	db, err := database.CreateSession(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("state: failed to create database session: %w", err)
	}

	mailer, err := email.NewEmailFromConfig(cfg.Email)
	if err != nil {
		database.CloseSession(db, cfg.Database)
		return nil, fmt.Errorf("state: failed to create email service: %w", err)
	}
	if err := mailer.Setup(); err != nil {
		database.CloseSession(db, cfg.Database)
		return nil, fmt.Errorf("state: failed to set up email service: %w", err)
	}

	return &State{
		Config:   cfg,
		Database: db,
		Logger:   logger,
		Storage:  storageProvider,
		Email:    mailer,
	}, nil
}

// Close gracefully closes application resources.
func (state *State) Close() error {
	if state == nil || state.Config == nil {
		return nil
	}
	return database.CloseSession(state.Database, state.Config.Database)
}

// CheckpointWAL flushes pending SQLite WAL contents to the database file.
func (state *State) CheckpointWAL() error {
	if state == nil || state.Config == nil {
		return nil
	}
	return database.CheckpointWAL(state.Database, state.Config.Database)
}
