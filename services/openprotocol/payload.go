package openprotocol

import "fmt"

const (
	LEN_HEADER  = 20
	DEFAULT_REV = "000"

	MID_0001_START     = "0001"
	MID_0002_START_ACK = "0002"
	MID_0003_STOP      = "0003"
	MID_0004_STOP      = "0004"
	MID_9999_ALIVE     = "9999"
)

type Header struct {
	LEN      int
	MID      string
	Revision string
	NoAck    string
	Station  string
	Spindle  string
	Spare    string
}

func (h *Header) Serialize() string {
	return fmt.Sprintf("%04d%04s%03s%-1s%-2s%-2s%-4s", h.LEN, h.MID, h.Revision, h.NoAck, h.Station, h.Spindle, h.Spare)
}

func GeneratePackage(mid string, rev string, data string, end string) string {
	h := Header{}
	switch mid {
	case MID_0001_START:
		h.MID = MID_0001_START
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize()
	}

	return ""
}
