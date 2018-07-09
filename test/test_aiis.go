package main

import (
	"github.com/masami10/aiis/services/pmon"
	"fmt"
	"flag"
	"time"
)

func main() {
	conf_path := flag.String("conf", "/PMON.CFG", "--conf")
	sleep := flag.Int("sleep", 1, "--sleep")
	flag.Parse()
	var v_conf = *conf_path
	var v_sleep = *sleep

	cfg := pmon.NewConfig()
	cfg.Path = v_conf
	fis, err := pmon.NewService(cfg, nil)
	if err != nil {
		fmt.Printf("fis open err: %s\n", err.Error())
		return
	}

	err = fis.Open()
	if err != nil {
		fmt.Printf("fis open err: %s\n", err.Error())
		return
	}

	for {
		data := fmt.Sprintf("SR1J--V001--C6-2018-1000001=9-01000001**LSV2A8CA7JN990001*BR24J3G0C 4X0 U80 Q89 1AS T1P 0EX 0MX D60 8RM 3KP 4UH UH2 4A4     7GR 26.04.2018*08.38.15**")
		err = fis.SendPmonMessage(pmon.PMONMSGSD, "21", data)
		if err != nil {
			fmt.Printf("send err: %s\n", err.Error())
			return
		}

		if v_sleep > 0 {
			time.Sleep(time.Duration(v_sleep) * time.Millisecond)
		}
	}
}