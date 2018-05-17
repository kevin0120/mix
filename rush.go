package main

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"github.com/masami10/rush/storage"
	"github.com/masami10/rush/db"
	"github.com/masami10/rush/core"
)


func main() {

	var err error

	db := rushdb.DB{}
	storage := rush_storage.Storage{}
	odoo := core.ODOO{}
	cvi3_service := core.CVI3Service{}
	api := core.APIServer{}

	fmt.Printf("初始化rush\n")
	masterpc_sn := "1qaz2wsx"

	fmt.Printf("初始化数据库\n")

	db.URL = "127.0.0.1"
	db.DBName = "rush"
	db.Port = 5432
	db.Pwd = "admin"
	db.User = "admin"

	fmt.Printf("初始化对象存储\n")

	storage.MinioURL = "127.0.0.1:39000"
	storage.MinioBacket = "curves"
	storage.MinioAccess = "TC1O0SUW08VRFBFYWK9C"
	storage.MinioSecret = "kAKzVLDEHeFI0xlBjkeVFKfkjPd/9JlBolz13UC8"

	fmt.Printf("初始化odoo服务\n")

	odoo.URL = "http://127.0.0.1:8069"
	odoo.DB = &db
	odoo.MasterPC_SN = masterpc_sn
	odoo.APIService = &api

	fmt.Printf("初始化cvi3服务\n")

	cvi3_service.DB = &db
	cvi3_service.Storage = &storage
	cvi3_service.Port = ":4710"
	cvi3_service.ODOO = &odoo
	cvi3_service.APIService = &api

	configs := []cvi3.CVI3Config{}
	configs = append(configs, cvi3.CVI3Config{"1", "192.168.1.200", 4700})
	cvi3_service.Config(configs)

	err = cvi3_service.StartService()
	if err != nil {
		fmt.Printf("cvi3_services service error:%s\n", err.Error())
	}

	fmt.Printf("odoo工单请求服务\n")
	go odoo.RequestWorkerOrders()

	fmt.Printf("初始化api服务\n")

	api.DB = &db
	api.Port = ":8080"
	api.CVI3 = &cvi3_service
	err = api.StartService()
	if err != nil {
		fmt.Printf("api service error:%s\n", err.Error())
	}
}
