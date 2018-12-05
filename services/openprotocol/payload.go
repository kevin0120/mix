package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/utils/biu"
	"strconv"
	"strings"
)

const (
	JOB_ACTION_ABORT = "abort"
)

const (
	IO_STATUS_ON       = "on"
	IO_STATUS_OFF      = "off"
	IO_STATUS_FLASHING = "flashing"
)

type IOStatus struct {
	No     int    `json:"no"`
	Status string `json:"status"`
}

const (
	LEN_HEADER         = 20
	DEFAULT_REV        = "000"
	DEFAULT_MSG_END    = string(0x00)
	LEN_SINGLE_SPINDLE = 18

	MID_0001_START                   = "0001"
	MID_0002_START_ACK               = "0002"
	MID_0003_STOP                    = "0003"
	MID_0004_CMD_ERR                 = "0004"
	MID_0005_CMD_OK                  = "0005"
	MID_0014_PSET_SUBSCRIBE          = "0014"
	MID_0018_PSET                    = "0018"
	MID_0012_PSET_DETAIL_REQUEST     = "0012"
	MID_0013_PSET_DETAIL_REPLY       = "0013"
	MID_0010_PSET_LIST_REQUEST       = "0010"
	MID_0011_PSET_LIST_REPLY         = "0011"
	MID_0034_JOB_INFO_SUBSCRIBE      = "0034"
	MID_0060_LAST_RESULT_SUBSCRIBE   = "0060"
	MID_7408_LAST_CURVE_SUBSCRIBE    = "7408"
	MID_0151_IDENTIFIER_SUBSCRIBE    = "0151"
	MID_0150_IDENTIFIER_SET          = "0150"
	MID_0038_JOB_SELECT              = "0038"
	MID_0064_OLD_SUBSCRIBE           = "0064"
	MID_0065_OLD_DATA                = "0065"
	MID_0130_JOB_OFF                 = "0130"
	MID_0250_SELECTOR_SUBSCRIBE      = "0250"
	MID_0042_TOOL_DISABLE            = "0042"
	MID_0043_TOOL_ENABLE             = "0043"
	MID_0030_JOB_LIST_REQUEST        = "0030"
	MID_0031_JOB_LIST_REPLY          = "0031"
	MID_0032_JOB_DETAIL_REQUEST      = "0032"
	MID_0033_JOB_DETAIL_REPLY        = "0033"
	MID_0200_CONTROLLER_RELAYS       = "0200"
	MID_0019_PSET_BATCH_SET          = "0019"
	MID_0020_PSET_BATCH_RESET        = "0020"
	MID_0035_JOB_INFO                = "0035"
	MID_0210_INPUT_SUBSCRIBE         = "0210"
	MID_0211_INPUT_MONITOR           = "0211"
	MID_0127_JOB_ABORT               = "0127"
	MID_0100_MULTI_SPINDLE_SUBSCRIBE = "0100"
	MID_0101_MULTI_SPINDLE_RESULT    = "0101"
	MID_0051_VIN_SUBSCRIBE           = "0051"
	MID_0052_VIN                     = "0052"

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

var request_errors = map[string]string{
	"00": "OK",
	"01": "Invalid data",
	"02": "Parameter set ID not present",
	"03": "Parameter set can not be set",
	"04": "Parameter set not running",
	"06": "VIN upload subscription already exists",
	"07": "VIN upload subscription does not exists",
	"08": "VIN input source not granted",
	"09": "Last tightening result subscription already exists",
	"10": "Last tightening result subscription does not exist",
	"11": "Alarm subscription already exists",
	"12": "Alarm subscription does not exist",
	"13": "Parameter set selection subscription already exists",
	"14": "Parameter set selection subscription does not exist",
	"15": "Tightening ID requested not found",
	"16": "Connection rejected protocol busy",
	"17": "Job ID not present",
	"18": "Job info subscription already exists",
	"19": "Job info subscription does not exist",
	"20": "Job can not be set",
	"21": "Job not running",
	"22": "Not possible to execute dynamic Job request",
	"23": "Job batch decrement failed",
	"24": "Not possible to create Pset",
	"25": "Programming control not granted",
	"30": "Controller is not a sync Master/station controller",
	"31": "Multi-spindle status subscription already exists",
	"32": "Multi-spindle status subscription does not exist",
	"33": "Multi-spindle result subscription already exists",
	"34": "Multi-spindle result subscription does not exist",
	"40": "Job line control info subscription already exists",
	"41": "Job line control info subscription does not exist",
	"42": "Identifier input source not granted",
	"43": "Multiple identifiers work order subscription already exists",
	"44": "Multiple identifiers work order subscription does not exist",
	"50": "Status external monitored inputs subscription already exists",
	"51": "Status external monitored inputs subscription does not exist",
	"52": "IO device not connected",
	"53": "Faulty IO device ID",
	"54": "Tool Tag ID unknown",
	"55": "Tool Tag ID subscription already exists",
	"56": "Tool Tag ID subscription does not exist",
	"58": "No alarm present",
	"59": "Tool currently in use",
	"60": "No histogram available",
	"70": "Calibration failed",
	"79": "Command failed",
	"80": "Audi emergency status subscription exists",
	"81": "Audi emergency status subscription does not exist",
	"82": "Automatic/Manual mode subscribe already exist",
	"83": "Automatic/Manual mode subscribe does not exist",
	"84": "The relay function subscription already exists",
	"85": "The relay function subscription does not exist",
	"86": "The selector socket info subscription already exist",
	"87": "The selector socket info subscription does not exist",
	"88": "The digin info subscription already exist",
	"89": "The digin info subscription does not exist",
	"90": "Lock at bach done subscription already exist",
	"91": "Lock at bach done subscription does not exist",
	"92": "Open protocol commands disabled",
	"93": "Open protocol commands disabled subscription already exists",
	"94": "Open protocol commands disabled subscription does not exist",
	"95": "Reject request, PowerMACS is in manual mode",
	"96": "Client already connected",
	"97": "MID revision unsupported",
	"98": "Controller internal request timeout",
	"99": "Unknown MID",
}

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
		h.LEN = LEN_HEADER + 100
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

	case MID_0032_JOB_DETAIL_REQUEST:
		h.MID = MID_0032_JOB_DETAIL_REQUEST
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0030_JOB_LIST_REQUEST:
		h.MID = MID_0030_JOB_LIST_REQUEST
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0042_TOOL_DISABLE:
		h.MID = MID_0042_TOOL_DISABLE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0043_TOOL_ENABLE:
		h.MID = MID_0043_TOOL_ENABLE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0200_CONTROLLER_RELAYS:
		h.MID = MID_0200_CONTROLLER_RELAYS
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0019_PSET_BATCH_SET:
		h.MID = MID_0019_PSET_BATCH_SET
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0020_PSET_BATCH_RESET:
		h.MID = MID_0020_PSET_BATCH_RESET
		h.LEN = LEN_HEADER + len(data)
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + data + end

	case MID_0210_INPUT_SUBSCRIBE:
		h.MID = MID_0210_INPUT_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0127_JOB_ABORT:
		h.MID = MID_0127_JOB_ABORT
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = ""
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0100_MULTI_SPINDLE_SUBSCRIBE:
		h.MID = MID_0100_MULTI_SPINDLE_SUBSCRIBE
		h.LEN = LEN_HEADER
		h.Revision = rev
		h.NoAck = "1"
		h.Station = ""
		h.Spindle = ""
		h.Spare = ""

		return h.Serialize() + end

	case MID_0051_VIN_SUBSCRIBE:
		h.MID = MID_0051_VIN_SUBSCRIBE
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

type IOMonitor struct {
	ControllerSN string `json:"controller_sn"`
	Inputs       string `json:"inputs"`
}

func (iom *IOMonitor) Deserialize(str string) error {

	iom.Inputs = str

	return nil
}

var result_errors = []string{
	"Rundown angle max shut off",
	"Rundown angle min shut off",
	"Torque max shut off",
	"Angle max shut off",
	"Selftap torque max shut off",
	"Selftap torque min shut off",
	"Prevail torque max shut off",
	"Prevail torque min shut off",
	"Prevail torque compensate overflow",
	"Current monitoring max shut off",
	"Post view torque min torque shut off",
	"Post view torque max torque shut off",
	"Post view torque Angle too small",
	"Trigger lost",
	"Torque less than target",
	"Tool hot",
	"Multistage abort",
	"Rehit",
	"DS measure failed",
	"Current limit reached",
	"EndTime out shutoff",
	"Remove fastener limit exceeded",
	"Disable drive",
	"Transducer lost",
	"Transducer shorted",
	"Transducer corrupt",
	"Sync timeout",
	"Dynamic current monitoring min",
	"Dynamic current monitoring max",
	"Angle max monitor",
	"Yield nut off",
	"Yield too few samples",
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
	TighteningErrorStatus         string
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
	rd.VIN = strings.TrimSpace(str[14:39])
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
	rd.ToolSerialNumber = strings.TrimSpace(str[171:285])
	rd.TorqueUnit = str[208:209]
	rd.ResultType = str[211:213]
	rd.ID2 = strings.TrimSpace(str[215:240])
	rd.ID3 = strings.TrimSpace(str[242:267])
	rd.ID4 = strings.TrimSpace(str[269:294])

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
	rd.VIN = strings.TrimSpace(str[39:64])

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

	error_status := str[127:137]
	error_value, err := strconv.ParseInt(error_status, 10, 32)
	if err != nil {
		return err
	}

	b_error := biu.ToBinaryString(error_value)
	b_error = strings.Trim(b_error, "[] ")
	b_error = strings.Replace(b_error, " ", "", -1)
	l := len(b_error)
	errs := []string{}
	for i := 0; i < l; i++ {
		v := b_error[l-1-i]

		if v == '1' {
			errs = append(errs, result_errors[i])
		}
	}

	rd.TighteningErrorStatus = strings.Join(errs, ",")

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

	rd.ToolSerialNumber = strings.TrimSpace(str[309:323])

	rd.TorqueUnit = str[394:395]
	rd.ResultType = str[397:399]
	rd.ID2 = strings.TrimSpace(str[401:426])
	rd.ID3 = strings.TrimSpace(str[428:453])
	rd.ID4 = strings.TrimSpace(str[455:480])

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

type JobList struct {
	num  int
	jobs []int
}

func (p *JobList) Deserialize(str string) error {
	var err error = nil
	p.num, err = strconv.Atoi(str[0:4])
	if err != nil {
		return err
	}

	for i := 0; i < p.num; i++ {
		job, _ := strconv.Atoi(str[(i+1)*4 : (i+1)*4+4])
		p.jobs = append(p.jobs, job)
	}

	return nil
}

type JobStep struct {
	StepName  string `json:"step_name"`
	ChannelID int    `json:"channel_id"`
	PSetID    int    `json:"pset_id"`
	BatchSize int    `json:"batch_size"`
	Socket    int    `json:"socket"`
}

type JobDetail struct {
	JobID         int    `json:"job"`
	JobName       string `json:"job_name"`
	OrderStrategy string `json:"order_strategy"`
	//MaxTimeforFirstTightening int
	//MaxTimetoCompleteJob int
	CountType         string    `json:"count_type"`
	LockAtJobDone     bool      `json:"lock_at_job_done"`
	UseLineControl    bool      `json:"use_line_control"`
	RepeatJob         bool      `json:"repeat_job"`
	LooseningStrategy string    `json:"loosening_strategy"`
	Steps             []JobStep `json:"steps"`
}

func (p *JobDetail) Deserialize(str string) error {
	var err error = nil

	p.JobID, err = strconv.Atoi(str[2:6])
	if err != nil {
		return err
	}

	p.JobName = strings.TrimSpace(str[8:33])

	order := str[35:36]
	switch order {
	case "0":
		p.OrderStrategy = "free"

	case "1":
		p.OrderStrategy = "forced"

	case "2":
		p.OrderStrategy = "free and forced"
	}

	//p.MaxTimeforFirstTightening, err = strconv.Atoi(str[38:42])
	//if err != nil {
	//	return err
	//}
	//
	//p.MaxTimetoCompleteJob, err = strconv.Atoi(str[44:49])
	//if err != nil {
	//	return err
	//}

	count_type := str[51:52]
	switch count_type {
	case "0":
		p.CountType = "only the OK tightenings are counted"

	case "1":
		p.CountType = "both the OK and NOK tightenings are counted"
	}

	if str[54:55] == "0" {
		p.LockAtJobDone = false
	} else {
		p.LockAtJobDone = true
	}

	if str[57:58] == "0" {
		p.UseLineControl = false
	} else {
		p.UseLineControl = true
	}

	if str[60:61] == "0" {
		p.RepeatJob = false
	} else {
		p.RepeatJob = true
	}

	loosening, err := strconv.Atoi(str[63:65])
	if err != nil {
		return err
	}

	switch loosening {
	case 0:
		p.LooseningStrategy = "enable"

	case 1:
		p.LooseningStrategy = "disable"

	case 2:
		p.LooseningStrategy = "enable only on NOK tightening"
	}

	step_str := str[75 : len(str)-1]
	steps := strings.Split(step_str, ";")
	job_step := JobStep{}
	for _, v := range steps {
		values := strings.Split(v, ":")

		job_step.ChannelID, _ = strconv.Atoi(values[0])
		job_step.PSetID, _ = strconv.Atoi(values[1])
		job_step.BatchSize, _ = strconv.Atoi(values[3])
		job_step.Socket, _ = strconv.Atoi(values[4])
		job_step.StepName = strings.TrimSpace(values[5])

		p.Steps = append(p.Steps, job_step)
	}

	return nil
}

type JobInfo struct {
	JobID           int
	JobStatus       int
	JobBatchMode    int
	JobBatchSize    int
	JobBatchCounter int
	Timestamp       string
	JobCurrentStep  int
	JobTotalStep    int
	JobStepType     int
}

func (ji *JobInfo) Deserialize(msg string) error {

	var err error

	ji.JobID, err = strconv.Atoi(msg[2:6])
	if err != nil {
		return err
	}

	ji.JobStatus, err = strconv.Atoi(msg[8:9])
	if err != nil {
		return err
	}

	ji.JobBatchMode, err = strconv.Atoi(msg[11:12])
	if err != nil {
		return err
	}

	ji.JobBatchSize, err = strconv.Atoi(msg[14:17])
	if err != nil {
		return err
	}

	ji.JobBatchCounter, err = strconv.Atoi(msg[20:24])
	if err != nil {
		return err
	}

	ji.Timestamp = msg[26:45]

	ji.JobCurrentStep, err = strconv.Atoi(msg[47:50])
	if err != nil {
		return err
	}

	ji.JobTotalStep, err = strconv.Atoi(msg[52:55])
	if err != nil {
		return err
	}

	ji.JobStepType, err = strconv.Atoi(msg[57:59])
	if err != nil {
		return err
	}

	return nil
}

func DeserializeIDS(str string) []string {
	rt := []string{}

	vin := strings.TrimSpace(str[2:27])
	rt = append(rt, vin)

	id2 := strings.TrimSpace(str[29:54])
	rt = append(rt, id2)

	id3 := strings.TrimSpace(str[56:71])
	rt = append(rt, id3)

	id4 := strings.TrimSpace(str[73:98])
	rt = append(rt, id4)

	return rt
}

type SingleSpindleResult struct {
	SpindleNo int
	ChannelID int
	Result    string
	Torque    float64
	Angle     float64
}

type MultiSpindleResult struct {
	TotalSpindleNumber int
	Vin                string
	JobID              int
	PSetID             int
	BatchSize          int
	BatchCount         int
	BatchStatus        int
	TorqueMin          float64
	TorqueMax          float64
	TorqueFinalTarget  float64
	AngleMin           float64
	AngleMax           float64
	FinalAngleTarget   float64

	Spindles []SingleSpindleResult
}

func (msr *MultiSpindleResult) Deserialize(str string) {

	sps := str[154:len(str)]

	sp_num := len(sps) / LEN_SINGLE_SPINDLE
	sp := SingleSpindleResult{}
	for i := 0; i < sp_num; i++ {
		target_sp := sps[i*LEN_SINGLE_SPINDLE : i*LEN_SINGLE_SPINDLE+LEN_SINGLE_SPINDLE]
		sp.SpindleNo, _ = strconv.Atoi(target_sp[0:2])
		if target_sp[4:5] == "0" {
			sp.Result = "NOK"
		} else {
			sp.Result = "OK"
		}

		sp.Torque, _ = strconv.ParseFloat(target_sp[6:12], 64)
		sp.Torque = sp.Torque / 100

		sp.Angle, _ = strconv.ParseFloat(target_sp[13:LEN_SINGLE_SPINDLE], 64)

		msr.Spindles = append(msr.Spindles, sp)
	}
}
