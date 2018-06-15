package fis

type Config struct {
	CH_RECV         string `yaml:"ch_recv"`
	CH_SEND         string `yaml:"ch_send"`
	SystemType      string `yaml:"system_type"`
	SoftwareVersion string `yaml:"software_version"`
	Mode            string `yaml:"mode"`
}

func NewConfig() Config {
	return Config{
		CH_SEND:         "02",
		CH_RECV:         "01",
		SystemType:      "screw",
		SoftwareVersion: "1.0",
		Mode:            "AUTO",
	}
}

func (c Config) Validate() error {

	return nil
}
