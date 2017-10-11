package controller

type Config struct {
	EnableGPIO   bool             `yaml:"enable_gpio"`
	Database     string           `yaml:"database"`
	HighRelay    bool             `yaml:"high_relay"`
	Lights       map[string]Light `yaml:"lights"`
	SlackToken   string           `yaml:"slack_token"`
	SlackChannel string           `yaml:"slack_channel"`
	DevMode      bool             `yaml:"dev_mode"`
}

var DefaultConfig = Config{
	Database:   "buildwatcher.db",
	EnableGPIO: true,
	Lights:     make(map[string]Light),
}
