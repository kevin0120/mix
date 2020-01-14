package grpc

import (
	"context"
	"fmt"
	"github.com/masami10/rush/services/transport"
	"github.com/masami10/rush/utils"
	microTransport "github.com/micro/go-micro/transport"
	_ "github.com/micro/go-micro/transport/grpc"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type Service struct {
	microTransport.Transport
	diag        Diagnostic
	addr        string
	workers     int
	opened      bool
	mtxHandlers sync.Mutex
	msgHandlers map[string]transport.OnMsgHandler
	status      atomic.String
	msgs        chan transport.Message
	exiting     chan chan struct{}
	respChannel chan RespStruct
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
	options := microTransport.Addrs(config.Address)
	return microTransport.NewTransport(options)
}

func NewService(config Config, d Diagnostic) *Service {
	s := &Service{
		addr:        config.Address,
		workers:     config.Workers,
		exiting:     make(chan chan struct{}, config.Workers),
		diag:        d,
		msgs:        make(chan transport.Message, 1024),
		respChannel: make(chan RespStruct, 1024),
		msgHandlers: map[string]transport.OnMsgHandler{},
	}
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
	for i := 0; i < s.workers; i++ {
		go s.doWork()
	}
	return nil
}

func (s *Service) setStatus(status string) {
	s.status.Store(status)
}

func (s *Service) Close() error {
	if !s.opened {
		return nil
	}
	for i := 0; i < s.workers; i++ {
		exit := make(chan struct{})
		s.exiting <- exit
		<-exit
	}
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
	msg := &transport.Message{
		Header: map[string]string{transport.HEADER_SUBJECT: subject},
		Body:   data,
	}

	return s.doSendMessage(msg)
}

func (s *Service) doSendMessage(msg *transport.Message) error {
	if s.client == nil {
		return errors.New("GRPC Client Is Empty")
	}
	return s.client.Send(msg)
}

func getMsgHeaderProperty(message *transport.Message, property string) string {
	header := message.Header
	pp := property
	switch property {
	case transport.HEADER_SUBJECT:
	case transport.HEADER_REPLY:
	case transport.HEADER_MSG_ID:
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

func getMessageID(message *transport.Message) string {
	return getMsgHeaderProperty(message, transport.HEADER_MSG_ID)
}

func (s *Service) Request(subject string, data []byte, timeOut time.Duration) ([]byte, error) {
	msgId := uuid.NewV4().String()
	msg := &transport.Message{
		Header: map[string]string{transport.HEADER_SUBJECT: subject, transport.HEADER_MSG_ID: msgId},
		Body:   data,
	}
	if err := s.doSendMessage(msg); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	respChannel := s.respChannel
	cc := make(chan *transport.Message, 2) //确保通道不会阻塞
	//defer close(cc)
	respChannel <- RespStruct{msgId, cc}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-cc:
		if msg == nil {
			return nil, errors.New(fmt.Sprintf("Opening, Channel Is %v", ok))
		}
		return msg.Body, nil
	}
}

func (s *Service) doWork() {
	var needRespMsgIds map[string]chan *transport.Message
	for ; ; {
		select {
		case exit := <-s.exiting:
			close(exit)
			return
		case resp := <-s.respChannel:
			msgId := resp.msgId
			needRespMsgIds[msgId] = resp.msg
		case msg := <-s.msgs:
			subject := getSubject(&msg)
			if subject == "" {
				continue
			}
			msgId := getMessageID(&msg)
			if e, ok := needRespMsgIds[msgId]; ok {
				if e != nil {
					e <- &msg
					close(e)
				}
				delete(needRespMsgIds, msgId)
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

func (s *Service) doOnMsg() {
	c := s.client
	if c == nil {
		return
	}
	var msg transport.Message
	for ; ; {
		if err := c.Recv(&msg); err != nil {
			s.msgs <- msg
		}
	}
}
