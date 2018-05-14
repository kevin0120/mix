package main

import (
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"fmt"
)

func OnStatus(sn string, status string) {
	fmt.Printf("%s:%s", sn, status)
}

func OnRecv(msg string) {
	fmt.Printf("%s", msg)
}

func main() {
	cvi3_service := cvi3.CVI3{}

	configs := []cvi3.CVI3Config{}
	configs = append(configs, cvi3.CVI3Config{"1", "192.168.1.200", 4700})

	cvi3_service.Config(configs)
	cvi3_service.RegisterCallBack(OnStatus, OnRecv)

	cvi3_service.StartService(4710)
}
