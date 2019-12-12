package io

const (

	OUTPUT_STATUS_OFF   = 0
	OUTPUT_STATUS_ON    = 1
	OUTPUT_STATUS_FLASH = 2
)

type IoStatus struct {
	SN     string `json:"sn"`
	Status string `json:"status"`
}

type IoContact struct {
	Src     string `json:"src"`
	SN      string `json:"sn"`
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
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
