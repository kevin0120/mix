package io

const (
	WS_IO_STATUS  = "WS_IO_STATUS"
	WS_IO_CONTACT = "WS_IO_CONTACT"
)

type IO_STATUS struct {
	SN     string `json:"sn"`
	Status string `json:"status"`
}

type IO_CONTACT struct {
	SN      string `json:"sn"`
	Type    string `json:"type"`
	CONTACT string `json:"contact"`
}
