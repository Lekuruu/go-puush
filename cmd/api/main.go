package main

import "github.com/Lekuruu/go-puush/api"

func main() {
	state := app.NewState()
	server := api.NewServer(state)
	server.Serve()
}
