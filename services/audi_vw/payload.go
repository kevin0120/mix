package audi_vw

import (
	"encoding/xml"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"strconv"
	"strings"
)

// header
const (
	HEADER_LEN = 32

	// HDR
	header_fixed = "55AA"

	// TYP
	Header_type_request_without_reply = 0
	Header_type_request_with_reply    = 1
	Header_type_reply                 = 2
	Header_type_keep_alive            = 3

	// COD
	Header_code_ok                       = 0
	Header_code_count_incorrect          = 1
	Header_code_reserved                 = 2
	Header_code_length_incorrect         = 3
	Header_code_xml_syntax_error         = 4
	Header_code_xml_ver_conflict         = 5
	Header_code_order_cannot_be_executed = 10
	Header_code_undefined_error          = 99

	// REV
	header_rev = "0000"

	// RSD
	header_rsd = "0000"

	XML_CURVE_KEY  = "<CUR>"
	XML_RESULT_KEY = "<MAR>"
	XML_EVENT_KEY  = "<EVT>"
	XML_NUT_KEY    = "<NUT"
	XML_STATUS_KEY = "<RDY>"
)

var request_errors = map[uint]string{
	1:  "Telegram count incorrect",
	2:  "Reserved",
	3:  "Incorrect length",
	4:  "XML syntax error",
	5:  "XML protocol: Version number conflict",
	10: "Order cannot be executed",
	99: "Undefined error",
}

var nut_ids = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"D": 4,
}

type HandlerPkg struct {
	IP  string
	Msg string
}

type Evt struct {
	XMLName xml.Name `xml:"ROOT"`
	MSL_MSG struct {
		EVT struct {
			STS struct {
				ONC struct {
					RDY int `xml:"RDY"`
					NUT struct {
						NIDs []string `xml:"NID"`
					} `xml:"NUT"`
				} `xml:"ONC"`
			} `xml:"STS"`
		} `xml:"EVT"`
	} `xml:"MSL_MSG"`
}

type StatusMsg struct {
	XMLName xml.Name `xml:"ROOT"`
	MSL_MSG struct {
		GRP struct {
			IPA string `xml:"IPA"`
			MSG struct {
				MGS string `xml:"MGS"`
			} `xml:"MSG"`
		} `xml:"GRP"`
	} `xml:"MSL_MSG"`
}

type PSetData struct {
	Name  string  `xml:"NAM"`
	Unit  string  `xml:"UNT"`
	Value float64 `xml:"VAL"`
}

type SMP struct {
	CUR_M string `xml:"Y1V"`
	CUR_W string `xml:"Y2V"`
	CUR_T string `xml:"XVA"`
}

type CUR struct {
	STP float64 `xml:"STP"`
	STV float64 `xml:"STV"`
	CNT int     `xml:"CNT"`
	SMP SMP     `xml:"SMP"`
}

type PRO struct {
	Strategy string     `xml:"PAP"`
	Values   []PSetData `xml:"MAR"`
}

type BLC struct {
	PRO PRO `xml:"PRO"`
	CUR CUR `xml:"CUR"`
}

type TIP struct {
	ToolSN string `xml:"SID"`
	TID    string `xml:"TID"`
	PSet   int    `xml:"PRG"`
	Date   string `xml:"DAT"`
	Time   string `xml:"TIM"`
	STA    string `xml:"STA"`
	BLC    []BLC  `xml:"BLC"`
}

type GRP struct {
	TIP TIP `xml:"TIP"`
}

type FAS struct {
	UserID int64 `xml:"FAP"`
	GRP    GRP   `xml:"GRP"`
}

type PAR struct {
	SN           string `xml:"PRT"`
	Workorder_id int64  `xml:"PI1"`
	Result_id    string `xml:"PI2"`
	Count        int    `xml:"STC"`
	Result       string `xml:"PSC"`
	FAS          FAS    `xml:"FAS"`
}

type PRC_SST struct {
	PAR PAR `xml:"PAR"`
}

type CVI3Result struct {
	XMLName xml.Name `xml:"ROOT"`
	PRC_SST PRC_SST  `xml:"PRC_SST"`
}

type CVI3Header struct {
	HDR string
	MID uint32
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

func (header *CVI3Header) Check() bool {
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

func (header *CVI3Header) Deserialize(headerStr string) {
	header.Init()

	var n uint64
	var err error

	n, err = strconv.ParseUint(headerStr[4:8], 10, 32)
	if err == nil {
		header.MID = uint32(n)
	}

	n, err = strconv.ParseUint(headerStr[8:16], 10, 32)
	if err == nil {
		header.SIZ = int(n)
	}

	n, err = strconv.ParseUint(headerStr[16:20], 10, 32)
	if err == nil {
		header.TYP = uint(n)
	}

	n, err = strconv.ParseUint(headerStr[20:24], 10, 32)
	if err == nil {
		header.COD = uint(n)
	}
}

func GeneratePacket(seq uint32, typ uint, xmlpacket string) (string, uint32) {
	header := CVI3Header{}
	header.Init()
	header.MID = seq
	header.SIZ = len(xmlpacket)
	header.TYP = typ
	headerStr := header.Serialize()

	return fmt.Sprintf("%s%s", headerStr, xmlpacket), header.MID
}

func XML2Curve(result *CVI3Result, cur_result *minio.ControllerCurve) {
	cur_result.CurveContent = minio.ControllerCurveFile{}
	cur_result.CurveContent.Result = result.PRC_SST.PAR.Result
	if cur_result.CurveContent.Result == "IO" {
		cur_result.CurveContent.Result = storage.RESULT_OK
	} else if cur_result.CurveContent.Result == "NIO" {
		cur_result.CurveContent.Result = storage.RESULT_NOK
	}

	blc := result.PRC_SST.PAR.FAS.GRP.TIP.BLC
	cur_ms := strings.Split(blc[len(blc)-1].CUR.SMP.CUR_M, " ")
	cur_result.CurveContent.CUR_M = make([]float64, blc[len(blc)-1].CUR.CNT)
	for k, v := range cur_ms {
		m, _ := strconv.ParseFloat(v, 64)
		cur_result.CurveContent.CUR_M[k] = m
	}

	cur_ws := strings.Split(blc[len(blc)-1].CUR.SMP.CUR_W, " ")
	cur_result.CurveContent.CUR_W = make([]float64, blc[len(blc)-1].CUR.CNT)
	for k, v := range cur_ws {
		w, _ := strconv.ParseFloat(v, 64)
		cur_result.CurveContent.CUR_W[k] = w
	}

	stp := blc[len(blc)-1].CUR.STP
	stv := blc[len(blc)-1].CUR.STV
	if blc[len(blc)-1].CUR.SMP.CUR_T == "" {
		for i := 0; i < blc[len(blc)-1].CUR.CNT; i++ {
			x := float64(i)*stp + stv
			//t,_ := big.NewFloat(x).SetPrec(5).Float64()
			t, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", x), 64)
			cur_result.CurveContent.CUR_T = append(cur_result.CurveContent.CUR_T, t)
		}
	} else {
		cur_ts := strings.Split(blc[len(blc)-1].CUR.SMP.CUR_T, " ")
		cur_result.CurveContent.CUR_T = make([]float64, blc[len(blc)-1].CUR.CNT)
		for k, v := range cur_ts {
			w, _ := strconv.ParseFloat(v, 64)
			cur_result.CurveContent.CUR_T[k] = w
		}
	}

}

func XML2Result(result *CVI3Result, rr *controller.ControllerResult) {

	blcs := result.PRC_SST.PAR.FAS.GRP.TIP.BLC

	rr.TighteningID = result.PRC_SST.PAR.FAS.GRP.TIP.TID
	rr.Controller_SN = result.PRC_SST.PAR.SN
	rr.GunSN = result.PRC_SST.PAR.FAS.GRP.TIP.ToolSN
	rr.Result = result.PRC_SST.PAR.Result
	if rr.Result == "IO" {
		rr.Result = storage.RESULT_OK
	} else if rr.Result == "NIO" {
		rr.Result = storage.RESULT_NOK
	}

	rr.PSet = result.PRC_SST.PAR.FAS.GRP.TIP.PSet
	rr.Workorder_ID = result.PRC_SST.PAR.Workorder_id
	rr.UserID = result.PRC_SST.PAR.FAS.UserID
	rr.Dat = fmt.Sprintf("%s %s", result.PRC_SST.PAR.FAS.GRP.TIP.Date, result.PRC_SST.PAR.FAS.GRP.TIP.Time)

	seq, _ := strconv.Atoi(result.PRC_SST.PAR.Result_id)
	rr.Seq = seq
	//rr.CurFile = fmt.Sprintf("%s_%d_%s_%s.json", rr.Controller_SN, rr.Workorder_ID, result_id, utils.GenerateID())

	if result.PRC_SST.PAR.FAS.GRP.TIP.STA == "LSN" {
		rr.PSetDefine.Strategy = controller.STRATEGY_LN
		rr.Result = "LSN"
	} else {
		rr.PSetDefine.Strategy = blcs[len(blcs)-1].PRO.Strategy
	}

	rr.Count = result.PRC_SST.PAR.Count

	result_values := blcs[len(blcs)-1].PRO.Values
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
}
