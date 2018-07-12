package openprotocol

import (
	"fmt"
	"strconv"
)

const (
	LEN_HEADER      = 20
	DEFAULT_REV     = "000"
	DEFAULT_MSG_END = string(0x00)

	MID_0001_START                 = "0001"
	MID_0002_START_ACK             = "0002"
	MID_0003_STOP                  = "0003"
	MID_0004_STOP                  = "0004"
	MID_0010_PSET_IDS              = "0010"
	MID_0014_PSET_SUBSCRIBE        = "0014"
	MID_0018_PSET                  = "0018"
	MID_0034_JOB_INFO_SUBSCRIBE    = "0034"
	MID_0060_LAST_RESULT_SUBSCRIBE = "0060"
	MID_7408_LAST_CURVE_SUBSCRIBE  = "7408"

	MID_0061_LAST_RESULT = "0061"
	MID_7410_LAST_CURVE  = "7410"

	MID_9999_ALIVE = "9999"
)

type OpenProtocolHeader struct {
	LEN      int
	MID      string
	Revision string
	NoAck    string
	Station  string
	Spindle  string
	Spare    string
}

func (h *OpenProtocolHeader) Serialize() string {
	return fmt.Sprintf("%04d%04s%03s%-1s%-2s%-2s%-4s", h.LEN, h.MID, h.Revision, h.NoAck, h.Station, h.Spindle, h.Spare)
}

func (h *OpenProtocolHeader) Deserialize(str string) {
	n, _ := strconv.ParseInt(str[0:4], 10, 32)
	h.LEN = int(n) - LEN_HEADER
	h.MID = str[4:8]
	h.Revision = str[8:10]
	h.NoAck = str[10:11]
	h.Station = str[11:13]
	h.Spindle = str[13:15]
	h.Spare = str[15:19]
}

func GeneratePackage(mid string, rev string, data string, end string) string {
	h := OpenProtocolHeader{}
	switch mid {
	case MID_9999_ALIVE:
		h.MID = MID_9999_ALIVE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0001_START:
		h.MID = MID_0001_START
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0010_PSET_IDS:
		h.MID = MID_0010_PSET_IDS
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0018_PSET:
		h.MID = MID_0018_PSET
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0014_PSET_SUBSCRIBE:
		h.MID = MID_0014_PSET_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0034_JOB_INFO_SUBSCRIBE:
		h.MID = MID_0034_JOB_INFO_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0060_LAST_RESULT_SUBSCRIBE:
		h.MID = MID_0060_LAST_RESULT_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_7408_LAST_CURVE_SUBSCRIBE:
		h.MID = MID_7408_LAST_CURVE_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end
	}

	return ""
}
