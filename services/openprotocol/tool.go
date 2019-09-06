package openprotocol

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
)

const (
	OPENPROTOCOL_MODE_JOB  = "1"
	OPENPROTOCOL_MODE_PSET = "0"
)

func NewTool(c *TighteningController, cfg tightening_device.ToolConfig, d Diagnostic) *TighteningTool {
	tool := TighteningTool{
		diag:       d,
		cfg:        cfg,
		parent:     c,
		BaseDevice: device.CreateBaseDevice(),
	}

	tool.UpdateStatus(device.STATUS_ONLINE)
	return &tool
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

	reply, err := s.parent.ProcessRequest(cmd, "", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 设置PSet
func (s *TighteningTool) SetPSet(pset int) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	data := fmt.Sprintf("%03d", pset)
	reply, err := s.parent.ProcessRequest(MID_0018_PSET, "", "", "", data)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 设置Job
func (s *TighteningTool) SetJob(job int) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	data := fmt.Sprintf("%04d", job)
	reply, err := s.parent.ProcessRequest(MID_0038_JOB_SELECT, "", "", "", data)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 模式选择: job/pset  0: pset 1: job
func (s *TighteningTool) ModeSelect(mode string) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	flag := OPENPROTOCOL_MODE_PSET
	if mode == tightening_device.MODE_JOB {
		flag = OPENPROTOCOL_MODE_JOB
	}

	reply, err := s.parent.ProcessRequest(MID_0130_JOB_OFF, "", "", "", flag)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 取消job
func (s *TighteningTool) AbortJob() error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	reply, err := s.parent.ProcessRequest(MID_0127_JOB_ABORT, "", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 设置pset次数
func (s *TighteningTool) SetPSetBatch(pset int, batch int) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	data := fmt.Sprintf("%03d%02d", pset, batch)
	reply, err := s.parent.ProcessRequest(MID_0019_PSET_BATCH_SET, "", "", "", data)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// pset列表
func (s *TighteningTool) GetPSetList() ([]int, error) {
	if s.Status() == device.STATUS_OFFLINE {
		return nil, errors.New(device.STATUS_OFFLINE)
	}

	reply, err := s.parent.ProcessRequest(MID_0010_PSET_LIST_REQUEST, "", "", "", "")
	if err != nil {
		return nil, err
	}

	return reply.(PSetList).psets, nil
}

// pset详情
func (s *TighteningTool) GetPSetDetail(pset int) (*tightening_device.PSetDetail, error) {
	if s.Status() == device.STATUS_OFFLINE {
		return nil, errors.New(device.STATUS_OFFLINE)
	}

	data := fmt.Sprintf("%03d", pset)
	reply, err := s.parent.ProcessRequest(MID_0012_PSET_DETAIL_REQUEST, "", "", "", data)
	if err != nil {
		return nil, err
	}

	switch v := reply.(type) {
	case string:
		return nil, errors.New(v)

	case tightening_device.PSetDetail:
		rt := reply.(tightening_device.PSetDetail)
		return &rt, nil
	}

	return nil, errors.New(controller.ERR_KNOWN)
}

// job列表
func (s *TighteningTool) GetJobList() ([]int, error) {
	if s.Status() == device.STATUS_OFFLINE {
		return nil, errors.New(device.STATUS_OFFLINE)
	}

	reply, err := s.parent.ProcessRequest(MID_0030_JOB_LIST_REQUEST, "", "", "", "")
	if err != nil {
		return nil, err
	}

	return reply.(JobList).jobs, nil
}

// job详情
func (s *TighteningTool) GetJobDetail(job int) (*tightening_device.JobDetail, error) {
	if s.Status() == device.STATUS_OFFLINE {
		return nil, errors.New(device.STATUS_OFFLINE)
	}

	data := fmt.Sprintf("%04d", job)
	reply, err := s.parent.ProcessRequest(MID_0032_JOB_DETAIL_REQUEST, "", "", "", data)
	if err != nil {
		return nil, err
	}

	switch v := reply.(type) {
	case string:
		return nil, errors.New(v)

	case tightening_device.JobDetail:
		rt := reply.(tightening_device.JobDetail)
		return &rt, nil
	}

	return nil, errors.New(controller.ERR_KNOWN)
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

// 处理结果
func (c *TighteningTool) OnResult(result interface{}) {
	if result == nil {
		c.diag.Error(fmt.Sprintf("Tool SN: %s", c.cfg.SN), errors.New("Result Is Nil"))
		return
	}

	tighteningResult := result.(*tightening_device.TighteningResult)
	dbResult := tighteningResult.ToDBResult()

	// 尝试获取最近一条没有对应结果的曲线并更新
	err := c.parent.Srv.DB.UpdateIncompleteCurve(dbResult.ToolSN, dbResult.TighteningID)
	if err == nil {
		c.diag.Debug("No Curve Need Update")
	}

	// 缓存结果
	err = c.parent.Srv.DB.Store(dbResult)
	if err != nil {
		c.diag.Debug("Save Result Failed")
	}

	// 分发结果
	c.parent.externalResultDispatch.Dispatch(tighteningResult)
}

// 处理曲线
func (c *TighteningTool) OnCurve(curve interface{}) {
	if curve == nil {
		c.diag.Error(fmt.Sprintf("Tool SN: %s", c.cfg.SN), errors.New("Curve Is Nil"))
		return
	}

	// TODO
	//tighteningCurve := curve.(*tightening_device.TighteningCurve)

	// 尝试获取最近一条没有对应曲线的结果并更新， 如果成功则上传曲线， 否则只缓存

	// 缓存曲线
}

func (c *TighteningTool) GetDispatch(name string) *utils.Dispatch {
	return nil
}
