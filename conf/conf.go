package conf

import (
	"github.com/masami10/rush/db"
	"github.com/masami10/rush/storage"
	"github.com/masami10/rush/core"
	"os"
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
)

const (
	DEFAULT_CONF_PATH = "./conf/rush.yml"
)

type ConfRush struct {
	MasterPC struct{
		SN string	`yaml:"sn"`
		Port int	`yaml:"api_port"`
	} `yaml:"masterpc"`

	DB rushdb.DB `yaml:"db"`
	MINIO rush_storage.StorageConf `yaml:"minio"`
	ODOO core.ODOOConf `yaml:"odoo"`
	CVI3 core.CVI3Conf `yaml:"cvi3"`
}

func file_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func CreateDefaultConf() ConfRush {
	conf := ConfRush{}
	conf.DB.URL = "127.0.0.1"
	conf.DB.User = "user"
	conf.DB.Pwd = "pwd"
	conf.DB.DBName = "dbname"
	conf.DB.Port = 5432

	conf.ODOO.Timeout = 3000
	conf.ODOO.Urls = append(conf.ODOO.Urls, "http://127.0.0.1:8069")

	conf.CVI3.Listen = 4710
	controller := core.ControllerConf{}
	controller.IP = "127.0.0.1"
	controller.Port = 4700
	controller.SN = "1"
	conf.CVI3.Controllers = append(conf.CVI3.Controllers, controller)

	conf.MasterPC.SN = "1"
	conf.MasterPC.Port = 8080

	conf.MINIO.URL = "127.0.0.1:9000"
	conf.MINIO.Backet = "backet"
	conf.MINIO.Access = "access"
	conf.MINIO.Secret = "secret"

	return conf
}

func InitConf() (ConfRush, error) {
	exist, err := file_exists(DEFAULT_CONF_PATH)
	var conf ConfRush = ConfRush{}

	if err != nil || !exist {
		// 创建缺省配置文件
		conf = CreateDefaultConf()
		conf_s, err := yaml.Marshal(conf)
		if err != nil {
			fmt.Printf("无法创建配置文件: %s\n", err.Error())
			return conf, err
		} else {
			os.Mkdir("./conf", os.ModePerm)
			f, _ := os.Create(DEFAULT_CONF_PATH)
			f.Close()
			err := ioutil.WriteFile(DEFAULT_CONF_PATH, conf_s, 0666)
			if err != nil {
				fmt.Printf("无法创建配置文件: %s\n", err.Error())
				return conf, err
			}
		}
	} else {
		// 读取配置文件
		yml, err := ioutil.ReadFile(DEFAULT_CONF_PATH)
		if err != nil {
			fmt.Printf("无法读取配置文件: %s\n", err.Error())
			return conf, err
		} else {
			err := yaml.Unmarshal(yml, &conf)
			if err != nil {
				fmt.Printf("无法读取配置文件: %s\n", err.Error())
				return conf, err
			}
		}
	}

	return conf, nil
}