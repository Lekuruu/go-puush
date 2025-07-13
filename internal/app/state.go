package app

import (
	"log/slog"

	"github.com/Lekuruu/go-puush/internal/storage"
	"github.com/sytallax/prettylog"
)

type State struct {
	Config  *Config
	Logger  *slog.Logger
	Storage storage.Storage
	// TODO: Add database, storage, ...
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
	handler := prettylog.NewHandler(nil)
	logger := slog.New(handler)

	return &State{
		Config:  config,
		Logger:  logger,
		Storage: fs,
	}
}
