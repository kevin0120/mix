package audi_vw

import (
	"github.com/masami10/rush/socket_listener"
	"sync"
	"github.com/masami10/rush/services/controller"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"bufio"
	"time"
	"log"
	"strings"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	StartManager()
}


type Service struct {

	configValue		atomic.Value

	name 			string
	listener		*socket_listener.SocketListener

	err 			chan error

	Controllers		map[string]Controller

	mux 			*sync.Mutex

	diag 			Diagnostic

}


func (s *Service) Err() <-chan error {
	return s.err
}


func NewService(c Config, d Diagnostic) *Service {
	addr := fmt.Sprintf("tcp://:%d",  c.Port)
	lis := socket_listener.NewSocketListener(addr)
	s := &Service{
		name: controller.AUDIPROTOCOL,
		listener: lis,
		mux: new(sync.Mutex),
		err: make(chan error ,1),
		diag: d,
	}

	s.configValue.Store(c)

	return s
}


func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}


func (p *Service) AddNewController(serial_no string)  *Controller{
	config := p.config()
	c := NewController(config)

	c.Srv = p //服务注入

	p.Controllers[serial_no] = c

	return &c
}

func (p *Service) Write(serial_no string ,buf []byte)  error{
	//p.mux.Lock()
	ws := p.Controllers
	//p.mux.Unlock()

	if _, ok := ws[serial_no]; !ok{
		return fmt.Errorf("Controller %s not found ", serial_no)
	}

	w := ws[serial_no]

	b, s := w.Write(data)

	w.Write(b,s)

	return nil
}

func (p *Service) Open() error {

	err := p.listener.Start()
	if err != nil {
		return errors.Wrapf(err, "Open Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Controllers{
		w.Connect()
	}


	return nil
}

func (p *Service) Close() error {
	err := p.listener.Close()
	if err != nil {
		return errors.Wrapf(err, "Close Protocol %s Listener fail", p.name)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	for _,w := range p.Controllers {
		err := w.Close()
		if err != nil {
			return errors.Wrapf(err, "Close Protocol %s Writer fail", p.name)
		}
	}

	return nil
}


// 此方法会被自动创建协程
func (p *Service) Read(c net.Conn) {
	ssl := p.listener
	defer ssl.InterListener.RemoveConnection(c)
	defer c.Close()

	buff := make([]byte, 4 * 4096) //创建一个打的buffer

	offset := 0

	scnr := bufio.NewScanner(c)
	for {
		if ssl.ReadTimeout.Nanoseconds() > 0 {
			c.SetReadDeadline(time.Now().Add(ssl.ReadTimeout))
		}
		if !scnr.Scan() {
			break
		}
		copy(buff[offset:], scnr.Bytes())
		buf, err := p.Parse(buff)
		if err != nil {
			log.Printf("unable to parse incoming line: %s", err)
			//TODO rate limit
			continue
		}
		if buf == nil {
			//猜测还需要继续读取数据
			offset = len(buf)
			continue
		}
		//解析成功，进行操作
		offset = 0 //偏移量回到0
	}

	if err := scnr.Err(); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("D! Timeout in plugin [input.socket_listener]: %s", err)
		} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
			log.Print(err)
		}
	}
}

func (p *Service) Parse(buf []byte)  ([]byte, error){

	return nil, nil
}


