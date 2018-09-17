package changan

type Config struct {
	AndonIP   string `yaml:"andon_ip"`
	AndonPort uint   `yaml:"andon_port"`
	Enable    bool   `yaml:"enable"`
}

func NewConfig() Config {
	prs := []string{}
	prs = append(prs, "GSP")

	return Config{
		AndonIP:   "192.168.1.1",
		AndonPort: 8888,
		Enable:    true,
	}
}

func (c Config) Validate() error {

	return nil
}
