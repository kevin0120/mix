package scanner

import (
	"fmt"
	"github.com/karalabe/hid"
	"testing"
	//"time"
)

func TestScanner(t *testing.T) {
	fmt.Println(hid.Supported())

	devinfos := hid.Enumerate(3118, 2311)
	devinfo := devinfos[0]

	dev, err := devinfo.Open()
	if err != nil {
		fmt.Printf("open failed:%s\n", err.Error())
		return
	}

	buf := make([]byte, 128)
	for {
		n, err := dev.Read(buf)
		if err != nil {
			fmt.Printf("read failed:%s\n", err.Error())
			break
		}

		if n > 0 {
			fmt.Printf("read:%s\n", string(buf)[0:n])
		}
	}

}
