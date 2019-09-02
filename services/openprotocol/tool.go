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
	diag   Diagnostic
	cfg    tightening_device.ToolConfig
	Mode   string
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

	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
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

	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
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

	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	if mode != tightening_device.MODE_JOB && mode != tightening_device.MODE_PSET {
		return errors.New("Mode Error")
	}

	s.parent.Response.Add(MID_0130_JOB_OFF, nil)
	defer s.parent.Response.remove(MID_0130_JOB_OFF)

	flag := "0"
	if mode == tightening_device.MODE_JOB {
		flag = "1"
	}

	s_off := GeneratePackage(MID_0130_JOB_OFF, rev, "", "", "", flag)

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

	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	if s.Mode != tightening_device.MODE_JOB {
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

// 设置pset次数
func (s *TighteningTool) SetPSetBatch(pset int, batch int) error {
	rev, err := GetVendorMid(s.parent.Model(), MID_0019_PSET_BATCH_SET)
	if err != nil {
		return err
	}

	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	s.parent.Response.Add(MID_0019_PSET_BATCH_SET, nil)
	defer s.parent.Response.remove(MID_0019_PSET_BATCH_SET)

	data := fmt.Sprintf("%03d%02d", pset, batch)
	ide := GeneratePackage(MID_0019_PSET_BATCH_SET, rev, "", "", "", data)

	s.parent.Write([]byte(ide))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0019_PSET_BATCH_SET)
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

// pset列表
func (s *TighteningTool) GetPSetList() ([]int, error) {
	var psets []int

	rev, err := GetVendorMid(s.parent.Model(), MID_0010_PSET_LIST_REQUEST)
	if err != nil {
		return psets, err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return psets, errors.New(controller.STATUS_OFFLINE)
	}

	defer s.parent.Response.remove(MID_0010_PSET_LIST_REQUEST)
	s.parent.Response.Add(MID_0010_PSET_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0010_PSET_LIST_REQUEST, rev, "", "", "", "")
	s.parent.Write([]byte(psets_request))

	var reply interface{} = nil

	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0010_PSET_LIST_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return psets, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	pset_list := reply.(PSetList)

	return pset_list.psets, nil
}

// pset详情
func (s *TighteningTool) GetPSetDetail(pset int) (tightening_device.PSetDetail, error) {
	var obj_pset_detail tightening_device.PSetDetail

	rev, err := GetVendorMid(s.parent.Model(), MID_0012_PSET_DETAIL_REQUEST)
	if err != nil {
		return obj_pset_detail, err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return obj_pset_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer s.parent.Response.remove(MID_0012_PSET_DETAIL_REQUEST)
	s.parent.Response.Add(MID_0012_PSET_DETAIL_REQUEST, nil)

	pset_detail := GeneratePackage(MID_0012_PSET_DETAIL_REQUEST, rev, "", "", "", fmt.Sprintf("%03d", pset))
	s.parent.Write([]byte(pset_detail))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0012_PSET_DETAIL_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return obj_pset_detail, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	switch v := reply.(type) {
	case string:
		return obj_pset_detail, errors.New(v)

	case tightening_device.PSetDetail:
		return reply.(tightening_device.PSetDetail), nil

	default:
		return obj_pset_detail, errors.New(controller.ERR_KNOWN)
	}
}

// job列表
func (s *TighteningTool) GetJobList() ([]int, error) {
	var jobs []int
	rev, err := GetVendorMid(s.parent.Model(), MID_0030_JOB_LIST_REQUEST)
	if err != nil {
		return jobs, err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return jobs, errors.New(controller.STATUS_OFFLINE)
	}

	defer s.parent.Response.remove(MID_0030_JOB_LIST_REQUEST)
	s.parent.Response.Add(MID_0030_JOB_LIST_REQUEST, nil)

	psetsRequest := GeneratePackage(MID_0030_JOB_LIST_REQUEST, rev, "", "", "", "")
	s.parent.Write([]byte(psetsRequest))

	var reply interface{} = nil

	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0030_JOB_LIST_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return jobs, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	jobList := reply.(JobList)

	return jobList.jobs, nil
}

// job详情
func (s *TighteningTool) GetJobDetail(job int) (tightening_device.JobDetail, error) {
	var objJobDetail tightening_device.JobDetail
	rev, err := GetVendorMid(s.parent.Model(), MID_0032_JOB_DETAIL_REQUEST)
	if err != nil {
		return objJobDetail, err
	}

	if s.parent.Status() == controller.STATUS_OFFLINE {
		return objJobDetail, errors.New(controller.STATUS_OFFLINE)
	}

	defer s.parent.Response.remove(MID_0032_JOB_DETAIL_REQUEST)
	s.parent.Response.Add(MID_0032_JOB_DETAIL_REQUEST, nil)

	job_detail := GeneratePackage(MID_0032_JOB_DETAIL_REQUEST, rev, "", "", "", fmt.Sprintf("%04d", job))
	s.parent.Write([]byte(job_detail))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = s.parent.Response.get(MID_0032_JOB_DETAIL_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return objJobDetail, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	switch v := reply.(type) {
	case string:
		return objJobDetail, errors.New(v)

	case tightening_device.JobDetail:
		return reply.(tightening_device.JobDetail), nil

	default:
		return objJobDetail, errors.New(controller.ERR_KNOWN)
	}
}

//func (s *TighteningTool) IdentifierSet(str string) error {
//	rev, err := GetVendorMid(c.Model(), MID_0150_IDENTIFIER_SET)
//	if err != nil {
//		return err
//	}
//
//	if c.Status() == controller.STATUS_OFFLINE {
//		return errors.New("status offline")
//	}
//
//	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, rev, "", "", "", str)
//
//	c.Write([]byte(ide))
//
//	return nil
//}

//func (s *TighteningTool) PSetBatchReset(pset int) error {
//	rev, err := GetVendorMid(c.Model(), MID_0020_PSET_BATCH_RESET)
//	if err != nil {
//		return err
//	}
//
//	if c.Status() == controller.STATUS_OFFLINE {
//		return errors.New("status offline")
//	}
//
//	s := fmt.Sprintf("%03d", pset)
//	ide := GeneratePackage(MID_0020_PSET_BATCH_RESET, rev, "", "", "", s)
//
//	c.Write([]byte(ide))
//
//	return nil
//}

func (s *TighteningTool) Status() string {
	if s.parent.Status() == device.STATUS_OFFLINE {
		return device.STATUS_OFFLINE
	}

	return s.BaseDevice.Status()
}

func (s *TighteningTool) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_TOOL
}
