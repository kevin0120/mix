package fis

type Config struct {
	chRecv          string `yaml:"ch_recv"`
	chSend          string `yaml:"ch_send"`
	SystemType      string `yaml:"system_type"`
	SoftwareVersion string `yaml:"software_version"`
	Mode            string `yaml:"mode"`
}

func NewConfig() Config {
	return Config{
		chSend:          "02",
		chRecv:          "01",
		SystemType:      "screw",
		SoftwareVersion: "1.0",
		Mode:            "AUTO",
	}
}

func (c Config) Validate() error {

	return nil
}
