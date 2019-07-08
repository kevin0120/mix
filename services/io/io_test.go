package io

import (
	"fmt"
	"github.com/goburrow/modbus"
	"testing"
	"time"
)

func TestIO(t *testing.T) {

	handler := modbus.NewTCPClientHandler("192.168.127.254:502")
	handler.Timeout = 3 * time.Second

	err := handler.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	client.WriteSingleCoil(0, 0xFF00)
	//client.write
	for {
		//outputs, err := client.ReadCoils(0, 8)
		//if err == nil {
		//	fmt.Println(outputs)
		//} else {
		//	break
		//}

		inputs, err := client.ReadDiscreteInputs(0, 8)
		if err == nil {
			fmt.Println(inputs)
		} else {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	//_, err = client.WriteMultipleCoils(0, 10, []byte{4, 3})

}
