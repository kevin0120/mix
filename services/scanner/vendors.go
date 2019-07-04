package scanner

import "fmt"

const (
	MODEL_HONEYWELL_1900    = "MODEL_HONEYWELL_1900"
	MODEL_DATALOGIC_GBT4400 = "MODEL_DATALOGIC_GBT4400"
)

var VENDOR_MODELS = map[string]Vendor{
	MODEL_HONEYWELL_1900:    &HONEYWELL{model: MODEL_HONEYWELL_1900},
	MODEL_DATALOGIC_GBT4400: &DATALOGIC{model: MODEL_DATALOGIC_GBT4400},
}

type Vendor interface {
	Parse([]byte) string

	// vendorID, productID
	ModelInfo() (uint16, uint16)
}

type HONEYWELL struct {
	model string
}

func (v *HONEYWELL) Parse(buf []byte) string {
	fmt.Printf("%d\n", len(buf))
	switch v.model {
	case MODEL_HONEYWELL_1900:
		return string(buf)
	}

	return string(buf)
}

func (v *HONEYWELL) ModelInfo() (uint16, uint16) {
	switch v.model {
	case MODEL_HONEYWELL_1900:
		return 3118, 2311
	}

	return 0, 0
}

type DATALOGIC struct {
	model string
}

func (v *DATALOGIC) Parse(buf []byte) string {
	return string(buf)
}

func (v *DATALOGIC) ModelInfo() (uint16, uint16) {
	switch v.model {
	case MODEL_DATALOGIC_GBT4400:
		return 1529, 8714
	}

	return 0, 0
}
