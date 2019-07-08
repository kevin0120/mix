package scanner

import (
	"errors"
	"fmt"
)

const (
	VendorHoneyWell = 3118
	VendorDataLogic = 1529
)

type commonHoneywellScanner struct {
	cfg *USBConfig
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
	cfg, err := dev.Config(2)
	if err != nil {
		return err
	}
	v.cfg = cfg
	intf, err := cfg.Interface(3, 0)
	if err != nil {
		return err
	}
	v.Interface = intf
	epIn, err := intf.InEndpoint(6)
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

func (v *commonHoneywellScanner) Parse(buf []byte) string {
	fmt.Printf("%d\n", len(buf))
	return string(buf)
}

func (v *commonDataLogicScanner) Parse(buf []byte) string {
	return string(buf)
}
