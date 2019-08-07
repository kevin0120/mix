package tightening_device

type DeviceConfig struct {
	Model    string `yaml:"model"`
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
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
