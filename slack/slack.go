package slack

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/stevebargelt/buildwatcher/controller"
)

const SlackBucket = "slack"

type Slack struct {
	stopCh chan struct{}
	config Config
	c      *controller.Controller
}

func NewSlack(c *controller.Controller, slackConfig Config) *Slack {
	return &Slack{
		config: slackConfig,
		c:      c,
	}
}

func (s *Slack) StartSlack() {
	s.stopCh = make(chan struct{})
	log.Println("Starting Slack")
	slackAPI := slack.New(s.config.SlackToken)
	f, err := os.OpenFile("slackapi.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger := log.New(f, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	slackAPI.SetDebug(true)

	rtm := slackAPI.NewRTM()

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
			case *slack.ConnectedEvent:
				log.Printf("Connected: %#v %v\n", ev.Info.User, ev.Info.Channels)
				log.Printf("Connection counter: %v\n\n", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage("Build Watcher Connected", s.config.SlackChannel))
			case *slack.MessageEvent:
				s.parseMessage(ev)
			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				log.Printf("SLACK: Invalid credentials")
				break
			default:
				// Ignore other events..
			}
		case <-s.stopCh:
			log.Println("Stopping Slack listener")
			return
		}
	}

}

func (s *Slack) Stop() {
	if s.stopCh == nil {
		log.Println("WARNING: stop channel is not initialized.")
		return
	}
	s.stopCh <- struct{}{}
	log.Println("Stopped Slack")
}

func (s *Slack) parseMessage(ev *slack.MessageEvent) {

	greenStatus := map[string]bool{
		"passed":    true,
		"success":   true,
		"succeeded": true,
	}
	yellowStatus := map[string]bool{
		"building": true,
		"started":  true,
	}
	redStatus := map[string]bool{
		"failed":  true,
		"failure": true,
		"failing": true,
	}

	var statusRegEx = `(?i)status=\w*`
	var jobRegex = `(?i)job=\w*`
	var buildRegex = `(?i)build=\w*`
	var URLRegex = `(?i)url=(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`

	message := fmt.Sprintf("%v", ev.Msg)
	status := matcher(message, statusRegEx)
	job := matcher(message, jobRegex)
	build := matcher(message, buildRegex)
	URL := matcher(message, URLRegex)

	fmt.Printf("\n\nstatus:%v,job:%v,build:%v,url:%v\n\n", status, job, build, URL)

	// proj, err := s.c.GetProject(job)
	// if err != nil {
	// 	log.Printf("Slack.go: Error trying to get project")
	// 	panic(err)
	// }

	// fmt.Printf("Project:%v|EOL", proj)

	if greenStatus[strings.ToLower(status)] {
		fmt.Printf("Green, baby\n")
		//s.c.LightOn("green")
	}

	if redStatus[strings.ToLower(status)] {
		fmt.Printf("Red, oh no\n")
		//s.c.LightOn("red")
	}

	if yellowStatus[strings.ToLower(status)] {
		fmt.Printf("Yellow, fingers crossed\n")
		//s.c.LightOn("yellow")
	}

	// Jenkins
	// ev.BotID: B7FDY3PJM

	// Travis
	// ev.BotID: B74V30W7J

}

func matcher(message string, match string) string {

	var kv []string

	theregex := regexp.MustCompile(match)
	if theregex.MatchString(message) {
		found := theregex.FindString(message)
		kv = strings.Split(found, "=")
		fmt.Printf("The key is: %v\n", kv[0])
		fmt.Printf("The value is: %v\n", kv[1])
	} else {
		fmt.Printf("No match Found\n")
		return ""
	}
	return kv[1]
}
