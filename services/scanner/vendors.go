package scanner

import (
	"errors"
)

const (
	VendorHoneyWell  = 3118
	ProductHoneyWell = 2305
	VendorDataLogic  = 1529
)

const (
	KEY_SHIFT    = 2
	INDEX_TARGET = 2
	INDEX_SHIFT  = 0
)

var keymap = map[byte][]string{
	4:  {"a", "A"},
	5:  {"b", "B"},
	6:  {"c", "C"},
	7:  {"d", "D"},
	8:  {"e", "E"},
	9:  {"f", "F"},
	10: {"g", "G"},
	11: {"h", "H"},
	12: {"i", "I"},
	13: {"j", "J"},
	14: {"k", "K"},
	15: {"l", "L"},
	16: {"m", "M"},
	17: {"n", "N"},
	18: {"o", "O"},
	19: {"p", "P"},
	20: {"q", "Q"},
	21: {"r", "R"},
	22: {"s", "S"},
	23: {"t", "T"},
	24: {"u", "U"},
	25: {"v", "V"},
	26: {"w", "W"},
	27: {"x", "X"},
	28: {"y", "Y"},
	29: {"z", "Z"},
	30: {"1", ""},
	31: {"2", ""},
	32: {"3", ""},
	33: {"4", ""},
	34: {"5", ""},
	35: {"6", ""},
	36: {"7", ""},
	37: {"8", ""},
	38: {"9", ""},
	39: {"0", ""},
}

type vendor struct {
	VendorID  ID
	ProductID ID
}

var vendors = []vendor{
	{
		VendorID:  VendorHoneyWell,
		ProductID: ProductHoneyWell,
	},
}

type commonHoneywellScanner struct {
	cfg        *USBConfig
	Interface  *USBInterface
	InEndpoint *USBInEndpoint
}

type commonDataLogicScanner struct {
}

func (d *DeviceInfo) updateDeviceService() {
	switch d.VendorID {
	case VendorHoneyWell:
		d.DeviceService = &commonHoneywellScanner{}
		//case VendorDataLogic:
		//	d.DeviceService = &commonDataLogicScanner{}
	}
}

func (v *commonHoneywellScanner) NewReader(dev *USBDevice) error {
	cfg, err := dev.Config(1)
	if err != nil {
		return err
	}
	v.cfg = cfg
	intf, err := cfg.Interface(0, 0)
	if err != nil {
		return err
	}
	v.Interface = intf
	epIn, err := intf.InEndpoint(4)
	if err != nil {
		return err
	}
	v.InEndpoint = epIn

	return nil
}

func (v *commonHoneywellScanner) Read(buf []byte) (int, error) {
	if v.InEndpoint == nil {
		return 0, errors.New("not Reader")
	}
	return v.InEndpoint.Read(buf)
}

func (v *commonHoneywellScanner) Close() error {
	var err error
	if v.cfg != nil {
		err = v.cfg.Close()
	}
	if v.Interface != nil {
		v.Interface.Close()
	}
	return err
}

func (v *commonHoneywellScanner) Parse(buf []byte) (string, error) {
	//fmt.Println(buf)
	//fmt.Printf("%d\n", len(buf))
	if buf[INDEX_TARGET] == 0 {
		return "", errors.New("invalid byte")
	}

	str := keymap[buf[INDEX_TARGET]][0]
	if buf[INDEX_SHIFT] == KEY_SHIFT {
		str = keymap[buf[INDEX_TARGET]][1]
	}

	return str, nil
}

func (v *commonDataLogicScanner) Parse(buf []byte) string {
	return string(buf)
}
