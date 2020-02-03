package io

import (
	"fmt"
	"github.com/masami10/rush/services/wsnotify"
	"gopkg.in/resty.v1"
	"testing"
	"time"
)

func getHttpClient() *resty.Client {
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(3 * time.Second)
	client.SetContentLength(true)
	// Headers for all request
	client.
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	return client
}

func wsTest(msg *wsnotify.WSMsg) {
	http := getHttpClient()
	resp, err := http.R().SetBody(msg).Put("http://127.0.0.1:8082/rush/v1/ws-test")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp.StatusCode())
}

func ioSet(sn string, index uint16, status uint16) {
	wsTest(&wsnotify.WSMsg{
		SN:   0,
		Type: wsnotify.WS_IO_SET,
		Data: IoSet{
			SN:     sn,
			Index:  index,
			Status: status,
		},
	})
}

func Test_SetPSet(t *testing.T) {

	ioSet("2_io", 0, OutputStatusOn)
}
