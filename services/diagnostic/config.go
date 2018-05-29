package diagnostic


type Config struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

func NewConfig() Config {
	return Config{
		File:  "STDERR",
		Level: "DEBUG",
	}
}
