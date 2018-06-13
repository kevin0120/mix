package audi_vw

import (
	"github.com/masami10/rush/socket_listener"
	"sync"
	"github.com/masami10/rush/services/controller"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync/atomic"
	"time"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"strings"
)

const (
	ERR_CVI3_NOT_FOUND = "CIV3 SN is invalid"
	ERR_CVI3_OFFLINE = "cvi3 offline"
	ERR_CVI3_REQUEST = "request to cvi3 failed"
	ERR_CVI3_REPLY_TIMEOUT = "cvi3 reply timeout"
	ERR_CVI3_REPLY = "cvi3 reply contains error"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
	StartManager()
}

type ControllerStatus struct {
	SN string `json:"controller_sn"`
	Status ControllerStatusType `json:"status"`
}


type Service struct {

	configValue		atomic.Value
	name 			string
	listener		*socket_listener.SocketListener
	err 			chan error
	Controllers		map[string]*Controller
	mux 			*sync.Mutex
	diag 			Diagnostic
	DB				*storage.Service
	WS				*wsnotify.Service
	Aiis			*aiis.Service
	Minio			*minio.Service
	handlers		Handlers
	handle_buffer	chan string
}


func (s *Service) Err() <-chan error {
	return s.err
}

func NewService(c Config, d Diagnostic) *Service {
	addr := fmt.Sprintf("tcp://:%d",  c.Port)

	s := &Service{
		Controllers: map[string]*Controller{},
		name: controller.AUDIPROTOCOL,
		mux: new(sync.Mutex),
		err: make(chan error ,1),
		diag: d,
		handlers: Handlers{},
	}

	s.handle_buffer = make(chan string, 1024)
	s.handlers.AudiVw = s
	lis := socket_listener.NewSocketListener(addr, s)
	s.listener = lis
	s.configValue.Store(c)

	return s
}


func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (p *Service) AddNewController(cfg controller.Config) {
	config := p.config()
	c := NewController(config)
	c.Srv = p //服务注入
	c.cfg = cfg
	p.Controllers[cfg.SN] = &c
}

func (p *Service) Write(serial_no string ,buf []byte)  error{
	return nil
}

func (p *Service) Open() error {

	p.diag.StartManager()

	err := p.listener.Start()
	if err != nil {
		return errors.Wrapf(err, "Open Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Controllers{
		go w.Start()
	}

	//cpu_num := runtime.NumCPU()
	for i := 0; i < p.config().Workers; i++ {
		go p.HandleProcess()
		p.diag.Debug(fmt.Sprintf("init handle process:%d", i + 1))
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

func (p *Service) NewConn(c net.Conn) {}

// 服务端读取
func (p *Service) Read(c net.Conn) {
	ssl := p.listener
	defer ssl.InterListener.RemoveConnection(c)
	defer c.Close()

	buffer := make([]byte, p.config().ReadBufferSize)
	for {

		n, err := c.Read(buffer)
		if err != nil {
			p.diag.Error("read err", err)
			break
		}

		msg := string(buffer[0:n])
		if len(msg) < HEADER_LEN {
			continue
		}

		header := CVI3Header{}
		header.Deserialize(msg[0: HEADER_LEN])
		//fmt.Printf("%d\n", header.SIZ)
		var body string = msg[HEADER_LEN: n]
		var rest int = header.SIZ - HEADER_LEN - n
		for {
			if rest <= 0 {
				break
			}
			n, err := c.Read(buffer)
			if err != nil {
				break
			}
			body += string(buffer[0:n])
			rest -= n
		}

		p.Parse([]byte(body))

		if header.TYP == Header_type_request_with_reply || header.TYP == Header_type_keep_alive {
			// 执行应答
			reply := CVI3Header{}
			reply.Init()
			reply.TYP = Header_type_reply
			reply.MID = header.MID
			reply_packet := reply.Serialize()

			_, err := c.Write([]byte(reply_packet))
			if err != nil {
				print("server reply err:%s\n", err.Error())
				break
			}
		}
	}
}

func (p *Service) Parse(buf []byte)  ([]byte, error){

	msg := string(buf)

	if strings.Contains(msg, XML_RESULT_KEY) {
		p.handle_buffer <- msg
	}

	return nil, nil
}

func (p *Service) HandleProcess() {
	var context = HandlerContext {
		cvi3_result: CVI3Result{},
		controller_curve: ControllerCurve{},
		controller_result: ControllerResult{},
		controller_curve_file: ControllerCurveFile{},
		db_curve: storage.Curves{},
		ws_result: wsnotify.WSResult{},
		aiis_result: aiis.AIISResult{},
		aiis_curve: aiis.CURObject{},
	}

	for {
		msg := <- p.handle_buffer
		p.handlers.HandleMsg(msg, &context)
	}
}

// 取得控制器状态
func (p *Service) GetControllersStatus(sn string) ([]ControllerStatus, error) {
	status := []ControllerStatus{}
	if sn != "" {
		c, exist := p.Controllers[sn]
		if !exist {
			return status, errors.New("controller not found")
		} else {
			s := ControllerStatus{}
			s.SN = sn
			s.Status = c.GetStatus()
			status = append(status, s)
			return status, nil
		}
	} else {
		for k, v := range p.Controllers {
			s := ControllerStatus{}
			s.SN = k
			s.Status = v.GetStatus()
			status = append(status, s)
		}

		return status, nil
	}
}

// 设置拧接程序
func (p *Service) PSet(sn string, pset int, workorder_id int64, result_id int64, count int, user_id int64) (error) {
	// 判断控制器是否存在
	c, exist := p.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(ERR_CVI3_NOT_FOUND)
	}

	if c.GetStatus() == STATUS_OFFLINE {
		// 控制器离线
		return errors.New(ERR_CVI3_OFFLINE)
	}

	// 设定pset并判断控制器响应
	serial, err := c.PSet(pset, workorder_id, result_id, count, user_id)
	if err != nil {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REQUEST)
	}

	var header_str string
	for i := 0; i < 6; i++ {
		header_str = c.Response.get(serial)
		if header_str != "" {
			c.Response.remove(serial)
			break
		}
		time.Sleep(time.Duration(c.req_timeout))
	}

	if header_str == "" {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := CVI3Header{}
	header.Deserialize(header_str)
	if !header.Check() {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY)
	}

	return nil
}
