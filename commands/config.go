package main

import (
	"io/ioutil"

	"github.com/kidoman/embd"
	"github.com/stevebargelt/buildwatcher/api"
	"github.com/stevebargelt/buildwatcher/controller"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	API        api.ServerConfig  `yaml:"api"`
}

var DefaultConfig = Config{
	Controller: controller.DefaultConfig,
	API:        api.DefaultConfig,
}

func ParseConfig(filename string) (*Config, error) {
	var c Config
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	for k, l := range c.Controller.Lights {
		l.ID = k
		c.Controller.Lights[l.ID] = l
		l.Dpin, err = embd.NewDigitalPin(l.GPIO)
		if err != nil {
			panic(err)
		}

	}
	return &c, nil
}
