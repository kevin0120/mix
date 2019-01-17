package main

import (
	"fmt"
)

func main() {
	//conf_path := flag.String("conf", "/PMON.CFG", "--conf")
	//sleep := flag.Int("sleep", 1, "--sleep")
	//flag.Parse()
	//var v_conf = *conf_path
	//var v_sleep = *sleep
	//
	//cfg := pmon.NewConfig()
	//cfg.Path = v_conf
	//fis, err := pmon.NewService(cfg, nil)
	//if err != nil {
	//	fmt.Printf("fis open err: %s\n", err.Error())
	//	return
	//}
	//
	//err = fis.Open()
	//if err != nil {
	//	fmt.Printf("fis open err: %s\n", err.Error())
	//	return
	//}
	//
	//for {
	//	data := fmt.Sprintf("SR1J--V001--C6-2018-1000001=9-01000001**LSV2A8CA7JN990001*BR24J3G0C 4X0 U80 Q89 1AS T1P 0EX 0MX D60 8RM 3KP 4UH UH2 4A4     7GR 26.04.2018*08.38.15**")
	//	err = fis.SendPmonMessage(pmon.PMONMSGSD, "21", data)
	//	if err != nil {
	//		fmt.Printf("send err: %s\n", err.Error())
	//		return
	//	}
	//
	//	if v_sleep > 0 {
	//		time.Sleep(time.Duration(v_sleep) * time.Millisecond)
	//	}
	//}

	sd := "\x02002000050093SD00020222000222SR3J--V001--C6-2019-9010001=4-01000220**ERGEBNIS00200182##001*RESULT*_VW416_1-62-7199-03 *Screwer   *1.0   *IO__*20190117180147*AUTO*VALUE *002*0000010.039               Nm        1*000002360.6               DEG       1***\x10\x039"
	//sdidx := strings.Index(sd, "SD")
	//bcStart := sdidx+len("SD")
	//newBC := fmt.Sprintf("%04d", 1)
	//for i := bcStart; i < bcStart+4 ; i++ {
	//
	//}

	//arr := strings.Split(sd, "SD")
	//new := arr[0] + "SD" + fmt.Sprintf("%04d", 1) + arr[1][4:]

	fmt.Printf("%s", sd)
}
