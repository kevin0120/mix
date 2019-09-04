package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/ascii"
	"testing"
)

const (
	TEST_STRINGS = "123" + "3.14" + "7410" + "15.3" + "-897" + "99" + "1" + "0" + "sn001"+"666"
)

type Header struct {
	TOOL string `start:"1"  end:"5"`
	Sn   int	`start:"6"  end:"8"`
}

type OpenProtocol struct {
	L   int64   `start:"1"  end:"3"`
	LEN float32 `start:"4"  end:"7"`
	MID string  `start:"8"  end:"11"`
	MD  float64 `start:"12" end:"15"`
	M   int     `start:"16" end:"19"`
	Faa uint    `start:"20" end:"21"`
	B   bool    `start:"22" end:"22"`
	C   bool    `start:"23" end:"23"`
	TO  Header  `start:"24" end:"..."`
}

func Test_Ascii(t *testing.T){
//var he =Header{}
	var testop = OpenProtocol{}
	err := ascii.Unmarshal(TEST_STRINGS, &testop)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", testop)
}
