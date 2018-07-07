package fis

type Config struct {
	CHRecv          string 		`yaml:"ch_recv"`
	CHSend          string 		`yaml:"ch_send"`
	SystemType      string 		`yaml:"system_type"`
	SoftwareVersion string 		`yaml:"software_version"`
	Mode            string 		`yaml:"mode"`
	FactoryCode		string 		`yaml:"factory_code"`
	PRS				[]string	`yaml:"prs"`
}

func NewConfig() Config {
	prs := []string{}
	prs = append(prs, "GSP")

	return Config{
		CHSend:          "02",
		CHRecv:          "01",
		SystemType:      "screw",
		SoftwareVersion: "1.0",
		Mode:            "AUTO",
		FactoryCode:	 "01",
		PRS:			 prs,
	}
}

func (c Config) Validate() error {

	return nil
}
