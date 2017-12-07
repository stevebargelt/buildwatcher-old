package main

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/stevebargelt/buildwatcher/api"
	"github.com/stevebargelt/buildwatcher/ciserver"
	"github.com/stevebargelt/buildwatcher/controller"
	"github.com/stevebargelt/buildwatcher/slack"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	API        api.ServerConfig  `yaml:"api"`
	Slack      slack.Config      `yaml:"slack"`
	CiServer   ciserver.Config   `yaml:"ciservers"`
}

func ParseConfig(filename string) (*Config, error) {

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))

	}
	
	fmt.Println(c)
	//fmt.Println(c.CiServer.CiServers[0].Name)
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
