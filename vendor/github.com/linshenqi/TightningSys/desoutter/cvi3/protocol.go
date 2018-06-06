package cvi3

import (
	"fmt"
	"strconv"
	"encoding/xml"
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

func GeneratePacket(serial uint, typ uint, xmlpacket string) (string) {
	header := CVI3Header{}
	header.Init()
	header.MID = serial
	header.SIZ = len(xmlpacket)
	header.TYP = typ

	header_str := header.Serialize()

	return fmt.Sprintf("%s%s", header_str, xmlpacket)
}


