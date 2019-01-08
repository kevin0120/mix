package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
}

func newWsClient(id int, url string) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("create ws client failed: %d\n", id)
	}

	defer c.Close()

	str, _ := json.Marshal(WSRegist{
		HMI_SN: fmt.Sprintf("%d", id),
	})

	c.WriteMessage(websocket.TextMessage, str)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("read err\n")
			return
		}

		fmt.Printf("recv: %s\n", message)
	}
}

func main() {
	url := flag.String("url", "ws://127.0.0.1:9092/aiis/v1/ws", "--url")
	num := flag.Int("num", 100, "--num")

	flag.Parse()

	vUrl := *url
	vNum := *num

	for i := 0; i < vNum; i++ {
		go newWsClient(i+1, vUrl)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
