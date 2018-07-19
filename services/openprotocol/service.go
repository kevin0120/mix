package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
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

	Parent *controller.Service
}

func NewService(c Config, d Diagnostic, parent *controller.Service) *Service {

	s := &Service{
		name:   controller.OPENPROTOCOL,
		diag:   d,
		Parent: parent,
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

func (p *Service) AddNewController(cfg controller.Config) controller.Controller {
	config := p.config()
	c := NewController(config)
	c.Srv = p //服务注入
	c.cfg = cfg

	return &c
}

func (p *Service) Open() error {
	for _, w := range p.Parent.Controllers {
		if w.Protocol() == controller.OPENPROTOCOL {
			go w.Start() //异步启动控制器
		}
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

func (p *Service) PSet(sn string, pset int, workorder_id int64, result_id int64, count int, user_id int64) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	if v.Status() == controller.STATUS_OFFLINE {
		// 控制器离线
		return errors.New(string(controller.STATUS_OFFLINE))
	}

	c := v.(*Controller)
	// 设定pset并判断控制器响应
	_, err := c.PSet(pset, workorder_id, result_id, count, user_id, c.cfg.ToolChannel)
	if err != nil {
		// 控制器请求失败
		return errors.New(controller.ERR_PSET_ERROR)
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
	id_info := fmt.Sprintf("%d-%d", workorder_id, user_id)

	err := c.JobSet(id_info, job)
	if err != nil {
		// 控制器请求失败
		return errors.New(controller.ERR_PSET_ERROR)
	}

	return nil
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
		return errors.New(controller.ERR_PSET_ERROR)
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

	pset_list, err := c.GetPSetList()
	if err != nil {
		return []int{}, err
	}

	return pset_list, nil
}

func (p *Service) GetPSetDetail(sn string, pset int) (PSetDetail, error) {
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return PSetDetail{}, errors.New(controller.ERR_CONTROLER_NOT_FOUND)
	}

	c := v.(*Controller)

	pset_detail, err := c.GetPSetDetail(pset)
	if err != nil {
		return PSetDetail{}, err
	}

	return pset_detail, nil
}