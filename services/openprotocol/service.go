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

	devices []*TighteningController
	tightening_device.ITighteningProtocol
}

func NewService(c Config, d Diagnostic, parent *controller.Service) *Service {

	s := &Service{
		name:    controller.OPENPROTOCOL,
		diag:    d,
		Parent:  parent,
		devices: []*TighteningController{},
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Name() string {
	return "OpenProtocol"
}

func (s *Service) Parse(msg string) ([]byte, error) {
	return nil, nil
}

func (s *Service) Write(sn string, buf []byte) error {
	return nil
}

func (s *Service) IsSupport(cfg *tightening_device.TighteningDeviceConfig) bool {
	_, err := GetModel(cfg.Model)
	if err != nil {
		return false
	}

	return true
}

func (s *Service) CreateController(cfg *tightening_device.TighteningDeviceConfig) (tightening_device.ITighteningController, error) {
	// 检测型号是否支持
	if !s.IsSupport(cfg) {
		return nil, errors.New("Not Supported")
	}

	protocolConfig := s.config()
	c := NewController(&protocolConfig, cfg, s.diag)
	c.Srv = s

	return c, nil
}

func (s *Service) Open() error {

	return nil
}

func (s *Service) Close() error {

	return nil
}

// TODO: 对外接口统一在tightening服务中实现
func (s *Service) ToolControl(sn string, tool_sn string, enable bool) error {
	//// 判断控制器是否存在
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	//var toolExist = false
	//
	//for _, t := range c.cfg.Tools {
	//	if t.SN == tool_sn {
	//		toolExist = true
	//		break
	//	}
	//}
	//
	//if !toolExist {
	//	return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	//}
	//
	//// 工具使能
	//err := c.ToolControl(enable)
	//if err != nil {
	//	// 控制器请求失败
	//	return err
	//}

	return nil
}

func (s *Service) PSet(sn string, tool_sn string, pset int, result_id int64, count int, user_id int64) error {
	// 判断控制器是否存在
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	//var toolExist = false
	//var toolChannel int
	//for _, t := range c.cfg.Tools {
	//	if t.SN == tool_sn {
	//		toolChannel = t.Channel
	//		toolExist = true
	//		break
	//	}
	//}
	//
	//if !toolExist {
	//	return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	//}
	//
	//ex_info := fmt.Sprintf("%s-%d-%d-%d", controller.AUTO_MODE, result_id, count, user_id)
	//
	//// 设定pset并判断控制器响应
	//_, err := c.PSet(pset, toolChannel, ex_info, 1)
	//if err != nil {
	//	// 控制器请求失败
	//	return err
	//}

	return nil
}

func (s *Service) PSetManual(sn string, tool_sn string, pset int, user_id int64, ex_info string, count int) error {
	// 判断控制器是否存在
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	//var toolExist = false
	//var toolChannel int
	//for _, t := range c.cfg.Tools {
	//	if t.SN == tool_sn {
	//		toolChannel = t.Channel
	//		toolExist = true
	//		break
	//	}
	//}
	//
	//if !toolExist {
	//	return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	//}
	//
	//// 设定pset并判断控制器响应
	//_, err := c.PSet(pset, toolChannel, ex_info, count)
	//if err != nil {
	//	// 控制器请求失败
	//	return err
	//}

	return nil
}

func (s *Service) JobSet(sn string, job int, workorder_id int64, user_id int64) error {
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	////workorder_id-user_id
	//id_info := fmt.Sprintf("%s-%d-%d", controller.AUTO_MODE, workorder_id, user_id)
	//
	//err := c.JobSet(id_info, job)
	//if err != nil {
	//	// 控制器请求失败
	//	return err
	//}

	return nil
}

func (s *Service) IDSet(sn string, str string) error {
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	//err := c.IdentifierSet(str)
	//if err != nil {
	//	// 控制器请求失败
	//	return err
	//}

	return nil
}

func (s *Service) JobSetManual(sn string, tool_sn string, job int, user_id int64, ex_info string, skip bool) error {
	//v, exist := s.Parent.Controllers[sn]
	//if !exist {
	//	// SN对应控制器不存在
	//	return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//}
	//
	//c := v.(*TighteningController)
	//
	//var toolExist = false
	//
	//for _, t := range c.cfg.Tools {
	//	if t.SN == tool_sn {
	//		toolExist = true
	//		break
	//	}
	//}
	//
	//if !toolExist && !skip {
	//	return errors.New(fmt.Sprintf(controller.ERR_TOOL_NOT_FOUND+" tool serial number:%s", tool_sn))
	//}
	//
	//return c.JobSet(ex_info, job)
	return nil
}

func (s *Service) JobOFF(sn string, off bool) error {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	s_off := "0"
	//	if off {
	//		s_off = "1"
	//	}
	//
	//	err := c.JobOff(s_off)
	//	if err != nil {
	//		// 控制器请求失败
	//		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	//	}
	//
	return nil
}

func (s *Service) JobControl(sn string, action string) error {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	var err error
	//	if action == JOB_ACTION_ABORT {
	//		err = c.JobAbort()
	//	}
	//
	//	if err != nil {
	//		// 控制器请求失败
	//		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	//	}
	//
	return nil
}

func (s *Service) GetPSetList(sn string) ([]int, error) {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return []int{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	psetLists, err := c.GetPSetList()
	//	if err != nil {
	//		return []int{}, err
	//	}
	//
	//	return psetLists, nil
	return []int{}, nil
}

//
func (s *Service) GetPSetDetail(sn string, pset int) (interface{}, error) {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return PSetDetail{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	psetDetail, err := c.GetPSetDetail(pset)
	//	if err != nil {
	//		return PSetDetail{}, err
	//	}
	//
	return nil, nil
}

//
func (s *Service) GetJobList(sn string) ([]int, error) {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return []int{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	jobLists, err := c.GetJobList()
	//	if err != nil {
	//		return []int{}, err
	//	}
	//
	return []int{}, nil
}

//
func (s *Service) GetJobDetail(sn string, job int) (interface{}, error) {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return JobDetail{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	jobDetail, err := c.GetJobDetail(job)
	//	if err != nil {
	//		return JobDetail{}, err
	//	}
	//
	return nil, nil
}

func (s *Service) IOSet(sn string, ios *[]IOStatus) error {
	//	v, exist := s.Parent.Controllers[sn]
	//	if !exist {
	//		// SN对应控制器不存在
	//		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	//	}
	//
	//	c := v.(*TighteningController)
	//
	//	return c.IOSet(ios)
	return nil
}

func (s *Service) generateIDInfo(info string) string {
	ids := ""
	for i := 0; i < MAX_IDS_NUM; i++ {
		if i == s.config().DataIndex {
			ids += fmt.Sprintf("%-25s", info)
		} else {
			ids += fmt.Sprintf("%25s", "")
		}
	}

	return ids
}

func (s *Service) TryCreateMaintenance(info ToolInfo) error {
	return s.Odoo.TryCreateMaintenance(info)
}
