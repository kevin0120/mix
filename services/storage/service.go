package storage

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/masami10/aiis/services/rush"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	OpenEngineSuccess(info string)
	UpdateResultSuccess(id int64)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	eng         *xorm.Engine
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Url,
		c.DBName)
	engine, err := xorm.NewEngine("postgres", info)

	if err != nil {
		return errors.Wrapf(err, "Create postgres engine fail")
	}

	exist, err := engine.IsTableExist("operation_result")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "operation_result")
	}
	if !exist {
		return errors.New("Check Table exist operation_result fail, Please start Odoo first")

	}

	engine.SetMaxOpenConns(c.MaxConnects) // always success

	s.eng = engine

	s.diag.OpenEngineSuccess(info)

	return nil
}

func (s *Service) Close() error {
	s.eng.Close()

	return nil
}

func (s *Service) UpdateResults(result *rush.OperationResult, id int64) error {

	affected, err := s.eng.Table("operation_result").Update(result)

	if err != nil {
		return errors.Wrapf(err, "Update result record %d fail", id)
	}

	s.diag.UpdateResultSuccess(affected)

	return nil
}
