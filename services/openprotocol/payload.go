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
	MID_0151_IDENTIFIER_SUBSCRIBE  = "0151"
	MID_0150_IDENTIFIER_SET        = "0150"
	MID_0038_JOB_SELECT            = "0038"
	MID_0064_OLD_TIGHTING          = "0064"
	MID_0130_JOB_OFF               = "0130"

	MID_0061_LAST_RESULT = "0061"
	MID_7410_LAST_CURVE  = "7410"

	MID_9999_ALIVE = "9999"
)

const (
	MID_0038_REV_1 = "001"
	MID_0038_REV_2 = "002"
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

	case MID_0151_IDENTIFIER_SUBSCRIBE:
		h.MID = MID_0151_IDENTIFIER_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0150_IDENTIFIER_SET:
		h.MID = MID_0150_IDENTIFIER_SET
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0038_JOB_SELECT:
		h.MID = MID_0038_JOB_SELECT
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0064_OLD_TIGHTING:
		h.MID = MID_0064_OLD_TIGHTING
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0130_JOB_OFF:
		h.MID = MID_0130_JOB_OFF
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end
	}

	return ""
}

type ResultData struct {

	//rev2
	CellID                        int
	ChannelID                     int
	ControllerName                string
	VIN                           string
	JobID                         string
	PsetID                        string
	Strategy                      string
	StrategyOption                []byte
	BatchSize                     int
	BatchCount                    int
	TighteningStatus              string
	BatchStatus                   string
	TorqueStatus                  string
	AngleStatus                   string
	RundownAngleStatus            string
	CurrentMonitoringStatus       string
	SelftapStatus                 string
	PrevailTorqueMonitoringStatus string
	PrevailTorqueCompensateStatus string
	TighteningErrorStatus         []byte
	TorqueMin                     float64
	TorqueMax                     float64
	TorqueFinalTarget             float64
	Torque                        float64
	AngleMin                      float64
	AngleMax                      float64
	FinalAngleTarget              float64
	Angle                         float64
	RundownAngleMin               float64
	RundownAngleMax               float64
	RundownAngle                  float64
	CurrentMonitoringMin          float64
	CurrentMonitoringMax          float64
	CurrentMonitoring             float64
	SelftapMin                    float64
	SelftapMax                    float64
	SelftapTorque                 float64
	PrevailTorqueMonitoringMin    float64
	PrevailTorqueMonitoringMax    float64
	PrevailTorque                 float64
	TightingID                    int64
	JobSequenceNumber             int32
	SyncTighteningID              int32
	ToolSerialNumber              string
	TimeStamp                     string
	TimeStampPSetLastChange       string

	//rev3
	PSetName   string
	TorqueUnit string
	ResultType string

	//rev4
	ID2 string
	ID3 string
	ID4 string

	//rev5
	CustomerErrorCode string

	//rev6

	//rev998

}

func (rd *ResultData) Deserialize(str string) error {

	//rd.CellID = str[0:1]

	return nil
}
