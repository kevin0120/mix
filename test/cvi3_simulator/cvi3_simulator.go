package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/masami10/rush/test/cvi3_simulator/controller"
	"gopkg.in/yaml.v2"
	"time"
)

const (
	PORT_START = 4720
	SN_START   = 1
)

type cvi3_cfg struct {
	SN       string `yaml:"serial_no"`
	Protocol string `yaml:"protocol"`
	IP       string `yaml:"remote_ip"`
	Port     int    `yaml:"port"`
}

func main() {

	num := flag.Int("controller", 100, "--controller")
	ip := flag.String("ip", "127.0.0.1", "--ip")
	workorder := flag.Int("workorder", 8, "--workorder")
	results := flag.String("results", "[28, 29, 30]", "--results")
	count := flag.Int("count", 3, "--count")

	flag.Parse()

	var v_workorder = *workorder
	var v_results = *results
	var v_count = *count

	var result_ids = []int64{}
	json.Unmarshal([]byte(v_results), &result_ids)

	var v_ip string = *ip
	var v_num int = *num
	cfgs := []cvi3_cfg{}

	servers := []*controller.CVI3Server{}

	var sn = SN_START
	for i := PORT_START; i < (PORT_START + v_num); i++ {
		cfg := cvi3_cfg{
			SN:       fmt.Sprintf("%04d", sn),
			Protocol: "Audi/VW",
			IP:       v_ip,
			Port:     i,
		}

		cfgs = append(cfgs, cfg)

		server := controller.CVI3Server{
			SN:           cfg.SN,
			Workorder_id: int64(v_workorder),
			Result_ids:   result_ids,
			Count:        v_count,
		}

		servers = append(servers, &server)
		server.Start(uint(i))

		sn++
	}

	s, _ := yaml.Marshal(cfgs)
	fmt.Printf("%s\n", s)

	for {
		time.Sleep(1 * time.Second)
	}
}
