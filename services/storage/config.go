package storage

type Config struct {
	Enable bool `yaml:"enable"`
	Url         string `yaml:"db_url"`
	DBName      string `yaml:"db_name"`
	User        string `yaml:"db_user"`
	Password    string `yaml:"db_pwd"`
	MaxConnects int    `yaml:"max_connection"`
}

func NewConfig() Config {
	return Config{
		Enable: false,
		Url:         "127.0.0.1:5432",
		DBName:      "dbname",
		User:        "user",
		Password:    "pwd",
		MaxConnects: 60,
	}
}

func (c Config) Validate() error {

	return nil
}
