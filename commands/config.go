package main

import (
	"io/ioutil"

	"github.com/stevebargelt/buildwatcher/api"
	"github.com/stevebargelt/buildwatcher/controller"
	"github.com/stevebargelt/buildwatcher/slack"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	API        api.ServerConfig  `yaml:"api"`
	Slack      slack.Config      `yaml:"slack"`
}

var DefaultConfig = Config{
	Controller: controller.DefaultConfig,
	API:        api.DefaultConfig,
	Slack:      slack.DefaultConfig,
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

	// // add the embd digital pins for each light that is configured
	// for _, l := range c.Controller.Lights {
	// 	// l.ID = k
	// 	// c.Controller.Lights[l.ID] = l
	// 	l.Dpin, err = embd.NewDigitalPin(l.GPIO)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// }
	return &c, nil
}
