package pmon

import (
	"fmt"
	"math"
	"strconv"
)

type PMONSMGTYPE = string

const (
	PMONMSGSO = "SO" //Submission-Open
	PMONMSGAO = "AO" //Acknowledge Submission Open
	PMONMSGSC = "SC" //Submission Close
	PMONMSGAC = "AC" //Acknowledge Submission Close
	PMONMSGSD = "SD" //Send data
	PMONMSGAD = "AD" //Acknowledge Data
	PMONMSGUN = "UNDEFINED"
)

const (
	PMONSTARTREQ = string(0x02) //STX
	PMONSTARTRES = string(0x06) //ACK
)

const (
	PMONSTOP1 = string(0x10) //DLE
	PMONSTOP2 = string(0x03) //ETX
)

const PMONEND = PMONSTOP1 + PMONSTOP2

type PmonPackage struct {
	t    PMONSMGTYPE
	data []byte
}

func ValidateChecksum(buf []byte) error {
	s := checksumByte(buf[:len(buf)-1])
	if buf[len(buf)-1] != s {
		return fmt.Errorf("Validate checksum Fail %0x ", buf[len(buf)-1])
	}
	return nil
}

func IsUDPResponse(buf []byte) bool {
	if len(buf) == 12 {
		return true
	}
	return false
}

func checksumByte(b []byte) byte {
	ret := make([]byte, 1)
	ret[0] = 0x00
	for _, d := range b {
		x := ret[0] ^ d
		ret[0] = x
	}
	return ret[0]
}

func checksum(s string) string {
	data := []byte(s)
	ret := make([]byte, 1)
	ret[0] = 0x00
	for _, d := range data {
		x := ret[0] ^ d
		ret[0] = x
	}
	return string(ret[0])
}

func GetMsgId(buf []byte) int {
	x, _ := strconv.Atoi(string(buf[1:5]))
	return x
}

func UdpResponse(msgNum int) string {
	s := PMONSTARTRES
	x := fmt.Sprintf("%04d", msgNum)
	s += x + "0000" + PMONEND
	s += checksum(s)
	return s
}

func (c *Channel) setPMONHeader(start string, msgNum int, t PMONSMGTYPE) string {
	s := start
	x := fmt.Sprintf("%04d", msgNum)
	s += x + c.SNoT + c.SNoR + t
	return s
}

func (c *Channel) generateSO(msgNum int) (string, error) {
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGSO)
	s += "00000000000000000000" //20个0 for Length of Object ID + Number of Records + Record Length + Generation Number
	s += c.Segment
	s += "00"                                      //data Security
	s += fmt.Sprintf("%04d", c.Buffer)             //四位补零
	s += fmt.Sprintf("%04d", c.RestartPointLength) //四位补零,Length of the starting point
	if c.RestartPointLength > 0 {
		s += c.GetRestartPoint()
	}
	s += "0000"
	s += PMONEND
	s += checksum(s)
	return s, nil
}

func (c *Channel) generateSC(msgNum int) (string, error) {
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGSC)
	s += "00" // close reason
	s += PMONEND
	s += checksum(s)
	return s, nil
}

func (c *Channel) generateAO(msgNum int) (string, error) {
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGAO)
	s += "0000000000" //8个0 for Length of Object ID +Generation Number
	s += c.Segment
	s += "00"                                      //data Security
	s += fmt.Sprintf("%04d", c.Buffer)             //四位补零
	s += fmt.Sprintf("%04d", c.RestartPointLength) //四位补零,Length of the starting point
	//for idx :=0; idx < c.RestartPointLength; idx ++ {
	//	s += "0"
	//}
	if c.RestartPointLength > 0 {
		s += c.GetRestartPoint()
	}
	s += "0000"
	s += PMONEND
	s += checksum(s)
	return s, nil
}

func (c *Channel) generateAC(msgNum int) (string, error) {
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGAC)
	s += PMONEND
	s += checksum(s)
	return s, nil
}

func (c *Channel) generateAD(msgNum int, rBlockCount string) (string, error) {
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGAD)
	s += "0000" + rBlockCount // Acknowledge Info + Recipient Info + rBlockCount
	s += PMONEND
	s += checksum(s)
	return s, nil
}

func (c *Channel) generateSD(msgNum int, data string) ([]string, error) {
	d := []byte(data)
	s := c.setPMONHeader(PMONSTARTREQ, msgNum, PMONMSGSD)
	if SEGMENT == c.Segment {
		eLen := c.Buffer - len(s)
		needSegment := int(math.Ceil(float64(len(d) / eLen)))
		var r []string
		off := 0
		for idx := 0; idx <= needSegment; idx++ {

			s := c.setPMONHeader(PMONSTARTREQ, msgNum+idx, PMONMSGSD)
			s += fmt.Sprintf("%04d", c.GetBlockCount()) // BlockCount
			residualLengh := len(d[off:])
			if residualLengh <= eLen {
				//一条完整的
				s += fmt.Sprintf("%04d", residualLengh) // effective byte length = Residual byte length
				s += fmt.Sprintf("%06d", residualLengh) // Residual byte length
				s += data[off:]
			} else {
				s += fmt.Sprintf("%04d", len(d[off:off+eLen])) // effective byte length
				s += fmt.Sprintf("%06d", residualLengh)        // Residual byte length
				s += data[off : off+eLen]
				off += eLen
			}
			s += PMONEND
			s += checksum(s)
			r = append(r, s) //添加到list中去
		}
		return r, nil
	} else {
		s += fmt.Sprintf("%04d", c.GetBlockCount()) // BlockCount
		s += fmt.Sprintf("%04d", len(d))            // byte length
		s += data
		s += PMONEND
		s += checksum(s)
		return []string{s}, nil
	}

}

func (c *Channel) PMONGenerateMsg(t PMONSMGTYPE, data string) ([]string, error) {
	msgid := c.conn.U.GetMsgNum()
	switch t {
	case PMONMSGSO:
		s, err := c.generateSO(msgid)
		return []string{s}, err
	case PMONMSGSC:
		s, err := c.generateSC(msgid)
		return []string{s}, err
	case PMONMSGAO:
		s, err := c.generateAO(msgid)
		return []string{s}, err
	case PMONMSGAC:
		s, err := c.generateAC(msgid)
		return []string{s}, err
	case PMONMSGSD:
		s, err := c.generateSD(msgid, data)
		return s, err
	case PMONMSGAD:
		err := fmt.Errorf("AD message is not support to generate")
		return []string{""}, err
	default:
		err := fmt.Errorf("message Type not found")
		return []string{""}, err
	}
}

func PMONParseMsg(buf []byte) PmonPackage {
	protocolElement := string(buf[13:15])
	switch protocolElement {
	case PMONMSGSO:
		return PmonPackage{
			t:    PMONMSGSO,
			data: nil,
		}
	case PMONMSGAO:
		return PmonPackage{
			t:    PMONMSGAO,
			data: nil,
		}
	case PMONMSGSC:
		return PmonPackage{
			t:    PMONMSGSC,
			data: nil,
		}
	case PMONMSGAC:
		return PmonPackage{
			t:    PMONMSGAC,
			data: nil,
		}
	case PMONMSGSD:
		return PmonPackage{
			t:    PMONMSGSD,
			data: buf[15:], //从block count 开始传送
		}
	case PMONMSGAD:
		return PmonPackage{
			t:    PMONMSGAD,
			data: buf[15:],
		}
	}
	return PmonPackage{
		t:    PMONMSGUN,
		data: nil,
	}
}
