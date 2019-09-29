package openprotocol

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
	"sync/atomic"
)

const (
	OPENPROTOCOL_MODE_JOB  = "1"
	OPENPROTOCOL_MODE_PSET = "0"
)

func CreateTool(c *TighteningController, cfg tightening_device.ToolConfig, d Diagnostic) *TighteningTool {
	tool := TighteningTool{
		diag:       d,
		cfg:        cfg,
		controller: c,
		BaseDevice: device.CreateBaseDevice(),
	}

	tool.UpdateStatus(device.STATUS_ONLINE)
	return &tool
}

type TighteningTool struct {
	diag       Diagnostic
	cfg        tightening_device.ToolConfig
	Mode       atomic.Value
	controller *TighteningController

	device.BaseDevice
}

func (s *TighteningTool) SetMode(mode string) {
	s.Mode.Store(mode)

	//s.controller.Srv.DB.UpdateTool(&storage.Guns{
	//	Serial: s.cfg.SN,
	//	Mode:   mode,
	//})
}

func (s *TighteningTool) GetMode() string {
	return s.Mode.Load().(string)
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

	reply, err := s.controller.ProcessRequest(cmd, "", "", "", "")
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
	reply, err := s.controller.ProcessRequest(MID_0018_PSET, "", "", "", data)
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
	reply, err := s.controller.ProcessRequest(MID_0038_JOB_SELECT, "", "", "", data)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

// 模式选择: job/pset
func (s *TighteningTool) ModeSelect(mode string) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	flag := OPENPROTOCOL_MODE_PSET
	if mode == tightening_device.MODE_JOB {
		flag = OPENPROTOCOL_MODE_JOB
	}

	reply, err := s.controller.ProcessRequest(MID_0130_JOB_OFF, "", "", "", flag)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	s.SetMode(mode)

	return nil
}

// 取消job
func (s *TighteningTool) AbortJob() error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	reply, err := s.controller.ProcessRequest(MID_0127_JOB_ABORT, "", "", "", "")
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
	reply, err := s.controller.ProcessRequest(MID_0019_PSET_BATCH_SET, "", "", "", data)
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

	reply, err := s.controller.ProcessRequest(MID_0010_PSET_LIST_REQUEST, "", "", "", "")
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
	reply, err := s.controller.ProcessRequest(MID_0012_PSET_DETAIL_REQUEST, "", "", "", data)
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

	reply, err := s.controller.ProcessRequest(MID_0030_JOB_LIST_REQUEST, "", "", "", "")
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
	reply, err := s.controller.ProcessRequest(MID_0032_JOB_DETAIL_REQUEST, "", "", "", data)
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

func (s *TighteningTool) TraceSet(str string) error {
	if s.Status() == device.STATUS_OFFLINE {
		return errors.New(device.STATUS_OFFLINE)
	}

	id := s.controller.Srv.generateIDInfo(str)
	reply, err := s.controller.ProcessRequest(MID_0150_IDENTIFIER_SET, "", "", "", id)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

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
	if s.controller.Status() == device.STATUS_OFFLINE {
		return device.STATUS_OFFLINE
	}

	return s.BaseDevice.Status()
}

func (s *TighteningTool) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_TOOL
}

// 处理结果
func (s *TighteningTool) OnResult(result interface{}) {
	if result == nil {
		s.diag.Error(fmt.Sprintf("Tool SN: %s", s.cfg.SN), errors.New("Result Is Nil"))
		return
	}

	tighteningResult := result.(*tightening_device.TighteningResult)
	dbTool, err := s.controller.Srv.DB.GetGun(s.cfg.SN)
	if err == nil {
		if s.Mode.Load().(string) == tightening_device.MODE_JOB {
			tighteningResult.Seq, tighteningResult.Count = s.controller.calBatch(dbTool.WorkorderID)
		} else {
			tighteningResult.Seq = dbTool.Seq
			tighteningResult.Count = dbTool.Count
		}

		tighteningResult.WorkorderID = dbTool.WorkorderID
		tighteningResult.UserID = dbTool.UserID
		tighteningResult.Batch = fmt.Sprintf("%d/%d", tighteningResult.Seq, dbTool.Total)

		dbWorkorder, err := s.controller.Srv.DB.GetWorkorder(dbTool.WorkorderID, true)
		if err != nil {
			s.diag.Error("Get Workorder Failed", err)
			return
		}

		consume, err := s.controller.Srv.Odoo.GetConsumeBySeq(&dbWorkorder, tighteningResult.Seq)
		if err != nil {
			s.diag.Error("Get Consume Failed", err)
			return
		}

		tighteningResult.NutNo = consume.NutNo
	}

	dbResult := tighteningResult.ToDBResult()

	// 尝试获取最近一条没有对应结果的曲线并更新, 同时缓存结果
	err = s.controller.Srv.DB.UpdateIncompleteCurveAndSaveResult(dbResult)
	if err != nil {
		s.diag.Error("Handle Result With Curve Failed", err)
	}

	// 分发结果
	s.controller.GetDispatch(tightening_device.DISPATCH_RESULT).Dispatch(tighteningResult)
}

// 处理曲线
func (s *TighteningTool) OnCurve(curve interface{}) {
	if curve == nil {
		s.diag.Error(fmt.Sprintf("Tool SN: %s", s.cfg.SN), errors.New("Curve Is Nil"))
		return
	}

	tighteningCurve := curve.(*tightening_device.TighteningCurve)
	dbCurves := tighteningCurve.ToDBCurve()

	// 尝试获取最近一条没有对应曲线的结果并更新, 同时缓存曲线
	err := s.controller.Srv.DB.UpdateIncompleteResultAndSaveCurve(dbCurves)
	if err != nil {
		s.diag.Error("Handle Curve With Result Failed", err)
	}
}

func (s *TighteningTool) GetDispatch(name string) *utils.Dispatcher {
	return nil
}
