package openprotocol

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/masami10/rush/services/tightening_device"
	"strconv"
	"strings"
	"time"
)

const (
	OP_TERMINAL = 0x00
)

const (
	JOB_INFO_NOT_COMPLETED = 0
	JOB_ACTION_ABORT       = "abort"
	MAX_IDS_NUM            = 4
)

const (
	IO_STATUS_ON       = "on"
	IO_STATUS_OFF      = "off"
	IO_STATUS_FLASHING = "flashing"
)

const (
	EVT_CONTROLLER_NO_ERR          = "E000"
	EVT_CONTROLLER_TOOL_DISCONNECT = "I003"
	EVT_CONTROLLER_TOOL_CONNECT    = "I002"
)

type IOStatus struct {
	No     int    `json:"no"`
	Status string `json:"status"`
}

const (
	LEN_HEADER         = 20
	DEFAULT_REV        = "000"
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
	MID_0040_TOOL_INFO_REQUEST       = "0040"
	MID_0041_TOOL_INFO_REPLY         = "0041"
	MID_0060_LAST_RESULT_SUBSCRIBE   = "0060"
	MID_7408_LAST_CURVE_SUBSCRIBE    = "7408"
	MID_0151_IDENTIFIER_SUBSCRIBE    = "0151"
	MID_0150_IDENTIFIER_SET          = "0150"
	MID_0038_JOB_SELECT              = "0038"
	MID_0064_OLD_SUBSCRIBE           = "0064"
	MID_0065_OLD_DATA                = "0065"
	MID_0070_ALARM_SUBSCRIBE         = "0070"
	MID_0071_ALARM                   = "0071"
	MID_0076_ALARM_STATUS            = "0076"
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

const ()

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
	"30": "TighteningController is not a sync Master/station controller",
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
	"52": "IO tightening_device not connected",
	"53": "Faulty IO tightening_device ID",
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
	"98": "TighteningController internal request timeout",
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
	if len(str) != LEN_HEADER {
		return
	}

	n, _ := strconv.ParseInt(str[0:4], 10, 32)
	h.LEN = int(n) - LEN_HEADER
	h.MID = str[4:8]
	h.Revision = str[8:10]
	h.NoAck = str[10:11]
	h.Station = str[11:13]
	h.Spindle = str[13:15]
	h.Spare = str[15:19]
}

func GeneratePackage(mid string, rev string, noack string, station string, spindle string, data string) string {
	h := OpenProtocolHeader{
		MID:      mid,
		LEN:      LEN_HEADER + len(data),
		Revision: rev,
		NoAck:    noack,
		Station:  station,
		Spindle:  spindle,
		Spare:    "",
	}

	return h.Serialize() + data + string(OP_TERMINAL)
}

type IOMonitor struct {
	ControllerSN string `json:"controller_sn"`
	Inputs       string `json:"inputs"`
}

func (iom *IOMonitor) Deserialize(str string) error {

	iom.Inputs = str

	return nil
}

func (iom *IOMonitor) ToTighteningControllerInput() *tightening_device.TighteningControllerInput {
	return &tightening_device.TighteningControllerInput{
		ControllerSN: iom.ControllerSN,
		Inputs:       iom.Inputs,
	}
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
	CellID         int    `start:"3"  end:"6"`
	ChannelID      int    `start:"9"  end:"10"`
	ControllerName string `start:"13"  end:"37"`
	VIN            string `start:"40"  end:"64"`
	JobID          int    `start:"67"  end:"70"`
	PSetID         int    `start:"73"  end:"75"`
	Strategy       string `start:"78"  end:"79"`
	//StrategyOption                []byte	`start:"3"  end:"6"`
	BatchSize                     int     `start:"89"  end:"92"`
	BatchCount                    int     `start:"95"  end:"98"`
	TighteningStatus              string  `start:"101"  end:"101"`
	BatchStatus                   string  `start:"104"  end:"104"`
	TorqueStatus                  string  `start:"107"  end:"107"`
	AngleStatus                   string  `start:"110"  end:"110"`
	RundownAngleStatus            string  `start:"113"  end:"113"`
	CurrentMonitoringStatus       string  `start:"116"  end:"116"`
	SelftapStatus                 string  `start:"119"  end:"119"`
	PrevailTorqueMonitoringStatus string  `start:"122"  end:"122"`
	PrevailTorqueCompensateStatus string  `start:"125"  end:"125"`
	TighteningErrorStatus         string  `start:"128"  end:"137"`
	TorqueMin                     float64 `start:"140"  end:"145"`
	TorqueMax                     float64 `start:"148"  end:"153"`
	TorqueFinalTarget             float64 `start:"156"  end:"161"`
	Torque                        float64 `start:"164"  end:"169"`
	AngleMin                      float64 `start:"172"  end:"176"`
	AngleMax                      float64 `start:"179"  end:"183"`
	FinalAngleTarget              float64 `start:"186"  end:"190"`
	Angle                         float64 `start:"193"  end:"197"`
	//RundownAngleMin               float64	`start:"3"  end:"6"`
	//RundownAngleMax               float64	`start:"3"  end:"6"`
	//RundownAngle                  float64	`start:"130"  end:"134"`
	//CurrentMonitoringMin          float64	`start:"3"  end:"6"`
	//CurrentMonitoringMax          float64	`start:"3"  end:"6"`
	//CurrentMonitoring             float64	`start:"137"  end:"139"`
	//SelftapMin                    float64	`start:"3"  end:"6"`
	//SelftapMax                    float64	`start:"3"  end:"6"`
	//SelftapTorque                 float64	`start:"142"  end:"147"`
	//PrevailTorqueMonitoringMin    float64	`start:"3"  end:"6"`
	//PrevailTorqueMonitoringMax    float64	`start:"3"  end:"6"`
	//PrevailTorque                 float64	`start:"150"  end:"155"`
	TightingID string `start:"284"  end:"293"`
	//JobSequenceNumber             int32		`start:"3"  end:"6"`
	//SyncTighteningID              int32		`start:"3"  end:"6"`
	ToolSerialNumber string `start:"310"  end:"323"`
	TimeStamp        string `start:"326"  end:"344"`
	//TimeStampPSetLastChange       string	`start:"3"  end:"6"`

	//rev3
	//PSetName   string	`start:"3"  end:"6"`
	TorqueUnit string `start:"395"  end:"395"`
	ResultType string `start:"398"  end:"399"`

	//rev4
	ID2 string `start:"402"  end:"426"`
	ID3 string `start:"429"  end:"453"`
	ID4 string `start:"456"  end:"480"`

	//rev5
	//CustomerErrorCode string	`start:"3"  end:"6"`

	//rev6

	//rev998
	NumberOfStages int `start:"3"  end:"6"`
	//NumberOfStageResults int	`start:"3"  end:"6"`
	StageResult string `start:"3"  end:"..."`
}

type ResultDataOld struct {

	//rev2
	VIN                           string  `start:"15"  end:"39"`
	JobID                         int     `start:"42"  end:"45"`
	PSetID                        int     `start:"48"  end:"50"`
	Strategy                      string  `start:"53"  end:"54"`
	BatchSize                     int     `start:"64"  end:"67"`
	BatchCount                    int     `start:"70"  end:"73"`
	TighteningStatus              string  `start:"76"  end:"76"`
	BatchStatus                   string  `start:"79"  end:"79"`
	TorqueStatus                  string  `start:"82"  end:"82"`
	AngleStatus                   string  `start:"85"  end:"85"`
	RundownAngleStatus            string  `start:"88"  end:"88"`
	CurrentMonitoringStatus       string  `start:"91"  end:"91"`
	SelftapStatus                 string  `start:"94"  end:"94"`
	PrevailTorqueMonitoringStatus string  `start:"97"  end:"97"`
	PrevailTorqueCompensateStatus string  `start:"100"  end:"100"`
	Torque                        float64 `start:"115"  end:"120"`
	Angle                         float64 `start:"123"  end:"127"`
	RundownAngle                  float64 `start:"130"  end:"134"`
	CurrentMonitoring             float64 `start:"137"  end:"139"`
	SelftapTorque                 float64 `start:"142"  end:"147"`
	PrevailTorque                 float64 `start:"150"  end:"155"`
	TightingID                    string  `start:"3"  end:"12"`
	ToolSerialNumber              string  `start:"172"  end:"285"`
	TimeStamp                     string  `start:"188"  end:"206"`
	//rev3
	TorqueUnit string `start:"209"  end:"209"`
	ResultType string `start:"212"  end:"213"`

	//rev4
	ID2 string `start:"216"  end:"240"`
	ID3 string `start:"243"  end:"267"`
	ID4 string `start:"270"  end:"294"`

	//rev5

	//rev6

	//rev998
	StageResult string `start:"3"  end:"..."`
}

// TODO
func (rd *ResultData) ToTighteningResult() *tightening_device.TighteningResult {
	measureResult := tightening_device.RESULT_OK
	if rd.TighteningStatus == "0" {
		measureResult = tightening_device.RESULT_NOK
	}

	strategy := ""

	switch rd.Strategy {
	case "01":
		strategy = tightening_device.STRATEGY_AW

	case "02":
		strategy = tightening_device.STRATEGY_AW

	case "03":
		strategy = tightening_device.STRATEGY_ADW

	case "04":
		strategy = tightening_device.STRATEGY_AD
	}

	if rd.ResultType == "02" {
		measureResult = tightening_device.RESULT_LSN
		strategy = tightening_device.STRATEGY_LN
	}

	return &tightening_device.TighteningResult{
		// 工具序列号
		ToolSN: strings.TrimSpace(rd.ToolSerialNumber),

		// 收到时间
		UpdateTime: time.Now(),

		// job号
		Job: rd.JobID,

		// pset号
		PSet: rd.PSetID,

		// 拧紧ID
		TighteningID: rd.TightingID,

		// 实际结果
		MeasureResult: measureResult,

		// 实际扭矩
		MeasureTorque: rd.Torque / 100,

		// 实际角度
		MeasureAngle: rd.Angle,

		// 实际耗时
		MeasureTime: 0,

		// 拧紧策略
		Strategy: strategy,

		// 最大扭矩
		TorqueMax: rd.TorqueMax / 100,

		// 最小扭矩
		TorqueMin: rd.TorqueMin / 100,

		// 扭矩阈值
		TorqueThreshold: rd.TorqueMax / 100,

		// 目标扭矩
		TorqueTarget: rd.TorqueFinalTarget / 100,

		// 最大角度
		AngleMax: rd.AngleMax,

		// 最小角度
		AngleMin: rd.AngleMin,

		// 目标角度
		AngleTarget: rd.FinalAngleTarget,
	}
}

func DeserializePSetDetail(str string) (*tightening_device.PSetDetail, error) {
	var p tightening_device.PSetDetail
	var err error = nil

	p.PSetID, err = strconv.Atoi(str[2:5])
	if err != nil {
		return nil, err
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
		return nil, err
	}

	p.TorqueMin, err = strconv.ParseFloat(str[41:47], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueMin = p.TorqueMin / 100

	p.TorqueMax, err = strconv.ParseFloat(str[49:55], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueMax = p.TorqueMax / 100

	p.TorqueTarget, err = strconv.ParseFloat(str[57:63], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueTarget = p.TorqueTarget / 100

	p.AngleMin, err = strconv.ParseFloat(str[65:70], 64)
	if err != nil {
		return nil, err
	}

	p.AngleMax, err = strconv.ParseFloat(str[72:77], 64)
	if err != nil {
		return nil, err
	}

	p.AngleTarget, err = strconv.ParseFloat(str[79:84], 64)
	if err != nil {
		return nil, err
	}

	return &p, nil
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

func DeserializeJobDetail(str string) (*tightening_device.JobDetail, error) {
	var err error = nil
	var p tightening_device.JobDetail

	p.JobID, err = strconv.Atoi(str[2:6])
	if err != nil {
		return nil, err
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
		return nil, err
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
	job_step := tightening_device.JobStep{}
	for _, v := range steps {
		values := strings.Split(v, ":")

		job_step.ChannelID, _ = strconv.Atoi(values[0])
		job_step.PSetID, _ = strconv.Atoi(values[1])
		job_step.BatchSize, _ = strconv.Atoi(values[3])
		job_step.Socket, _ = strconv.Atoi(values[4])
		job_step.StepName = strings.TrimSpace(values[5])

		p.Steps = append(p.Steps, job_step)
	}

	return &p, nil
}

type AlarmInfo struct {
	ErrorCode      string `start:"3"  end:"6"`
	isControllerOK bool   `start:"9"  end:"9"`
	isToolOK       bool   `start:"12"  end:"12"`
}

type ToolInfo struct {
	ToolSN               string `json:"serial_no"`
	ControllerSN         string `json:"controller_sn"`
	TotalTighteningCount int    `json:"times"`
	CountSinLastService  int    `json:"sin_last_service"`
}

func (ti *ToolInfo) Deserialize(msg string) error {

	var err error = nil

	if len(msg) < 20 {
		return errors.New("tool info msg len is not enough")
	}

	ti.ToolSN = strings.TrimSpace(msg[2:16]) //去除空格获取序列号

	ti.TotalTighteningCount, err = strconv.Atoi(msg[18:28])

	ti.ControllerSN = strings.TrimSpace(msg[51:61])

	ti.CountSinLastService, err = strconv.Atoi(msg[92:102])

	return err
}

type JobInfo struct {
	JobID               int    `start:"3"  end:"6"`
	JobStatus           int    `start:"9"  end:"9"`
	JobBatchMode        int    `start:"12"  end:"12"`
	JobBatchSize        int    `start:"15"  end:"17"`
	JobBatchCounter     int    `start:"21"  end:"24"`
	Timestamp           string `start:"27"  end:"45"`
	JobCurrentStep      int    `start:"48"  end:"50"`
	JobTotalStep        int    `start:"53"  end:"55"`
	JobStepType         int    `start:"58"  end:"59"`
	JobTighteningStatus int    `start:"62"  end:"63"`
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

type AlarmStatus struct {
	Status         string `start:"3"  end:"3"`
	ErrorCode      string `start:"6"  end:"9"`
	IsControllerOK bool   `start:"12"  end:"12"`
	IsToolOK       bool   `start:"15"  end:"15"`
}

type CurveBody struct {
	ToolNumber    int    `start:"3"  end:"4"`
	TorqueString  string `start:"28"  end:"41"`
	AngleString   string `start:"44"  end:"57"`
	MeasurePoints int    `start:"60"  end:"63"`
	Num           string `start:"66"  end:"67"` //曲线总共分几段
	Id            string `start:"70"  end:"71"` //当前为第几段
	Data          string `start:"72"  end:"..."`
}
