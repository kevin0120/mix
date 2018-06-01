package pmon

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"github.com/masami10/aiis/services/httpd"
	"unsafe"
)

type TraceLogLevel int

const (
	PmonLogall TraceLogLevel = iota
	PmonLogaError
	PmonLogaMsg
)

const (
	MsgTypeSO = "SO"
	MsgTypeAO = "AO"
	MsgTypeSD = "SD"
	MsgTypeAD = "AD"
	MsgTypeSC = "SC"
	MsgTypeAC = "AC"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type EventHandler func(int, *C.char, *C.char)

type Service struct {
	configValue atomic.Value
	pmonLib     *syscall.DLL

	expatLib *syscall.DLL

	HTTPD *httpd.Service

	diag Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag: d,
	}
	s.configValue.Store(c)
	return s

}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

//请查阅 https://golang.org/cmd/cgo
func h(i int, ch *C.char, msg *C.char) {
	sCh := C.GoString(ch)
	sMsg := C.GoString(msg)
	fmt.Printf("%d, %s, %s", i, sCh, sMsg)
}

func (s *Service) Open() error {
	c := s.config()
	lib_pmon_path := filepath.Join(c.PmonDir, "pmon_lib.dll")
	lib_expat_path := filepath.Join(c.PmonDir, "libexpat.dll")
	lib_expat, err := syscall.LoadDLL(lib_expat_path)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("Load expat Lib fail: %s", lib_expat_path)
	}
	lib_pmon, err := syscall.LoadDLL(lib_pmon_path)
	if err != nil {
		return fmt.Errorf("Load Pmon Lib fail: %s", lib_pmon_path)
	}
	s.pmonLib = lib_pmon
	s.expatLib = lib_expat

	_, err = s.init(h, PmonLogall)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Close() error {
	defer s.pmonLib.Release()
	defer s.expatLib.Release()
	return nil
}

func (s *Service) init(e EventHandler, tracelevel TraceLogLevel) (bool, error) {
	if s.pmonLib == nil {
		return false, errors.New("Pmon lib is not current load")
	}
	proc, err := s.pmonLib.FindProc("PmonInit")

	if err != nil {
		return false, errors.New("fail to find the entrypoin PmonInit")
	}

	fn := syscall.NewCallback(e) //创建一个stdcall的回掉函数
	r, _, err := proc.Call(fn, uintptr(tracelevel))

	if err != nil {
		return false, err
	}
	if r != 0 {
		return false, fmt.Errorf("call PmonInit fail %d", r)
	}

	return false, nil
}

func (s *Service) SendData(msgId int, channelNumber string, data string) (bool, error) {
	if s.pmonLib == nil {
		return false, errors.New("Pmon lib is not current load")
	}
	proc, err := s.pmonLib.FindProc("SendData")

	if err != nil {
		return false, errors.New("fail to find the entrypoin SendData")
	}

	n := C.CString(channelNumber)
	defer C.free(unsafe.Pointer(n))

	d := C.CString(channelNumber)
	defer C.free(unsafe.Pointer(d))

	r, _, err := proc.Call(uintptr(msgId), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(d)))

	if err != nil {
		return false, err
	}
	if r != 0 {
		return false, fmt.Errorf("call PmonInit fail %d", r)
	}

	return false, nil
}

func (s *Service) SendPmonMessage(msgType string, channelNumber string, data string) (bool, error) {
	if s.pmonLib == nil {
		return false, errors.New("Pmon lib is not current load")
	}
	proc, err := s.pmonLib.FindProc("SendPmonMessage")

	if err != nil {
		return false, errors.New("fail to find the entrypoint SendPmonMessage")
	}

	m := C.CString(msgType)
	defer C.free(unsafe.Pointer(m))

	n := C.CString(channelNumber)
	defer C.free(unsafe.Pointer(n))

	d := C.CString(channelNumber)
	defer C.free(unsafe.Pointer(d))

	r, _, err := proc.Call(uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(d)))

	if err != nil {
		return false, err
	}
	if r != 0 {
		return false, fmt.Errorf("call PmonInit fail %d", r)
	}

	return false, nil
}
