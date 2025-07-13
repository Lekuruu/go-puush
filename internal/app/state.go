package app

type State struct {
	Config *Config
	// TODO: Add database, storage, ...
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return &State{Config: config}
}
