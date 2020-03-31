package scanner

import "github.com/google/gousb"

//type USBDevice = gousb.BaseDevice

type USBDevice interface {
	//SetAutoDetach(bool ) error
	Close() error
}
type ID = gousb.ID
type USBConfig = gousb.Config
type USBInterface = gousb.Interface
type USBInEndpoint = gousb.InEndpoint
