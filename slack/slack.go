package slack

import (
	"fmt"
	"log"
	"os"
	"regexp"

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
	fmt.Printf("TOKEN: %s\n\n", s.config.SlackToken)
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
			//fmt.Print("Event Received:\n\n")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
			case *slack.ConnectedEvent:
				log.Printf("Connected: %#v %v\n", ev.Info.User, ev.Info.Channels)
				log.Printf("Connection counter: %v\n\n", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage("Build Watcher Connected", s.config.SlackChannel))
			case *slack.MessageEvent:
				parseMessage(ev)
				// if Match(ev, "(?i)green(.*)") {
				// 	s.c.LightOn("green")
				// }
				// if Match(ev, "(?i)red(.*)") {
				// 	s.c.LightOn("red")
				// }
			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				log.Printf("SLACK: Invalid credentials")
				break
			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
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

func parseMessage(ev *slack.MessageEvent) {

	var statusRegEx = regexp.MustCompile(`(?i)status=\w*`)

	message := fmt.Sprintf("%v", ev.Msg)
	if statusRegEx.MatchString(message) {
		fmt.Printf("Found Status: %v\n\n\n", statusRegEx.FindString(message))
	}
	// fmt.Println("\n\nALL THE INFO FROM parseMessage")
	// fmt.Printf("ev.Name: %s\n", ev.Msg.Name)
	// fmt.Printf("ev.BotID: %s\n", ev.Msg.BotID)
	// fmt.Printf("ev.User: %s\n", ev.Msg.User)
	// fmt.Printf("ev.Username: %s\n", ev.Msg.Username)
	// fmt.Printf("ev.Type: %s\n", ev.Msg.Type)
	// fmt.Printf("ev.Msg.Text: %s\n", ev.Msg.Text)
	// fmt.Printf("ev.Msg.Topic: %s\n", ev.Msg.Topic)
	// fmt.Printf("ev.Msg.Comment: %s\n", ev.Msg.Comment)
	// fmt.Printf("ev.Msg.Attachments: %v\n\n", ev.Msg.Attachments)
	// if len(ev.Msg.Attachments) > 0 {
	// 	fmt.Printf("ev.Msg.Attachments.Text: %v\n\n", ev.Msg.Attachments[0].Text)

	// 	if len(ev.Msg.Attachments[0].Fields) > 0 {
	// 		fmt.Printf("Breaking it down...\n")

	// 		for k, v := range ev.Msg.Attachments[0].Fields {
	// 			fmt.Printf("ev.Msg.Attachments[0].Field[%v]: %v\n", k, v)
	// 			fmt.Printf("ev.Msg.Attachments[0].Field[%v].Value: %v\n", k, v.Value)
	// 		}

	// 		m := make(map[string]string)
	// 		text := ev.Msg.Attachments[0].Fields[0].Value
	// 		fields := strings.Split(text, ",")

	// 		for _, pair := range fields {
	// 			z := strings.Split(pair, "=")
	// 			z[0] = strings.TrimSuffix(z[0], ", ")
	// 			z[1] = strings.TrimSpace(z[1])
	// 			m[z[0]] = z[1]
	// 		}

	// 		fmt.Printf("\n\nThe MAP\n")
	// 		for a, b := range m {
	// 			fmt.Printf("a=%v, b=%v\n", a, b)
	// 		}
	// 	}
	// }

	// Jenkins
	// ev.Name:
	// ev.BotID: B7FDY3PJM
	// ev.User:
	// ev.Username:
	// ev.Type: message

	// Travis
	// ev.Name:
	// ev.BotID: B74V30W7J
	// ev.User:
	// ev.Username:
	// ev.Type: message

}
