package openprotocol

import (
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
	Controllers map[string]*Controller
	DB          *storage.Service
	WS          *wsnotify.Service
	Aiis        *aiis.Service
	Minio       *minio.Service
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		Controllers: map[string]*Controller{},
		name:        controller.OPENPROTOCOL,
		diag:        d,
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

func (p *Service) AddNewController(cfg controller.Config) {
	config := p.config()
	c := NewController(config)
	c.Srv = p //服务注入
	c.cfg = cfg
	p.Controllers[cfg.SN] = &c
}

func (p *Service) Open() error {
	for _, w := range p.Controllers {
		go w.Start() //异步启动控制器
	}

	return nil
}

func (p *Service) Close() error {
	for _, w := range p.Controllers {
		err := w.Close()
		if err != nil {
			return errors.Wrapf(err, "Close Protocol %s Writer fail", p.name)
		}
	}

	return nil
}
