package scanner

type DeviceConfig struct {
	VendorID  int `yaml:"vendorID"`
	ProductID int `yaml:"productID"`
}

type Config struct {
	Enable  bool           `yaml:"enable"`
	Devices []DeviceConfig `yaml:"devices"`
}

func NewConfig() Config {
	return Config{
		Enable:  true,
		Devices: []DeviceConfig{{1529, 8714}},
	}
}

func (c Config) Validate() error {

	return nil
}
