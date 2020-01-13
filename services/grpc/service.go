package grpc

import (
	"fmt"
	"github.com/masami10/rush/services/transport"
	"github.com/masami10/rush/utils"
	microTransport "github.com/micro/go-micro/transport"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
)

type Service struct {
	microTransport.Transport
	diag        Diagnostic
	addr        string
	opened      bool
	mtxHandlers sync.Mutex
	msgHandlers map[string]transport.OnMsgHandler
	status      atomic.String
	client      transport.Client
}

func (s *Service) msgHandler(subject string) transport.OnMsgHandler {
	s.mtxHandlers.Lock()
	defer s.mtxHandlers.Unlock()
	if h, ok := s.msgHandlers[subject]; ok {
		return h
	}
	return nil
}

func newGRPCTransport(config Config) transport.Transport {
	options := microTransport.Addrs(config.addr)
	return microTransport.NewTransport(options)
}

func NewService(config Config, d Diagnostic) *Service {
	s := &Service{addr: config.addr, diag: d, msgHandlers: map[string]transport.OnMsgHandler{}}
	s.Transport = newGRPCTransport(config)
	s.status.Store(utils.STATUS_OFFLINE)
	return s
}

func (s *Service) Open() error {
	addr := s.addr
	if addr == "" {
		return errors.New("GRPC Address Is Empty")
	}
	c, err := s.doConnect(addr)
	if err != nil {
		return errors.Wrap(err, "Open")
	}
	s.client = c
	s.opened = true
	s.setStatus(utils.STATUS_ONLINE)
	return nil
}

func (s *Service) setStatus(status string) {
	s.status.Store(status)
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) doConnect(addr string) (transport.Client, error) {
	c, err := s.Dial(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s *Service) GetServerAddress() []string {
	return s.Transport.Options().Addrs
}

func (s *Service) OnMessage(subject string, handler transport.OnMsgHandler) error {
	s.mtxHandlers.Lock()
	defer s.mtxHandlers.Unlock()
	if _, ok := s.msgHandlers[subject]; ok {
		s.diag.Info(fmt.Sprintf("Handler For Subject: %s Is Existed, Do U Want To OverWrite!!!", subject))
	}
	s.msgHandlers[subject] = handler // overwrite
	return nil
}

func (s *Service) TransportForceOpen() error {
	if s.opened {
		return nil
	}
	return s.Open()
}

func (s *Service) Status() string {
	return s.status.Load()
}

func (s *Service) onStatus(status string) {
	s.setStatus(status)
}

func (s *Service) SendMessage(subject string, data []byte) error {
	if s == nil {
		return nil
	}

	return s.client.Send(&transport.Message{
		Header: map[string]string{transport.HEADER_SUBJECT: subject},
		Body:   data,
	})
}

func getMsgHeaderProperty(message *transport.Message, property string) string {
	header := message.Header
	pp := property
	switch property {
	case transport.HEADER_SUBJECT:
	case transport.HEADER_REPLY:
		pp = property
	default:
		return ""
	}
	if ret, ok := header[pp]; ok {
		return ret
	} else {
		return ""
	}
}

func getSubject(message *transport.Message) string {
	return getMsgHeaderProperty(message, transport.HEADER_SUBJECT)
}

func getReply(message *transport.Message) string {
	return getMsgHeaderProperty(message, transport.HEADER_REPLY)
}

func (s *Service) doOnMsg() {
	c := s.client
	if c == nil {
		return
	}
	var msg transport.Message
	for ; ; {
		if err := c.Recv(&msg); err != nil {
			subject := getSubject(&msg)
			if subject == "" {
				continue
			}
			h := s.msgHandler(subject)
			if h == nil {
				continue
			}
			if resp, err := h(&msg); err != nil {
				s.diag.Error(fmt.Sprintf("Handler Subject: %s Message Error", subject), err)
			} else {
				reply := getReply(&msg)
				if reply != "" {
					if err := s.SendMessage(reply, resp); err != nil {
						s.diag.Error(fmt.Sprintf("Reply Message To Subject: %s Error", reply), err)
					}
				}
			}
		}
	}
}
