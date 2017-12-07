package ciserver

import (
	"log"
	"os"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/stevebargelt/buildwatcher/controller"
)

var _STATUS = map[string]Status{
	"aborted":        ABORTED,
	"aborted_anime":  BUILDING_FROM_ABORTED,
	"blue":           SUCCESS,
	"blue_anime":     BUILDING_FROM_SUCCESS,
	"disabled":       DISABLED,
	"disabled_anime": BUILDING_FROM_DISABLED,
	"grey":           UNKNOWN,
	"grey_anime":     BUILDING_FROM_UNKNOWN,
	"notbuilt":       NOT_BUILT,
	"notbuilt_anime": BUILDING_FROM_NOT_BUILT,
	"red":            FAILURE,
	"red_anime":      BUILDING_FROM_FAILURE,
	"yellow":         UNSTABLE,
	"yellow_anime":   BUILDING_FROM_UNSTABLE,
}

type Jenkins struct {
	stopCh chan struct{}
	config Config
	c      *controller.Controller
}

func NewJenkins(c *controller.Controller, jenkinsConfig Config) *Jenkins {
	return &Jenkins{
		config: jenkinsConfig,
		c:      c,
	}
}

func (j *Jenkins) StartJenkins() {
	j.stopCh = make(chan struct{})
	log.Println("Starting Jenkins")
	//test := j.config.CiServers[0].URL
	jenkins, _ := gojenkins.CreateJenkins(j.config.CiServers[0].URL, j.config.CiServers[0].Username, j.config.CiServers[0].Password).Init()
	f, err := os.OpenFile("jenkins.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	job, _ := jenkins.GetJob(j.config.CiServers[0].Jobs[0].Name)

	ticker := time.NewTicker(time.Millisecond * 10000)
	go func() {
		for t := range ticker.C {
			job.Poll()
			status := _STATUS[job.GetDetails().Color]
			log.Println("Status = ", status)
			log.Println("Tick at = ", t)
		}
	}()

	select {
	case <-j.stopCh:
		log.Println("Stopping Slack polling")
		return
	}

}

func (j *Jenkins) Stop() {
	if j.stopCh == nil {
		log.Println("WARNING: stop channel is not initialized.")
		return
	}
	j.stopCh <- struct{}{}
	log.Println("Stopped Jenkins")
}
