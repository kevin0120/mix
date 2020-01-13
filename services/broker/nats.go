package broker

import (
	"fmt"
	"github.com/masami10/rush/services/transport"
	"github.com/masami10/rush/toml"
	"github.com/masami10/rush/utils"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
	"sync"
	"time"
)

type Nats struct {
	sync.Mutex
	addrs         []string
	conn          *nats.Conn
	diag          Diagnostic
	opts          transport.Options
	nopts         []nats.Option
	subscribes    map[string]*nats.Subscription
	loadBalancers map[string][]string
	respSubject   string
	workGroups    map[string][]string
	handler       StatusHandler
}

func NewNats(d Diagnostic, c Config) *Nats {
	addrs := c.ConnectUrls
	options := c.Options
	s := &Nats{
		diag:          d,
		subscribes:    map[string]*nats.Subscription{},
		loadBalancers: map[string][]string{},
		workGroups:    map[string][]string{},
		respSubject:   fmt.Sprintf("%s.responses", c.Name),
	}

	opts := generateConnectOptions(options)

	s.nopts = opts

	s.addrs = setAddrs(addrs)

	return s

}

func generateConnectOptions(options map[string]string) (opts []nats.Option) {
	for key, value := range options {
		switch key {
		case "name":
			opts = append(opts, nats.Name(value))
		case "token":
			opts = append(opts, nats.Token(value))
		case "ping":
			var dur toml.Duration
			if err := dur.UnmarshalText([]byte(value)); err == nil {
				opts = append(opts, nats.PingInterval(time.Duration(dur)))
			}
		}
	}
	return
}

func (s *Nats) Address() string {
	return strings.Join(s.addrs, ",")
}

func (s *Nats) SetStatusHandler(handler StatusHandler) {
	s.handler = handler
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

func (s *Nats) handleStatus(status nats.Status) {
	switch status {
	case nats.CONNECTED:
		s.handler(utils.STATUS_ONLINE)
	case nats.DISCONNECTED:
		s.handler(utils.STATUS_OFFLINE)
	}
}

func (s *Nats) statusHandler(conn *nats.Conn) {
	s.handleStatus(conn.Status())
	if cid, err := conn.GetClientID(); err == nil {
		s.diag.Debug(fmt.Sprintf("Client %d is %s ", cid, STATUS_BROKER[conn.Status()]))
	} else {
		s.diag.Error("statusHandler", err)
	}
}

func (s *Nats) statusErrHandler(conn *nats.Conn, err error) {
	s.handleStatus(conn.Status())
	if cid, err := conn.GetClientID(); err == nil {
		s.diag.Error(fmt.Sprintf("Client %d is %s ", cid, STATUS_BROKER[conn.Status()]), err)
	} else {
		s.diag.Error("statusErrHandler", err)
	}
}

func (s *Nats) setDefaultHandlers() {
	nc := s.conn
	nc.SetReconnectHandler(s.statusHandler)
	nc.SetDisconnectErrHandler(s.statusErrHandler)
	nc.SetClosedHandler(s.statusHandler)
}

func (s *Nats) Connect(urls []string) error {
	nc, err := nats.Connect(strings.Join(urls, ","), s.nopts...)
	if err != nil {
		s.diag.Error(fmt.Sprintf("connect To Urls: %s Error", urls), err)
		return err
	}
	s.conn = nc

	//最后设置默认句柄
	s.setDefaultHandlers()

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
		_ = s.unsubscribeAll()
		s.conn.Close()
	}
	return nil
}

func (s *Nats) isExist(subject string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.subscribes[subject]
	return ok
}

func (s *Nats) addSub(subject string, sub *nats.Subscription) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.subscribes[subject]; !ok {
		s.subscribes[subject] = sub
	} else {
		//已经存在
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

func (s *Nats) AppendWorkGroup(subject string, group string, handler transport.OnMsgHandler) error {
	nc := s.conn
	if nc == nil {
		err := errors.New("Nats Is Not Connected!")
		s.diag.Error("AppendWorkGroup Error", err)
		return err
	}
	if group == "" {
		err := errors.New("Group Name Is Required")
		s.diag.Error("AppendWorkGroup Error", err)
		return err
	}
	r := s.ensureWorkGroup(subject, group)

	if id, err := s.subscribe(subject, group, handler); err != nil {
		return err
	} else {
		r = append(r, id)

	}

	return nil
}

func (s *Nats) ensureWorkGroup(subject string, group string) []string {
	name := fmt.Sprintf("%s@%s", subject, group)
	if _, ok := s.loadBalancers[name]; !ok {
		s.loadBalancers[name] = []string{}
	}

	return s.loadBalancers[name]
}

func (s *Nats) subscribe(subject, group string, handler transport.OnMsgHandler) (registerName string, err error) {
	nc := s.conn
	if nc == nil {
		err = errors.New("Nats Is Not Connected!")
		s.diag.Error("Subscribe Error", err)
		return
	}
	if s.isExist(subject) {
		err = errors.Errorf("Subject: %s Is Been Subscribed, Please Unsubscribe First", subject)
		s.diag.Error("Subscribe Error", err)
		return
	}

	fn := func(msg *nats.Msg) {
		d := &transport.Message{
			Body:   msg.Data,
			Header: map[string]string{transport.HEADER_SUBJECT: msg.Subject, transport.HEADER_REPLY: msg.Reply},
		}
		if resp, err := handler(d); err != nil {
			s.diag.Error("Subscribe Handler Error", err)
		} else {
			if len(resp) == 0 {
				return
			}
			if err := nc.Publish(msg.Reply, resp); err != nil {
				s.diag.Error("Subscribe Handler Publish Error", err)
			}
		}
	}

	var ss *nats.Subscription
	registerName = subject
	if group == "" {
		ss, err = nc.Subscribe(subject, fn)

	} else {
		ss, err = nc.QueueSubscribe(subject, group, fn)
		registerName = uuid.NewV4().String()
	}

	if err != nil {
		return
	}
	s.addSub(registerName, ss)
	return
}

// 创建一个新的订阅，然后对其订阅发起返回
func (s *Nats) Subscribe(subject string, handler transport.OnMsgHandler) error {
	_, err := s.subscribe(subject, "", handler)
	return err
}

func (s *Nats) UnSubscribe(subject string) error {
	return s.removeSub(subject)
}

func (s *Nats) Publish(subject string, data []byte) error {
	nc := s.conn
	if nc == nil {
		err := errors.New("Nats Is Not Connected!")
		s.diag.Error("Publish Error", err)
		return err
	}

	return nc.Publish(subject, data)
}

func (s *Nats) DoRequest(subject string, data []byte, timeOut time.Duration) (resp []byte, err error) {
	nc := s.conn
	if nc == nil {
		err = errors.New("Nats Is Not Connected!")
		s.diag.Error("DoRequest Error", err)
		return
	}
	if timeOut.Nanoseconds() == 0 {
		timeOut = time.Second * 10
	}

	if r, rErr := nc.Request(subject, data, timeOut); rErr != nil {
		s.diag.Error("Nats Request Error", err)
		err = rErr
		return
	} else {
		s.diag.Debug(fmt.Sprintf("Get Resp With Msg: %v", resp))
		resp = r.Data
	}

	return
}
