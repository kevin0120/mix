package cvi

import (
	"fmt"
	"strconv"
	"time"
	"strings"
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
)


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

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}
