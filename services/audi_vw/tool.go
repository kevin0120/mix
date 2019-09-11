package audi_vw

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
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
	return nil
}

// 设置PSet
func (s *TighteningTool) SetPSet(pset int) error {
	return nil
}

// 设置Job
func (s *TighteningTool) SetJob(job int) error {
	return errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// 模式选择: job/pset  0: pset 1: job
func (s *TighteningTool) ModeSelect(mode string) error {
	return errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// 取消job
func (s *TighteningTool) AbortJob() error {
	return errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// 设置pset次数
func (s *TighteningTool) SetPSetBatch(pset int, batch int) error {
	return errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// pset列表
func (s *TighteningTool) GetPSetList() ([]int, error) {
	return nil, errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// pset详情
func (s *TighteningTool) GetPSetDetail(pset int) (*tightening_device.PSetDetail, error) {
	return nil, errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// job列表
func (s *TighteningTool) GetJobList() ([]int, error) {
	return nil, errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
}

// job详情
func (s *TighteningTool) GetJobDetail(job int) (*tightening_device.JobDetail, error) {
	return nil, errors.New(tightening_device.TIGHTENING_ERR_NOT_SUPPORTED)
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

// 处理结果
func (c *TighteningTool) OnResult(result interface{}) {
	if result == nil {
		c.diag.Error(fmt.Sprintf("Tool SN: %s", c.cfg.SN), errors.New("Result Is Nil"))
		return
	}

}

// 处理曲线
func (c *TighteningTool) OnCurve(curve interface{}) {
	if curve == nil {
		c.diag.Error(fmt.Sprintf("Tool SN: %s", c.cfg.SN), errors.New("Curve Is Nil"))
		return
	}
}

func (c *TighteningTool) GetDispatch(name string) *utils.Dispatcher {
	return nil
}
