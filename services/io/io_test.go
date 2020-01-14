package io

import (
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/stretchr/testify/assert"
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

func getIO() *IOModule {
	cfg := ConfigIO{
		SN:      "1",
		Model:   "MOXA_E1212",
		Address: "192.168.127.201:502",
	}
	return &IOModule{
		closing: make(chan struct{}, 1),
		cfg:     cfg,
		client: &ModbusTcp{
			cfg:    cfg,
			vendor: VendorModels["MOXA_E1212"],
		},

		flashInterval: 1 * time.Second,
	}
}

func TestStart(t *testing.T) {
	io := getIO()
	err := io.Start(nil)
	assert.Nil(t, err)
}

func TestStop(t *testing.T) {
	io := getIO()
	err := io.Stop()
	assert.Nil(t, err)

	io.Start(nil)
	err = io.Stop()
	assert.Nil(t, err)
}

func TestWrite(t *testing.T) {
	io := getIO()
	io.Start(nil)
	err := io.Write(0, OutputStatusOff)
	assert.NotNil(t, err)

	err = io.Write(0, OutputStatusOn)
	assert.NotNil(t, err)

	err = io.Write(0, OutputStatusFlash)
	assert.Nil(t, err)
}
