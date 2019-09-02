package openprotocol

import (
	"errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/tightening_device"
	"time"
)

func NewTool(c tightening_device.ITighteningController, cfg tightening_device.ToolConfig, d Diagnostic) TighteningTool {
	tool := TighteningTool{
		diag:       d,
		cfg:        cfg,
		parent:     c.(*TighteningController),
		BaseDevice: device.CreateBaseDevice(),
	}

	return tool
}

type TighteningTool struct {
	diag Diagnostic
	cfg  tightening_device.ToolConfig
	//status atomic.Value
	parent *TighteningController

	device.BaseDevice
}

// 工具使能控制
func (s *TighteningTool) ToolControl(enable bool) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	cmd := MID_0042_TOOL_DISABLE
	if enable {
		cmd = MID_0043_TOOL_ENABLE
	}

	rev, err := GetVendorMid(s.parent.Model(), cmd)
	if err != nil {
		return err
	}

	s.parent.Response.Add(cmd, nil)
	defer s.parent.Response.remove(cmd)

	sSend := GeneratePackage(cmd, rev, "", "", "", "")

	s.parent.Write([]byte(sSend))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(cmd)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	sReply := reply.(string)
	if sReply != request_errors["00"] {
		return errors.New(sReply)
	}

	return nil
}

// 设置PSet
func (s *TighteningTool) SetPSet(pset int) error { return nil }

// 设置Job
func (s *TighteningTool) SetJob(job int) error { return nil }

// 模式选择: job/pset
func (s *TighteningTool) ModeSelect(mode string) error { return nil }

// 取消job
func (s *TighteningTool) AbortJob() error { return nil }

func (s *TighteningTool) Status() string {
	if s.parent.Status() == device.STATUS_OFFLINE {
		return device.STATUS_OFFLINE
	}

	return s.BaseDevice.Status()
}

func (s *TighteningTool) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_TOOL
}
