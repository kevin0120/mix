package main

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"github.com/masami10/rush/storage"
	"github.com/masami10/rush/db"
	"github.com/masami10/rush/core"
	"github.com/masami10/rush/conf"
)



func main() {

	// 初始化配置文件
	c, err := conf.InitConf()
	if err != nil {
		fmt.Printf("初始化配置文件失败\n")
		return
	}

	//var err error

	db := rushdb.DB{}
	storage := rush_storage.Storage{}
	odoo := core.ODOO{}
	cvi3_service := core.CVI3Service{}
	API := core.APIServer{}

	fmt.Printf("初始化rush\n")
	masterpc_sn := c.MasterPC.SN

	fmt.Printf("初始化数据库\n")
	db = c.DB

	fmt.Printf("初始化对象存储\n")
	storage.Conf = &c.MINIO

	fmt.Printf("初始化odoo服务\n")
	odoo.URL = c.ODOO.Urls[0]
	odoo.DB = &db
	odoo.MasterPC_SN = masterpc_sn
	odoo.APIService = &API

	fmt.Printf("初始化cvi3服务\n")

	cvi3_service.DB = &db
	cvi3_service.Storage = &storage
	cvi3_service.Port = fmt.Sprintf(":%d", c.CVI3.Listen)
	cvi3_service.ODOO = &odoo
	cvi3_service.APIService = &API

	configs := []cvi3.CVI3Config{}
	configs = append(configs, cvi3.CVI3Config{c.CVI3.Controllers[0].SN, c.CVI3.Controllers[0].IP, c.CVI3.Controllers[0].Port})
	cvi3_service.Config(configs)

	err = cvi3_service.StartService()
	if err != nil {
		fmt.Printf("cvi3_services service error:%s\n", err.Error())
	}

	fmt.Printf("初始化api服务\n")

	API.DB = &db
	API.Port = fmt.Sprintf(":%d", c.MasterPC.Port)
	API.CVI3 = &cvi3_service
	err = API.StartService()
	if err != nil {
		fmt.Printf("api service error:%s\n", err.Error())
	}
}
