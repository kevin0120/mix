package broker

import (
	"strings"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn *nats.Client
	diag Diagnostic
}

func NewNats(d Diagnostic) *Nats  {
	return &Nats{
		diag: d,
	}
}

func (s *Nats)Connect(urls []string) error  {
	nc, err := nats.Connect(strings.Join(urls, ","))
	if err != nil {
		return err
	}
	s.conn = nc
	return nil
}
