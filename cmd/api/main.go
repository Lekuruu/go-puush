package main

import (
	"github.com/Lekuruu/go-puush/api"
	"github.com/Lekuruu/go-puush/internal/app"
)

func main() {
	state := app.NewState()
	if state == nil {
		return
	}

	server := api.NewServer(state)
	server.Serve()
}
