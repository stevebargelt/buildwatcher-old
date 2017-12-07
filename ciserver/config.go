package ciserver

// Config is a generic structure for Ci Server configs
type Config struct {
	CiServers []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		Pollrate int    `yaml:"pollrate"`
		Jobs     []struct {
			Name   string `yaml:"name"`
			Branch string `yaml:"branch"`
		} `yaml:"jobs"`
	} `yaml:"ciservers"`
}

// "light": {
//     "type": "console",
//     "num_leds": 32
// },
