package main

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"github.com/masami10/rush/storage"
	"github.com/masami10/rush/db"
	"github.com/masami10/rush/core"
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
)



type conf_rush struct {
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

func main() {

	conf_path := "./conf/rush.yml"
	exist, err := file_exists(conf_path)
	var conf conf_rush = conf_rush{}
	if err != nil || !exist {
		// 创建缺省配置文件
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
		controller.Port = 4710
		controller.SN = "1"
		conf.CVI3.Controllers = append(conf.CVI3.Controllers, controller)

		conf.MasterPC.SN = "1"
		conf.MasterPC.Port = 8080

		conf.MINIO.URL = "127.0.0.1:9000"
		conf.MINIO.Backet = "backet"
		conf.MINIO.Access = "access"
		conf.MINIO.Secret = "secret"

		conf_s, err := yaml.Marshal(conf)
		if err != nil {
			fmt.Printf("无法创建配置文件: %s\n", err.Error())
			return
		} else {
			os.Mkdir("./conf", os.ModePerm)
			f, _ := os.Create(conf_path)
			f.Close()
			err := ioutil.WriteFile(conf_path, conf_s, 0666)
			if err != nil {
				fmt.Printf("无法创建配置文件: %s\n", err.Error())
				return
			}
		}
	} else {
		// 读取配置文件
		yml, err := ioutil.ReadFile(conf_path)
		if err != nil {
			fmt.Printf("无法读取配置文件: %s\n", err.Error())
			return
		} else {
			err := yaml.Unmarshal(yml, &conf)
			if err != nil {
				fmt.Printf("无法读取配置文件: %s\n", err.Error())
				return
			}
		}
	}

	//var err error

	db := rushdb.DB{}
	storage := rush_storage.Storage{}
	odoo := core.ODOO{}
	cvi3_service := core.CVI3Service{}
	API := core.APIServer{}

	fmt.Printf("初始化rush\n")
	masterpc_sn := conf.MasterPC.SN

	fmt.Printf("初始化数据库\n")
	db = conf.DB

	fmt.Printf("初始化对象存储\n")
	storage.Conf = &conf.MINIO

	fmt.Printf("初始化odoo服务\n")

	odoo.URL = "http://10.1.1.31:8069"
	//odoo.URL = "http://127.0.0.1:8069"
	odoo.DB = &db
	odoo.MasterPC_SN = masterpc_sn
	odoo.APIService = &API

	fmt.Printf("初始化cvi3服务\n")

	cvi3_service.DB = &db
	cvi3_service.Storage = &storage
	cvi3_service.Port = ":4710"
	cvi3_service.ODOO = &odoo
	cvi3_service.APIService = &API

	configs := []cvi3.CVI3Config{}
	configs = append(configs, cvi3.CVI3Config{"1", "192.168.1.200", 4700})
	cvi3_service.Config(configs)

	err = cvi3_service.StartService()
	if err != nil {
		fmt.Printf("cvi3_services service error:%s\n", err.Error())
	}

	fmt.Printf("初始化api服务\n")

	API.DB = &db
	API.Port = ":" + string(conf.MasterPC.Port)
	API.CVI3 = &cvi3_service
	err = API.StartService()
	if err != nil {
		fmt.Printf("api service error:%s\n", err.Error())
	}
}
