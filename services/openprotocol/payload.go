package openprotocol

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	LEN_HEADER      = 20
	DEFAULT_REV     = "000"
	DEFAULT_MSG_END = string(0x00)

	MID_0001_START                 = "0001"
	MID_0002_START_ACK             = "0002"
	MID_0003_STOP                  = "0003"
	MID_0004_STOP                  = "0004"
	MID_0014_PSET_SUBSCRIBE        = "0014"
	MID_0018_PSET                  = "0018"
	MID_0012_PSET_DETAIL_REQUEST   = "0012"
	MID_0013_PSET_DETAIL_REPLY     = "0013"
	MID_0010_PSET_LIST_REQUEST     = "0010"
	MID_0011_PSET_LIST_REPLY       = "0011"
	MID_0034_JOB_INFO_SUBSCRIBE    = "0034"
	MID_0060_LAST_RESULT_SUBSCRIBE = "0060"
	MID_7408_LAST_CURVE_SUBSCRIBE  = "7408"
	MID_0151_IDENTIFIER_SUBSCRIBE  = "0151"
	MID_0150_IDENTIFIER_SET        = "0150"
	MID_0038_JOB_SELECT            = "0038"
	MID_0064_OLD_SUBSCRIBE         = "0064"
	MID_0065_OLD_DATA              = "0065"
	MID_0130_JOB_OFF               = "0130"
	MID_0250_SELECTOR_SUBSCRIBE    = "0250"

	MID_0008_DATA_SUB = "0008"

	MID_0061_LAST_RESULT = "0061"
	MID_7410_LAST_CURVE  = "7410"

	MID_9999_ALIVE = "9999"
)

const (
	ROTATION_CW  = "CW"
	ROTATION_CCW = "CCW"
)

const (
	MODE_PSET = "pset"
	MODE_JOB  = "job"
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

	case MID_0250_SELECTOR_SUBSCRIBE:
		h.MID = MID_0250_SELECTOR_SUBSCRIBE
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

	case MID_0064_OLD_SUBSCRIBE:
		h.MID = MID_0064_OLD_SUBSCRIBE
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

	case MID_0008_DATA_SUB:
		h.MID = MID_0008_DATA_SUB
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0012_PSET_DETAIL_REQUEST:
		h.MID = MID_0012_PSET_DETAIL_REQUEST
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0010_PSET_LIST_REQUEST:
		h.MID = MID_0010_PSET_LIST_REQUEST
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end
	}

	return ""
}

type ResultData struct {

	//rev2
	CellID                        int
	ChannelID                     int
	ControllerName                string
	VIN                           string
	JobID                         int
	PSetID                        int
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
	TightingID                    string
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
	NumberOfStages       int
	NumberOfStageResults int
	StageResult          string
}

func (rd *ResultData) DeserializeOld(str string) error {
	var err error = nil

	rd.TightingID = str[2:12]
	rd.VIN = str[14:39]
	rd.JobID, err = strconv.Atoi(str[41:45])
	if err != nil {
		return err
	}

	rd.PSetID, err = strconv.Atoi(str[47:50])
	if err != nil {
		return err
	}

	rd.Strategy = str[52:54]
	rd.BatchSize, err = strconv.Atoi(str[63:67])
	if err != nil {
		return err
	}

	rd.BatchCount, err = strconv.Atoi(str[69:73])
	if err != nil {
		return err
	}

	rd.TighteningStatus = str[75:76]
	rd.BatchStatus = str[78:79]
	rd.TorqueStatus = str[81:82]
	rd.AngleStatus = str[84:85]
	rd.RundownAngleStatus = str[87:88]
	rd.CurrentMonitoringStatus = str[90:91]
	rd.SelftapStatus = str[93:94]
	rd.PrevailTorqueMonitoringStatus = str[96:97]
	rd.PrevailTorqueCompensateStatus = str[99:100]

	rd.Torque, err = strconv.ParseFloat(str[114:120], 64)
	if err != nil {
		return err
	}

	rd.Angle, err = strconv.ParseFloat(str[122:127], 64)
	if err != nil {
		return err
	}

	rd.RundownAngle, err = strconv.ParseFloat(str[129:134], 64)
	if err != nil {
		return err
	}

	rd.CurrentMonitoring, err = strconv.ParseFloat(str[136:139], 64)
	if err != nil {
		return err
	}

	rd.SelftapTorque, err = strconv.ParseFloat(str[141:147], 64)
	if err != nil {
		return err
	}

	rd.PrevailTorque, err = strconv.ParseFloat(str[149:155], 64)
	if err != nil {
		return err
	}

	rd.TimeStamp = str[187:206]
	rd.TorqueUnit = str[208:209]
	rd.ResultType = str[211:213]
	rd.ID2 = str[215:240]
	rd.ID3 = str[242:267]
	rd.ID4 = str[269:294]

	return nil
}

func (rd *ResultData) Deserialize(str string) error {

	var err error = nil
	rd.CellID, err = strconv.Atoi(str[2:6])
	if err != nil {
		return err
	}

	rd.ChannelID, err = strconv.Atoi(str[8:10])
	if err != nil {
		return err
	}

	rd.ControllerName = str[12:37]
	rd.VIN = str[39:64]

	rd.JobID, err = strconv.Atoi(str[66:70])
	if err != nil {
		return err
	}

	rd.PSetID, err = strconv.Atoi(str[72:75])
	if err != nil {
		return err
	}

	rd.Strategy = str[77:79]

	rd.BatchSize, err = strconv.Atoi(str[88:92])
	if err != nil {
		return err
	}

	rd.BatchCount, err = strconv.Atoi(str[94:98])
	if err != nil {
		return err
	}

	rd.TighteningStatus = str[100:101]
	rd.BatchStatus = str[103:104]
	rd.TorqueStatus = str[106:107]
	rd.AngleStatus = str[109:110]
	rd.RundownAngleStatus = str[112:113]
	rd.CurrentMonitoringStatus = str[115:116]
	rd.SelftapStatus = str[118:119]
	rd.PrevailTorqueMonitoringStatus = str[121:122]
	rd.PrevailTorqueCompensateStatus = str[124:125]

	//error_status := str[127:137]
	//for _, v := range error_status {
	//	rd.TighteningErrorStatus = append(rd.TighteningErrorStatus, v)
	//}

	rd.TorqueMin, err = strconv.ParseFloat(str[139:145], 64)
	if err != nil {
		return err
	}

	rd.TorqueMax, err = strconv.ParseFloat(str[147:153], 64)
	if err != nil {
		return err
	}

	rd.TorqueFinalTarget, err = strconv.ParseFloat(str[155:161], 64)
	if err != nil {
		return err
	}

	rd.Torque, err = strconv.ParseFloat(str[163:169], 64)
	if err != nil {
		return err
	}

	rd.AngleMin, err = strconv.ParseFloat(str[171:176], 64)
	if err != nil {
		return err
	}

	rd.AngleMax, err = strconv.ParseFloat(str[178:183], 64)
	if err != nil {
		return err
	}

	rd.FinalAngleTarget, err = strconv.ParseFloat(str[185:190], 64)
	if err != nil {
		return err
	}

	rd.Angle, err = strconv.ParseFloat(str[192:197], 64)
	if err != nil {
		return err
	}

	rd.TightingID = str[283:293]

	rd.TimeStamp = str[325:344]

	rd.TorqueUnit = str[394:395]
	rd.ResultType = str[397:399]
	rd.ID2 = str[401:426]
	rd.ID3 = str[428:453]
	rd.ID4 = str[455:480]

	rd.NumberOfStages, err = strconv.Atoi(str[508:510])
	if err != nil {
		return err
	}

	rd.NumberOfStageResults, err = strconv.Atoi(str[512:514])
	if err != nil {
		return err
	}

	rd.StageResult = str[516:527]

	return nil
}

type PSetDetail struct {
	PSetID            int     `json:"pset"`
	PSetName          string  `json:"pset_name"`
	RotationDirection string  `json:"rotation_direction"`
	BatchSize         int     `json:"batch_size"`
	TorqueMin         float64 `json:"torque_min"`
	TorqueMax         float64 `json:"torque_max"`
	TorqueTarget      float64 `json:"torque_target"`
	AngleMin          float64 `json:"angle_min"`
	AngleMax          float64 `json:"angle_max"`
	AngleTarget       float64 `json:"angle_target"`
}

func (p *PSetDetail) Deserialize(str string) error {
	var err error = nil

	p.PSetID, err = strconv.Atoi(str[2:5])
	if err != nil {
		return err
	}

	p.PSetName = strings.TrimSpace(str[7:32])
	p.RotationDirection = str[34:35]

	switch p.RotationDirection {
	case "1":
		p.RotationDirection = ROTATION_CW

	case "2":
		p.RotationDirection = ROTATION_CCW
	}

	p.BatchSize, err = strconv.Atoi(str[37:39])
	if err != nil {
		return err
	}

	p.TorqueMin, err = strconv.ParseFloat(str[41:47], 64)
	if err != nil {
		return err
	}

	p.TorqueMin = p.TorqueMin / 100

	p.TorqueMax, err = strconv.ParseFloat(str[49:55], 64)
	if err != nil {
		return err
	}

	p.TorqueMax = p.TorqueMax / 100

	p.TorqueTarget, err = strconv.ParseFloat(str[57:63], 64)
	if err != nil {
		return err
	}

	p.TorqueTarget = p.TorqueTarget / 100

	p.AngleMin, err = strconv.ParseFloat(str[65:70], 64)
	if err != nil {
		return err
	}

	p.AngleMax, err = strconv.ParseFloat(str[72:77], 64)
	if err != nil {
		return err
	}

	p.AngleTarget, err = strconv.ParseFloat(str[79:84], 64)
	if err != nil {
		return err
	}

	return nil
}

type PSetList struct {
	num   int
	psets []int
}

func (p *PSetList) Deserialize(str string) error {
	var err error = nil
	p.num, err = strconv.Atoi(str[0:3])
	if err != nil {
		return err
	}

	for i := 0; i < p.num; i++ {
		pset, _ := strconv.Atoi(str[(i+1)*3 : (i+1)*3+3])
		p.psets = append(p.psets, pset)
	}

	return nil
}
