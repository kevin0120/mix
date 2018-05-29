package diagnostic

import (
	dia "github.com/influxdata/kapacitor/services/diagnostic"
	"io"
	"os"
)

type Service struct {
	*dia.Service
}

func NewService(c Config, stdout, stderr io.Writer) *Service {
	diac := dia.Config{c.File, c.Level}
	s := dia.NewService(diac, stdout, stderr)
	return &Service{ s}
}

func (s *Service) NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		l: s.Logger.With( dia.String("service", "run")),
	}
}

func BootstrapMainHandler() *CmdHandler {
	s := NewService(NewConfig(), nil, os.Stderr)
	// Should never error
	_ = s.Open()

	return s.NewCmdHandler()
}
