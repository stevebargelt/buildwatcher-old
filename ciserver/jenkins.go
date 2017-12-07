package ciserver

import (
	"log"
	"os"

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

//NewJenkins initializes the Jenkins client - connects to jenkins instance
func OldStartJenkins(jenkinsURL string, username string, password string) (*gojenkins.Jenkins, error) {
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, username, password).Init()
	return jenkins, err
}

func (j *Jenkins) StartJenkins() {
	j.stopCh = make(chan struct{})
	log.Println("Starting Jenkins")
	//test := j.config.CiServers[0].URL
	jenkins, err := gojenkins.CreateJenkins(j.config.CiServers[0].URL, j.config.CiServers[0].Username, j.config.CiServers[0].Password).Init()
	f, err := os.OpenFile("jenkins.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	job, _ := jenkins.GetJob(j.config.CiServers[0].Jobs[0].Name)
	job.Poll()

	//logger := log.New(f, "jenkins: ", log.Lshortfile|log.LstdFlags)

	// slackAPI.SetDebug(true)

	// rtm := slackAPI.NewRTM()

	//go rtm.ManageConnection()

	for {
		select {
		//case msg := <-rtm.IncomingEvents:

		case <-j.stopCh:
			log.Println("Stopping Jenkins polling")
			return
		}
	}

}
