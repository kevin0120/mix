package io

import "github.com/masami10/rush/services/device"

const (
	IO_STATUS_ONLINE  = device.BaseDeviceStatusOnline
	IO_STATUS_OFFLINE = device.BaseDeviceStatusOffline

	IO_TYPE_INPUT  = "input"
	IO_TYPE_OUTPUT = "output"

	IO_MODBUSTCP = "modbustcp"
)
