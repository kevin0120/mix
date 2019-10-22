package io

const (
	WS_IO_STATUS  = "WS_IO_STATUS"
	WS_IO_CONTACT = "WS_IO_CONTACT"
	WS_IO_SET     = "WS_IO_SET"

	OUTPUT_STATUS_OFF   = 0
	OUTPUT_STATUS_ON    = 1
	OUTPUT_STATUS_FLASH = 2
)

type IoStatus struct {
	SN     string `json:"sn"`
	Status string `json:"status"`
}

type IoContact struct {
	SN      string `json:"sn"`
	Type    string `json:"type"`
	CONTACT string `json:"contact"`
}

type IoSet struct {
	SN     string `json:"sn"`
	Index  uint16 `json:"index"`
	Status uint16 `json:"status"`
}

type IoData struct {
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}

type IoConfig struct {
	InputNum  uint16 `json:"input_num"`
	OutputNum uint16 `json:"output_num"`
}
