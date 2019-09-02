package openprotocol

import (
	"errors"
	"fmt"
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
func (s *TighteningTool) SetPSet(pset int) error {
	rev, err := GetVendorMid(s.parent.Model(), MID_0018_PSET)
	if err != nil {
		return err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s.parent.Response.Add(MID_0018_PSET, nil)
	defer s.parent.Response.remove(MID_0018_PSET)

	s_pset := GeneratePackage(MID_0018_PSET, rev, "", "", "", fmt.Sprintf("%03d", pset))

	s.parent.Write([]byte(s_pset))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0018_PSET)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

// 设置Job
func (s *TighteningTool) SetJob(job int) error {
	rev, err := GetVendorMid(s.parent.Model(), MID_0038_JOB_SELECT)
	if err != nil {
		return err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s.parent.Response.Add(MID_0038_JOB_SELECT, nil)
	defer s.parent.Response.remove(MID_0038_JOB_SELECT)

	s_job := GeneratePackage(MID_0038_JOB_SELECT, rev, "", "", "", fmt.Sprintf("%04d", job))

	s.parent.Write([]byte(s_job))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0038_JOB_SELECT)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}
	return nil
}

// 模式选择: job/pset  0: pset 1: job
func (s *TighteningTool) ModeSelect(mode string) error {
	rev, err := GetVendorMid(s.parent.Model(), MID_0130_JOB_OFF)
	if err != nil {
		return err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s.parent.Response.Add(MID_0130_JOB_OFF, nil)
	defer s.parent.Response.remove(MID_0130_JOB_OFF)

	s_off := GeneratePackage(MID_0130_JOB_OFF, rev, "", "", "", mode)

	s.parent.Write([]byte(s_off))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0130_JOB_OFF)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	return nil
}

// 取消job
func (s *TighteningTool) AbortJob() error {
	rev, err := GetVendorMid(s.parent.Model(), MID_0127_JOB_ABORT)
	if err != nil {
		return err
	}

	if s.parent.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	s.parent.Response.Add(MID_0127_JOB_ABORT, nil)
	defer s.parent.Response.remove(MID_0127_JOB_ABORT)

	s_job := GeneratePackage(MID_0127_JOB_ABORT, rev, "", "", "", "")

	s.parent.Write([]byte(s_job))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0127_JOB_ABORT)
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

func (s *TighteningTool) Status() string {
	if s.parent.Status() == device.STATUS_OFFLINE {
		return device.STATUS_OFFLINE
	}

	return s.BaseDevice.Status()
}

func (s *TighteningTool) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_TOOL
}
