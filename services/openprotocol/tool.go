package openprotocol

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/tightening_device"
	"go.uber.org/atomic"
)

func NewTool(c *TighteningController, cfg tightening_device.ToolConfig, d Diagnostic) *TighteningTool {
	tool := TighteningTool{
		diag:       d,
		cfg:        cfg,
		controller: c,
		BaseDevice: device.CreateBaseDevice(device.BaseDeviceTighteningTool, d, c),
	}
	tool.SetSerialNumber(cfg.SN)
	tool.BaseDevice.UpdateStatus(device.BaseDeviceStatusOffline)
	tool.SetMode(c.ProtocolService.GetDefaultMode())
	return &tool
}

type TighteningTool struct {
	device.BaseDevice
	diag       Diagnostic
	cfg        tightening_device.ToolConfig
	mode       atomic.Value
	controller *TighteningController
}

func (s *TighteningTool) SetMode(mode string) {
	s.mode.Store(mode)
}

func (s *TighteningTool) Mode() string {
	return s.mode.Load().(string)
}

// 工具使能控制
func (s *TighteningTool) ToolControl(enable bool) error {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	cmd := MID_0042_TOOL_DISABLE
	if enable {
		cmd = MID_0043_TOOL_ENABLE
	}

	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(cmd, "", "", "", "")
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	data := fmt.Sprintf("%03d", pset)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0018_PSET, "", "", "", data)
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	data := fmt.Sprintf("%04d", job)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0038_JOB_SELECT, "", "", "", data)
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	flag := OpenprotocolModePset
	if mode == tightening_device.MODE_JOB {
		flag = OpenprotocolModeJob
	}

	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0130_JOB_OFF, "", "", "", flag)
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0127_JOB_ABORT, "", "", "", "")
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	data := fmt.Sprintf("%03d%02d", pset, batch)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0019_PSET_BATCH_SET, "", "", "", data)
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
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0010_PSET_LIST_REQUEST, "", "", "", "")
	if err != nil {
		return nil, err
	}

	return reply.(PSetList).psets, nil
}

// pset详情
func (s *TighteningTool) GetPSetDetail(pset int) (*tightening_device.PSetDetail, error) {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	data := fmt.Sprintf("%03d", pset)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0012_PSET_DETAIL_REQUEST, "", "", "", data)
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

	return nil, errors.New(tightening_device.TIGHTENING_ERR_UNKNOWN)
}

// job列表
func (s *TighteningTool) GetJobList() ([]int, error) {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0030_JOB_LIST_REQUEST, "", "", "", "")
	if err != nil {
		return nil, err
	}

	return reply.(JobList).jobs, nil
}

// job详情
func (s *TighteningTool) GetJobDetail(job int) (*tightening_device.JobDetail, error) {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	data := fmt.Sprintf("%04d", job)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0032_JOB_DETAIL_REQUEST, "", "", "", data)
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

	return nil, errors.New(tightening_device.TIGHTENING_ERR_UNKNOWN)
}

func (s *TighteningTool) TraceSet(str string) error {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return errors.New(device.BaseDeviceStatusOffline)
	}

	id := s.controller.ProtocolService.generateIDInfo(str)
	reply, err := s.controller.getClient(s.SerialNumber()).ProcessRequest(MID_0150_IDENTIFIER_SET, "", "", "", id)
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
//	if c.Status() == controller.BaseDeviceStatusOffline {
//		return errors.New("status offline")
//	}
//
//	s := fmt.Sprintf("%03d", pset)
//	ide := GeneratePackage(MID_0020_PSET_BATCH_RESET, rev, "", "", "", s)
//
//	c.IOWrite([]byte(ide))
//
//	return nil
//}

func (s *TighteningTool) Status() string {
	if s.controller.Status() == device.BaseDeviceStatusOffline {
		return device.BaseDeviceStatusOffline
	}

	return s.BaseDevice.Status()
}

func (s *TighteningTool) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_TOOL
}

// 处理结果
func (s *TighteningTool) onResult(result interface{}) {
	if result == nil {
		s.diag.Error(fmt.Sprintf("Tool SerialNumber: %s", s.cfg.SN), errors.New("Result Is Nil "))
		return
	}

	tighteningResult := result.(tightening_device.TighteningResult)
	dbTool, err := s.controller.ProtocolService.storageService.GetTool(s.cfg.SN)
	if err == nil && dbTool.CurrentWorkorderID != 0 {
		if s.Mode() == tightening_device.MODE_JOB {
			tighteningResult.Seq, tighteningResult.Count = s.controller.calBatch(dbTool.CurrentWorkorderID)
		} else {
			tighteningResult.Seq = dbTool.Seq
			tighteningResult.Count = dbTool.Count
		}

		tighteningResult.WorkorderID = dbTool.CurrentWorkorderID
		tighteningResult.UserID = dbTool.UserID
		tighteningResult.Batch = fmt.Sprintf("%d/%d", tighteningResult.Seq, dbTool.Total)

		dbStep, err := s.controller.ProtocolService.storageService.GetStep(dbTool.StepID)
		if err != nil {
			s.diag.Error("Get Step Failed", err)
			return
		}

		consume, err := s.controller.ProtocolService.backendService.GetConsumeBySeqInStep(&dbStep, tighteningResult.Seq)
		if err != nil {
			s.diag.Error("Get Consume Failed", err)
			return
		}

		tighteningResult.NutNo = consume.NutNo
	}

	dbResult := tighteningResult.ToDBResult()

	// 尝试获取最近一条没有对应结果的曲线并更新, 同时缓存结果
	err = s.controller.ProtocolService.storageService.UpdateIncompleteCurveAndSaveResult(dbResult)
	if err != nil {
		s.diag.Error("Handle Result With Curve Failed", err)
	}

	// 分发结果
	tighteningResult.ID = dbResult.Id
	s.controller.doDispatch(dispatcherbus.DispatcherResult, tighteningResult)
}

// 处理曲线
func (s *TighteningTool) onCurve(curve interface{}) {
	if curve == nil {
		s.diag.Error(fmt.Sprintf("Tool SerialNumber: %s", s.cfg.SN), errors.New("Curve Is Nil "))
		return
	}

	tighteningCurve := curve.(*tightening_device.TighteningCurve)
	dbCurves := tighteningCurve.ToDBCurve()

	// 尝试获取最近一条没有对应曲线的结果并更新, 同时缓存曲线
	err := s.controller.ProtocolService.storageService.UpdateIncompleteResultAndSaveCurve(dbCurves)
	if err != nil {
		s.diag.Error("Handle Curve With Result Failed", err)
	} else {
		s.diag.Info(fmt.Sprintf("缓存曲线成功 SN:%s", s.cfg.SN))
	}

	// 分发曲线
	s.controller.doDispatch(dispatcherbus.DispatcherCurve, tighteningCurve)
	s.diag.Info(fmt.Sprintf("缓存曲线成功 工具:%s 对应拧紧ID:%s", dbCurves.ToolSN, dbCurves.TighteningID))
}

func (s *TighteningTool) UpdateStatus(status string) {

	s.BaseDevice.UpdateStatus(status)
	toolStatus := []device.Status{{
		Type:   tightening_device.TIGHTENING_DEVICE_TYPE_TOOL,
		SN:     s.SerialNumber(),
		Status: status,
	}}

	s.controller.doDispatch(dispatcherbus.DispatcherDeviceStatus, toolStatus)
}

func (s *TighteningTool) Config() interface{} {
	return s.cfg
}
