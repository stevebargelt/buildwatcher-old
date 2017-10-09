package controller

type Config struct {
	EnableGPIO bool             `yaml:"enable_gpio"`
	HighRelay  bool             `yaml:"high_relay"`
	Database   string           `yaml:"database"`
	Lights     map[string]Light `yaml:"lights"`
	DevMode    bool             `yaml:"dev_mode"`
}

var DefaultConfig = Config{
	Database:   "buildwatcher.db",
	EnableGPIO: true,
	Lights:     make(map[string]Light),
}
