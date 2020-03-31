package main

// CGO_ENABLED=0 GOARCH=amd64 go build ./stress_test.go

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var gAiisUrl = "http://10.1.1.31:8069"
var gTaskNum = 8
var gResultBatch = false

var gReqInteval = 200
var g_genetime_mtx = sync.Mutex{}

type Service struct {
	wg      sync.WaitGroup
	closing chan chan struct{}
}

var gService = Service{}

func Close() {
	for i := 0; i < gTaskNum; i++ {
		closed := make(chan struct{})
		gService.closing <- closed
		<-closed
	}
}

func task() {
	//nextWriteThreshold := time.Now()

	for {
		select {
		case <-time.After(time.Duration(gReqInteval) * time.Millisecond):
			//开始执行压力测试任务
			fmt.Println("stress test going\n")
		case stopDone := <-gService.closing:
			close(stopDone)
			fmt.Println("remove a new task for stress task")
			return //退出协程
		}
	}
}

func main() {

	fmt.Printf("start\n")

	aiis := flag.String("odoo", "http://127.0.0.1:9092", "--aiis")
	taskNum := flag.Int("masterpc", 10, "--masterpc")
	reqItv := flag.Int("interval", 2000, "--interval")
	batch := flag.Bool("batch", true, "--batch")
	flag.Parse()

	gAiisUrl = *aiis
	gTaskNum = *taskNum
	gReqInteval = *reqItv
	gResultBatch = *batch

	fmt.Println(fmt.Sprintf("odoo:%s, task:%d, inteval:%d", gAiisUrl, gTaskNum, gReqInteval))

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	for i := 0; i < gTaskNum; i++ {
		gService.wg.Add(1)
		go task()
		fmt.Println("add a new task for stress task")
	}

	for {
		for s := range signalCh {
			switch s.String() {
			case syscall.SIGTERM.String():
				fmt.Println("SIGTERM received, initializing clean shutdown...")
				go func() {
					Close()
				}()
				return

			case syscall.SIGHUP.String():
				fmt.Println("SIGHUP received, reloading tasks/templates/handlers directory...")
				go func() {
					Close()
				}()
				return

			default:
				fmt.Println("signal received, initializing clean shutdown...")
				go func() {
					Close()
				}()
				return
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}

}
