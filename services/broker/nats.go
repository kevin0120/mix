package broker

import (
	"fmt"
	"github.com/masami10/rush/toml"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type Nats struct {
	sync.Mutex
	addrs      []string
	conn       *nats.Conn
	diag       Diagnostic
	opts       brokerOptions
	nopts      []nats.Option
	subscribes map[string]*nats.Subscription
}

func NewNats(d Diagnostic, addrs []string, options map[string]string) *Nats {
	s := &Nats{
		diag:       d,
		subscribes: map[string]*nats.Subscription{},
	}

	opts := generateConnectOptions(options)

	s.nopts = opts

	s.addrs = setAddrs(addrs)

	return s

}

func generateConnectOptions(options map[string]string) (opts []nats.Option)  {
	for key, value := range options {
		switch key {
		case "name":
			opts = append(opts, nats.Name(value))
		case "token":
			opts = append(opts, nats.Token(value))
		case "ping":
			var dur toml.Duration
			if err :=dur.UnmarshalText([]byte(value)); err == nil {
				opts = append(opts, nats.PingInterval(time.Duration(dur)))
			}
		}
	}
	return
}

func (s *Nats) Address() string {
	return strings.Join(s.addrs, ",")
}

func setAddrs(addrs []string) []string {
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}
	return cAddrs
}

func (s *Nats) Connect(urls []string) error {
	nc, err := nats.Connect(strings.Join(urls, ","), s.nopts...)
	if err != nil {
		s.diag.Error(fmt.Sprintf("connect To Urls: %s Error", urls), err)
		return err
	}
	s.conn = nc
	return nil
}

func (s *Nats) unsubscribeAll() error {
	var unSubs []string
	s.Lock()
	for subject, sub := range s.subscribes {
		if err := sub.Unsubscribe(); err != nil {
			s.diag.Error("Unsubscribe Error", err)
			continue
		}
		unSubs = append(unSubs, subject)
	}
	s.Unlock()

	for _, subject := range unSubs {
		delete(s.subscribes, subject)
	}
	return nil
}

func (s *Nats) Close() error {
	if s.conn != nil {
		s.unsubscribeAll()
		s.conn.Close()
	}
	return nil
}

func (s *Nats) addSub(subject string, sub *nats.Subscription) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.subscribes[subject]; !ok {
		s.subscribes[subject] = sub
	}
}

func (s *Nats) removeSub(subject string) error {
	s.Lock()
	defer s.Unlock()
	if ss, ok := s.subscribes[subject]; ok {
		return ss.Unsubscribe()
	}
	delete(s.subscribes, subject)
	return nil
}

func (s *Nats) Subscribe(subject string, handler SubscribeHandler) error {
	if s.conn == nil {
		err := errors.New("Nats Is Not Connected!")
		s.diag.Error("Subscribe Error", err)
		return err
	}
	nc := s.conn
	fn := func(msg *nats.Msg) {
		d := &brokerMessage{
			Body:   msg.Data,
			Header: map[string]string{"subject": msg.Subject, "Reply": msg.Reply},
		}
		handler(d)
	}
	ss, err := nc.Subscribe(subject, fn)
	if err != nil {
		return err
	}
	s.addSub(subject, ss)
	return nil
}

func (s *Nats) UnSubscribe(subject string) error {
	return s.removeSub(subject)

}

func (s *Nats) Publish(subject string, data[]byte) error {
	if s.conn == nil {
		err := errors.New("Nats Is Not Connected!")
		s.diag.Error("Subscribe Error", err)
		return err
	}
	nc := s.conn
	return nc.Publish(subject, data)
}
