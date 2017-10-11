package slack

type Config struct {
	SlackToken   string `yaml:"slack_token"`
	SlackChannel string `yaml:"slack_channel"`
}

var DefaultConfig = Config{
	SlackChannel: "#general",
}
