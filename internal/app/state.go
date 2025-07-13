package app

import (
	"log/slog"

	"github.com/sytallax/prettylog"
)

type State struct {
	Config *Config
	Logger *slog.Logger
	// TODO: Add database, storage, ...
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	handler := prettylog.NewHandler(nil)
	logger := slog.New(handler)

	return &State{
		Config: config,
		Logger: logger,
	}
}
