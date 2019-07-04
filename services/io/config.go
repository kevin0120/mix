package io

type ConfigIO struct {
	SN      string `yaml:"sn"`
	Model   string `yaml:"model"`
	Address string `yaml:"address"`
}

type Config struct {
	Enable bool       `yaml:"enable"`
	IOS    []ConfigIO `yaml:"ios"`
}

func NewConfig() Config {
	return Config{
		Enable: true,
		IOS: []ConfigIO{
			{
				SN:      "1",
				Model:   MODEL_MOXA_E1212,
				Address: "modbustcp://127.0.0.1:502",
			},
		},
	}
}

func (c Config) Validate() error {

	return nil
}
