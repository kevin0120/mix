package tightening_device

type TighteningDeviceConfig struct {
	// 控制器型号
	Model string `yaml:"model"`

	// 控制器协议类型
	Protocol string `yaml:"protocol"`

	// 连接地址
	Endpoint string `yaml:"endpoint"`

	// 控制器序列号
	SN string `yaml:"sn"`

	// 工具列表
	Tools []ToolConfig `yaml:"tools"`
}

type ToolConfig struct {
	// 工具序列号
	SN string `yaml:"sn"`

	// 工具通道号
	Channel int `yaml:"channel"`
}

type Config struct {
	Enable  bool                     `yaml:"enable"`
	Devices []TighteningDeviceConfig `yaml:"devices"`
}

func NewConfig() Config {

	return Config{
		Enable:  true,
		Devices: []TighteningDeviceConfig{},
	}
}

func (c Config) Validate() error {

	return nil
}
