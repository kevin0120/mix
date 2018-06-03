package storage

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/masami10/rush/services/httpd"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	OpenEngineSuccess(info string)
	Close()
	Closed()
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	eng         *xorm.Engine
	httpd       *httpd.Service
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

	exist, err := engine.IsTableExist("Workorders")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Workorders")
	}
	if !exist {
		if err := engine.Sync2(new(Workorders)); err != nil {
			return errors.Wrapf(err, "Create Table Workorders fail")
		}

	}

	exist, err = engine.IsTableExist("Results")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Results")
	}
	if !exist {
		if err := engine.Sync2(new(Results)); err != nil {
			return errors.Wrapf(err, "Create Table Results fail")
		}

	}

	exist, err = engine.IsTableExist("Curves")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Curves")
	}
	if !exist {
		if err := engine.Sync2(new(Curves)); err != nil {
			return errors.Wrapf(err, "Create Table Curves fail")
		}

	}

	s.eng = engine

	s.diag.OpenEngineSuccess(info)
	return nil
}

func (s *Service) Close() error {
	s.diag.Close()

	s.eng.Close()

	s.diag.Closed()

	return nil
}

func (s *Service) Store(data interface{}) error {
	session := s.eng.NewSession()
	defer session.Close()

	// add Begin() before any action
	err := session.Begin()
	_, err = session.Insert(data)
	if err != nil {
		session.Rollback()
		return errors.Wrapf(err, "store data fail")
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil
}
