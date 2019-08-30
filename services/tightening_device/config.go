package tightening_device

type DeviceConfig struct {
	Model    string       `yaml:"model"`
	Protocol string       `yaml:"protocol"`
	Endpoint string       `yaml:"endpoint"`
	SN       string       `yaml:"sn"`
	Tools    []ToolConfig `yaml:"tools"`
}

type ToolConfig struct {
	SN      string `yaml:"sn"`
	Channel int    `yaml:"channel"`
}

type Config struct {
	Enable  bool           `yaml:"enable"`
	Devices []DeviceConfig `yaml:"devices"`
}

func NewConfig() Config {

	return Config{
		Enable:  true,
		Devices: []DeviceConfig{},
	}
}

func (c Config) Validate() error {

	return nil
}
