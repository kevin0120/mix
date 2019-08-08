package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	name        string

	DB    *storage.Service
	WS    *wsnotify.Service
	Aiis  *aiis.Service
	Minio *minio.Service
	Odoo  *odoo.Service

	Parent *controller.Service

	devices []Controller
}

func NewService(c Config, d Diagnostic, parent *controller.Service) *Service {

	s := &Service{
		name:    controller.OPENPROTOCOL,
		diag:    d,
		Parent:  parent,
		devices: []Controller{},
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (p *Service) Parse(msg string) ([]byte, error) {
	return nil, nil
}

func (p *Service) Write(sn string, buf []byte) error {
	return nil
}

func (p *Service) AddNewController(cfg controller.ControllerConfig) controller.Controller {
	config := p.config()
	c := NewController(config, p.diag)
	c.Srv = p //服务注入
	c.cfg = cfg

	return &c
}

func (p *Service) AddDevice(cfg controller.DeviceConfig, ts interface{}) controller.Controller {
	config := p.config()
	c := NewController(config, p.diag)
	c.Srv = p //服务注入
	c.cfg = controller.ControllerConfig{
		RemoteIP: cfg.Endpoint,
	}
	c.SetModel(cfg.Model)
	c.tighteningDevice = ts.(*tightening_device.Service)

	p.devices = append(p.devices, c)

	return nil
}

func (p *Service) Open() error {
	//for _, w := range p.Parent.Controllers {
	//	if w.Protocol() == controller.OPENPROTOCOL {
	//		go w.Start() //异步启动控制器
	//	}
	//}

	for _, v := range p.devices {
		go v.Start()
	}

	return nil
}

func (p *Service) Close() error {
	for _, w := range p.Parent.Controllers {
		err := w.Close()
		if err != nil {
			return errors.Wrapf(err, "Close Protocol %s Writer fail", p.name)
		}
	}

	return nil
}

func (p *Service) ToolControl(sn string, tool_sn string, enable bool) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	var t controller.ToolConfig

	var toolExist = false

	for _, t = range c.cfg.Tools {
		if t.SerialNO == tool_sn {
			toolExist = true
			break
		}
	}

	if !toolExist {
		return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	}

	// 工具使能
	err := c.ToolControl(enable)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil
}

func (p *Service) PSet(sn string, tool_sn string, pset int, result_id int64, count int, user_id int64) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	var t controller.ToolConfig

	var toolExist = false

	for _, t = range c.cfg.Tools {
		if t.SerialNO == tool_sn {
			toolExist = true
			break
		}
	}

	if !toolExist {
		return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	}

	ex_info := fmt.Sprintf("%s-%d-%d-%d", controller.AUTO_MODE, result_id, count, user_id)

	// 设定pset并判断控制器响应
	_, err := c.PSet(pset, t.ToolChannel, ex_info, 1)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil
}

func (p *Service) PSetManual(sn string, tool_sn string, pset int, user_id int64, ex_info string, count int) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	var t controller.ToolConfig

	var toolExist = false

	for _, t = range c.cfg.Tools {
		if t.SerialNO == tool_sn {
			toolExist = true
			break
		}
	}

	if !toolExist {
		return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	}

	// 设定pset并判断控制器响应
	_, err := c.PSet(pset, t.ToolChannel, ex_info, count)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil
}

func (p *Service) JobSet(sn string, job int, workorder_id int64, user_id int64) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	//workorder_id-user_id
	id_info := fmt.Sprintf("%s-%d-%d", controller.AUTO_MODE, workorder_id, user_id)

	err := c.JobSet(id_info, job)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil
}

func (p *Service) IDSet(sn string, str string) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	err := c.IdentifierSet(str)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil
}

func (p *Service) JobSetManual(sn string, tool_sn string, job int, user_id int64, ex_info string, skip bool) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	var t controller.ToolConfig

	var toolExist = false

	for _, t = range c.cfg.Tools {
		if t.SerialNO == tool_sn {
			toolExist = true
			break
		}
	}

	if !toolExist && !skip {
		return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	}

	return c.JobSet(ex_info, job)
}

func (p *Service) JobOFF(sn string, off bool) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	s_off := "0"
	if off {
		s_off = "1"
	}

	err := c.JobOff(s_off)
	if err != nil {
		// 控制器请求失败
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	return nil
}

func (p *Service) JobControl(sn string, action string) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	var err error
	if action == JOB_ACTION_ABORT {
		err = c.JobAbort()
	}

	if err != nil {
		// 控制器请求失败
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	return nil
}

func (p *Service) GetPSetList(sn string) ([]int, error) {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return []int{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	psetLists, err := c.GetPSetList()
	if err != nil {
		return []int{}, err
	}

	return psetLists, nil
}

func (p *Service) GetPSetDetail(sn string, pset int) (PSetDetail, error) {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return PSetDetail{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	psetDetail, err := c.GetPSetDetail(pset)
	if err != nil {
		return PSetDetail{}, err
	}

	return psetDetail, nil
}

func (p *Service) GetJobList(sn string) ([]int, error) {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return []int{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	jobLists, err := c.GetJobList()
	if err != nil {
		return []int{}, err
	}

	return jobLists, nil
}

func (p *Service) GetJobDetail(sn string, job int) (JobDetail, error) {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return JobDetail{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	jobDetail, err := c.GetJobDetail(job)
	if err != nil {
		return JobDetail{}, err
	}

	return jobDetail, nil
}

func (p *Service) IOSet(sn string, ios *[]IOStatus) error {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	return c.IOSet(ios)
}

func (p *Service) GenerateIDInfo(info string) string {
	ids := ""
	for i := 0; i < 4; i++ {
		if i == p.config().DataIndex {
			ids += fmt.Sprintf("%-25s", info)
		} else {
			ids += fmt.Sprintf("%25s", "")
		}
	}

	return ids
}

func (p *Service) TryCreateMaintenance(info ToolInfo) error {
	return p.Odoo.TryCreateMaintenance(info)
}
