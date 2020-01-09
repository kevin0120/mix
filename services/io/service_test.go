package io

import (
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getSrv() *Service {
	srv := NewService(NewConfig(), nil)
	srv.WS = wsnotify.NewService(wsnotify.NewConfig(), nil)
	srv.deviceService, _ = device.NewService(device.NewConfig(), nil)

	return srv
}

func Test_Open(t *testing.T) {
	srv := getSrv()
	err := srv.Open()
	assert.Nil(t, err)
}

func Test_Close(t *testing.T) {
	srv := getSrv()
	err := srv.Close()
	assert.Nil(t, err)
}

func Test_Read(t *testing.T) {
	srv := getSrv()
	_, _, err := srv.Read("1")
	assert.NotNil(t, err)
}

func Test_Write(t *testing.T) {
	srv := getSrv()
	err := srv.Write("1", 0, 0)
	assert.NotNil(t, err)
}
