package api

import (
	"log"
	"net/http"

	"github.com/stevebargelt/buildwatcher/controller"
)

type ServerConfig struct {
	Address string `yaml:"address"`
	Display bool   `yaml:"display"`
}

var DefaultConfig = ServerConfig{
	Address: ":9002",
}

type Server struct {
	config ServerConfig
}

func SetupServer(config ServerConfig, c *controller.Controller) error {
	server := &Server{
		config: config,
	}

	http.Handle("/api/", NewAPIHandler(c))
	log.Printf("Starting http server at: %s\n", config.Address)

	go http.ListenAndServe(config.Address, nil)

	log.Printf("http server started at: %s\n", server.config.Address)
	return nil
}
