package internal

type State struct {
	config *Config
	// TODO: Add database, storage, ...
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return &State{config: config}
}
