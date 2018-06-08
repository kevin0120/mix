package audi_vw

import (
	"fmt"
	"strconv"
	"encoding/xml"
	"strings"
	"github.com/masami10/rush/utils"
)

const (
	RESULT_NONE = "NONE"
	RESULT_OK   = "OK"
	RESULT_NOK  = "NOK"
)

const (
	RESULT_STAGE_INIT  = "init"
	RESULT_STAGE_FINAL = "final"
)

// header
const (
	HEADER_LEN = 32

	// HDR
	header_fixed = "55AA"

	// TYP
	Header_type_request_without_reply = 0
	Header_type_request_with_reply = 1
	Header_type_reply = 2
	Header_type_keep_alive = 3

	// COD
	Header_code_ok = 0
	Header_code_count_incorrect = 1
	Header_code_reserved = 2
	Header_code_length_incorrect = 3
	Header_code_xml_syntax_error = 4
	Header_code_xml_ver_conflict = 5
	Header_code_order_cannot_be_executed = 10
	Header_code_undefined_error = 99

	// REV
	header_rev = "0000"

	// RSD
	header_rsd = "0000"

	XML_RESULT_KEY = "<CUR>"
)

type PSetData struct {
	Name string `xml:"NAM"`
	Unit string `xml:"UNT"`
	Value float64 `xml:"VAL"`
}

type SMP struct {
	CUR_M string `xml:"Y1V"`
	CUR_W string `xml:"Y2V"`
}

type CUR struct {
	SMP SMP `xml:"SMP"`
}

type PRO struct {
	Strategy string `xml:"PAP"`
	Values []PSetData `xml:"MAR"`
}

type BLC struct {
	PRO PRO `xml:"PRO"`
	CUR CUR `xml:"CUR"`
}

type TIP struct {
	PSet int `xml:"PRG"`
	Date string `xml:"DAT"`
	Time string `xml:"TIM"`
	BLC BLC `xml:"BLC"`
}

type GRP struct {
	TIP TIP `xml:"TIP"`
}

type FAS struct {
	GRP GRP `xml:"GRP"`
}

type PAR struct {
	SN string `xml:"PRT"`
	Workorder_id int64 `xml:"PI1"`
	Result_id string `xml:"PI2"`
	Count int `xml:"STC"`
	Result string `xml:"PSC"`
	FAS FAS `xml:"FAS"`
}

type PRC_SST struct {
	PAR PAR `xml:"PAR"`
}

type CVI3Result struct {
	XMLName     xml.Name `xml:"ROOT"`
	PRC_SST		PRC_SST `xml:"PRC_SST"`
}

type CVI3Header struct {
	HDR string
	MID uint
	SIZ int
	TYP uint
	COD uint
	REV string
	RSD string
}

func (header *CVI3Header) Init() {
	header.HDR = header_fixed
	header.REV = header_rev
	header.RSD = header_rsd
}

func (header *CVI3Header) Check() (bool) {
	if header.COD == Header_code_ok {
		return true
	} else {
		return false
	}
}

func (header *CVI3Header) Serialize() string {
	return fmt.Sprintf("%s%04d%08d%04d%04d%s%s",
		header.HDR,
		header.MID,
		header.SIZ,
		header.TYP,
		header.COD,
		header.REV,
		header.RSD)
}

func (header *CVI3Header) Deserialize(header_str string) {
	header.Init()

	var n uint64
	var err error

	n, err = strconv.ParseUint(header_str[4:8], 10, 32)
	if err == nil {
		header.MID = uint(n)
	}

	n, err = strconv.ParseUint(header_str[8:16], 10, 32)
	if err == nil {
		header.SIZ = int(n)
	}

	n, err = strconv.ParseUint(header_str[16:20], 10, 32)
	if err == nil {
		header.TYP = uint(n)
	}

	n, err = strconv.ParseUint(header_str[20:24], 10, 32)
	if err == nil {
		header.COD = uint(n)
	}
}

func GeneratePacket(seq uint, typ uint, xmlpacket string) (string, uint) {
	header := CVI3Header{}
	header.Init()
	header.MID = seq
	header.SIZ = len(xmlpacket)
	header.TYP = typ
	header_str := header.Serialize()

	return fmt.Sprintf("%s%s", header_str, xmlpacket), header.MID
}

type ControllerCurve struct {
	ResultID  int64
	CurveFile string
	CurveData string
	Count     int
}

type ControllerCurveFile struct {
	Result string    `json:"result"`
	CUR_M  []float64 `json:"cur_m"`
	CUR_W  []float64 `json:"cur_w"`
}

type ControllerResult struct {
	Result_id     int64      `json:"result_id"`
	Controller_SN string     `json:"controller_sn"`
	Workorder_ID  int64      `json:"workorder_id"`
	CurFile       string     `json:"cur_file"`
	Result        string     `json:"result"`
	Dat           string     `json:"dat"`
	PSet          int        `json:"pset"`
	Count         int        `json:"count"`
	PSetDefine    PSetDefine `json:"pset_define"`

	ResultValue ResultValue `json:"result_value"`
}

type PSetDefine struct {
	Strategy string  `json:"strategy"`
	Mp       float64 `json:"M+"`
	Mm       float64 `json:"M-"`
	Ms       float64 `json:"MS"`
	Ma       float64 `json:"MA"`
	Wp       float64 `json:"W+"`
	Wm       float64 `json:"W-"`
	Wa       float64 `json:"WS"`
}

type ResultValue struct {
	Mi float64 `json:"MI"`
	Wi float64 `json:"WI"`
	Ti float64 `json:"TI"`
}

func XML2Curve(result CVI3Result) ControllerCurveFile {
	cur_result := ControllerCurveFile{}
	cur_result.Result = result.PRC_SST.PAR.Result
	if cur_result.Result == "IO" {
		cur_result.Result = RESULT_OK
	} else if cur_result.Result == "NIO" {
		cur_result.Result = RESULT_NOK
	}

	cur_ms := strings.Split(result.PRC_SST.PAR.FAS.GRP.TIP.BLC.CUR.SMP.CUR_M, " ")
	for i := range cur_ms {
		v, _ := strconv.ParseFloat(cur_ms[i], 64)
		cur_result.CUR_M = append(cur_result.CUR_M, v)
	}

	cur_ws := strings.Split(result.PRC_SST.PAR.FAS.GRP.TIP.BLC.CUR.SMP.CUR_W, " ")
	for i := range cur_ws {
		v, _ := strconv.ParseFloat(cur_ws[i], 64)
		cur_result.CUR_W = append(cur_result.CUR_W, v)
	}

	return cur_result
}

func XML2Result(result CVI3Result) ControllerResult {
	rr := ControllerResult{}

	rr.Controller_SN = result.PRC_SST.PAR.SN
	rr.Result = result.PRC_SST.PAR.Result
	if rr.Result == "IO" {
		rr.Result = RESULT_OK
	} else if rr.Result == "NIO" {
		rr.Result = RESULT_NOK
	}

	rr.PSet = result.PRC_SST.PAR.FAS.GRP.TIP.PSet
	rr.Workorder_ID = result.PRC_SST.PAR.Workorder_id
	rr.Dat = fmt.Sprintf("%s %s", result.PRC_SST.PAR.FAS.GRP.TIP.Date, result.PRC_SST.PAR.FAS.GRP.TIP.Time)
	result_id := result.PRC_SST.PAR.Result_id
	rid, _ := strconv.Atoi(result_id)
	rr.Result_id = int64(rid)
	rr.CurFile = fmt.Sprintf("%s_%d_%s_%s.json", rr.Controller_SN, rr.Workorder_ID, result_id, utils.GenerateID())
	rr.PSetDefine.Strategy = result.PRC_SST.PAR.FAS.GRP.TIP.BLC.PRO.Strategy
	rr.Count = result.PRC_SST.PAR.Count

	result_values := result.PRC_SST.PAR.FAS.GRP.TIP.BLC.PRO.Values
	for i := range result_values {
		switch result_values[i].Name {
		case "M+":
			rr.PSetDefine.Mp = result_values[i].Value
		case "M-":
			rr.PSetDefine.Mm = result_values[i].Value
		case "MS":
			rr.PSetDefine.Ms = result_values[i].Value
		case "MA":
			rr.PSetDefine.Ma = result_values[i].Value
		case "W+":
			rr.PSetDefine.Wp = result_values[i].Value
		case "W-":
			rr.PSetDefine.Wm = result_values[i].Value
		case "WA":
			rr.PSetDefine.Wa = result_values[i].Value
		case "MI":
			rr.ResultValue.Mi = result_values[i].Value
		case "WI":
			rr.ResultValue.Wi = result_values[i].Value
		case "tI":
			if result_values[i].Unit == "s" {
				rr.ResultValue.Ti = result_values[i].Value * 1000
			} else {
				rr.ResultValue.Ti = result_values[i].Value
			}
		}
	}

	return rr
}