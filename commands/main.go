package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
	"github.com/stevebargelt/buildwatcher/api"
	"github.com/stevebargelt/buildwatcher/ciserver"
	"github.com/stevebargelt/buildwatcher/controller"
	"github.com/stevebargelt/buildwatcher/slack"
)

//Version is the version... not implemented yet
var Version string

func main() {

	//create your file with desired read/write permissions
	f, err := os.OpenFile("buildwatcher.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	configFile := flag.String("config", "", "Build Watcher configuration file path")
	version := flag.Bool("version", false, "Print version information")
	flag.Usage = func() {
		text := `
    Usage: buildwatcher [OPTIONS]

    Options:

      -config string
          Configuration file path
      -version
			    Print version information
    `
		fmt.Println(strings.TrimSpace(text))
	}
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}
	config := &DefaultConfig
	if *configFile != "" {
		conf, err := ParseConfig(*configFile)
		if err != nil {
			log.Fatal("Failed to parse config file", err)
		}
		config = conf
	}

	// Initialize the controller
	c, err := controller.New(config.Controller)
	if err != nil {
		log.Fatal("Failed to initialize controller. ERROR:", err)
	}
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	// Initialize Jenkins
	//var jenkinsClient *gojenkins.Jenkins
	jenk := ciserver.NewJenkins()
	jenk := ciserver.NewJenkins(c, config.Jenkins)

	// Initialize the Slack controller
	sl := slack.NewSlack(c, config.Slack)
	go sl.StartSlack()
	log.Println("Starting Slack subsystem")

	//Initialize the API server
	if err := api.SetupServer(config.API, c); err != nil {
		log.Fatal("ERROR:", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGUSR2)
	for {
		select {
		case s := <-ch:
			switch s {
			case os.Interrupt:
				c.Stop()
				sl.Stop()
				return
				// case syscall.SIGUSR2:
				// 	c.DumpTelemetry()
			}
		}
	}
}

// Red    "12" //GPIO18
// Yellow "18" //GPIO24
// Green  "13" //GPIO27
// buzzer "16" //GPIO23
