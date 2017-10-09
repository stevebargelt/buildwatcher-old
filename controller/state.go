package controller

import (
	"log"

	"github.com/kidoman/embd"
)

type State struct {
	config Config
	store  *Store
}

func NewState(c Config, store *Store) *State {
	return &State{
		config: c,
		store:  store,
	}
}

func (s *State) Bootup() error {
	if s.config.EnableGPIO {
		log.Println("Enabled GPIO subsystem")
		embd.InitGPIO()
	}
	return nil
}

func (s *State) TearDown() {
	if s.config.EnableGPIO {
		embd.CloseGPIO()
		log.Println("Stopping GPIO subsystem")
	}
}
