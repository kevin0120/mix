package main

import (
	"time"
	"github.com/masami10/rush/test/cvi3_simulator/controller"
	"flag"
)

const (
	DEFAULT_PORT = 4720
)

func main() {

	num := flag.Int("controller", 10, "--controller")
	flag.Parse()

	var v_num int = *num

	for i := DEFAULT_PORT; i < (DEFAULT_PORT + v_num); i++ {
		server := controller.CVI3Server{}
		server.Start(uint(i))
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
