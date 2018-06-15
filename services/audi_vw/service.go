package audi_vw

import (
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/socket_listener"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	//"time"
	"time"
)

const (
	ERR_CVI3_NOT_FOUND     = "CIV3 SN is invalid"
	ERR_CVI3_OFFLINE       = "cvi3 offline"
	ERR_CVI3_REQUEST       = "request to cvi3 failed"
	ERR_CVI3_REPLY_TIMEOUT = "cvi3 reply timeout"
	ERR_CVI3_REPLY         = "cvi3 reply contains error"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
	StartManager()
}

type ControllerStatus struct {
	SN     string               `json:"controller_sn"`
	Status ControllerStatusType `json:"status"`
}

type Service struct {
	configValue   atomic.Value
	name          string
	listener      *socket_listener.SocketListener
	err           chan error
	Controllers   map[string]*Controller
	diag          Diagnostic
	DB            *storage.Service
	WS            *wsnotify.Service
	Aiis          *aiis.Service
	Minio         *minio.Service
	handlers      Handlers
	wg            sync.WaitGroup
	closing       chan struct{}
	handle_buffer chan string
}

func (s *Service) Err() <-chan error {
	return s.err
}

func NewService(c Config, d Diagnostic) *Service {
	addr := fmt.Sprintf("tcp://:%d", c.Port)

	s := &Service{
		Controllers: map[string]*Controller{},
		name:        controller.AUDIPROTOCOL,
		err:         make(chan error, 1),
		diag:        d,
		wg:          sync.WaitGroup{},
		closing:     make(chan struct{}, 1),
		handlers:    Handlers{},
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

func (p *Service) Write(sn string, buf []byte) error {
	if _, ok := p.Controllers[sn]; !ok {
		return fmt.Errorf("can not found controller :%s", sn)
	}
	c := p.Controllers[sn]
	s := c.Sequence()
	c.Write(buf, s)
	return nil
}

func (p *Service) Open() error {

	p.diag.StartManager()

	err := p.listener.Start()
	if err != nil {
		return errors.Wrapf(err, "Open Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Controllers {
		go w.Start() //异步启动控制器
	}

	for i := 0; i < p.config().Workers; i++ {
		p.wg.Add(1)
		go p.HandleProcess()
		p.diag.Debug(fmt.Sprintf("init handle process:%d", i+1))
	}

	return nil
}

func (p *Service) Close() error {
	err := p.listener.Close()
	if err != nil {
		return errors.Wrapf(err, "Close Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Controllers {
		err := w.Close()
		if err != nil {
			return errors.Wrapf(err, "Close Protocol %s Writer fail", p.name)
		}
	}
	for i := 0; i < p.config().Workers; i++ {
		p.closing <- struct{}{}
		p.diag.Debug(fmt.Sprintf("Close AUDIVW Server Handler:%d", i+1))
	}

	p.wg.Wait() //阻塞 等待全部handler关闭

	return nil
}

func (p *Service) NewConn(c net.Conn) {}

// 服务端读取
func (p *Service) Read(c net.Conn) {
	ssl := p.listener
	defer ssl.InterListener.RemoveConnection(c)
	defer c.Close()

	rest := 0
	body := ""
	header := CVI3Header{}
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
		off := 0 //循环前偏移为0
		for off < n {
			if rest == 0 {
				header.Deserialize(msg[off : off+HEADER_LEN])
				if n-off > HEADER_LEN+header.SIZ {
					//粘包
					body = msg[off+HEADER_LEN : off+HEADER_LEN+header.SIZ]
					p.Parse([]byte(body))
					p.CVIResponse(&header, c)
					off += HEADER_LEN + header.SIZ
					rest = 0 //同样解析头
				} else {
					body = msg[off+HEADER_LEN : n]
					rest = header.SIZ - (n - (off + HEADER_LEN))
					break
				}
			} else {
				if n-off > rest {
					//粘包
					body += string(buffer[off : off+rest]) //已经是完整的包
					p.Parse([]byte(body))
					p.CVIResponse(&header, c)
					off += rest
					rest = 0 //进入解析头
				} else {
					body += string(buffer[off:n])
					rest -= n - off
					break
				}
			}
		}

		if rest == 0 {
			p.Parse([]byte(body))
			p.CVIResponse(&header, c)
		}
	}
}

func (p *Service) CVIResponse(header *CVI3Header, c net.Conn) {
	if header.TYP == Header_type_request_with_reply || header.TYP == Header_type_keep_alive {
		// 执行应答
		reply := CVI3Header{}
		reply.Init()
		reply.TYP = Header_type_reply
		reply.MID = header.MID
		replyPacket := reply.Serialize()

		_, err := c.Write([]byte(replyPacket))
		if err != nil {
			print("server reply err:%s\n", err.Error())
		}
	}
}

func (p *Service) Parse(buf []byte) ([]byte, error) {

	msg := string(buf)

	if strings.Contains(msg, XML_RESULT_KEY) {
		p.handle_buffer <- msg
	}

	return nil, nil
}

func (p *Service) HandleProcess() {
	var context = HandlerContext{
		cvi3Result:          CVI3Result{},
		controllerCurve:     ControllerCurve{},
		controllerResult:    ControllerResult{},
		controllerCurveFile: ControllerCurveFile{},
		dbCurve:             storage.Curves{},
		wsResult:            wsnotify.WSResult{},
		aiisResult:          aiis.AIISResult{},
		aiisCurve:           aiis.CURObject{},
	}
	for {
		select {
		case msg := <-p.handle_buffer:
			p.handlers.HandleMsg(msg, &context)

		case <-p.closing:
			p.wg.Done()
			return
		}
	}

}

// 取得控制器状态
func (p *Service) GetControllersStatus(sn string) ([]ControllerStatus, error) {
	var status []ControllerStatus
	if sn != "" {
		c, exist := p.Controllers[sn]
		if !exist {
			return status, errors.New("controller not found")
		} else {
			s := ControllerStatus{}
			s.SN = sn
			s.Status = c.Status()
			status = append(status, s)
			return status, nil
		}
	} else {
		for k, v := range p.Controllers {
			s := ControllerStatus{}
			s.SN = k
			s.Status = v.Status()
			status = append(status, s)
		}

		return status, nil
	}
}

// 设置拧接程序
func (p *Service) PSet(sn string, pset int, workorder_id int64, result_id int64, count int, user_id int64) error {
	// 判断控制器是否存在
	c, exist := p.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(ERR_CVI3_NOT_FOUND)
	}

	if c.Status() == STATUS_OFFLINE {
		// 控制器离线
		return errors.New(ERR_CVI3_OFFLINE)
	}

	// 设定pset并判断控制器响应
	seq, err := c.PSet(pset, workorder_id, result_id, count, user_id)
	if err != nil {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REQUEST)
	}

	c.Response.Add(seq, "")

	defer c.Response.remove(seq)

	//i := 0

	//for {
	//	select {
	//	//case <- time.After(time.Duration(c.req_timeout)):
	//	//	i += 1
	//	//	if i >= 6 {
	//	//		return errors.New(ERR_CVI3_REPLY_TIMEOUT)
	//	//	}
	//	case headerStr := <-c.response:
	//		header := CVI3Header{}
	//		header.Deserialize(headerStr)
	//		if !header.Check() {
	//			// 控制器请求失败
	//			return errors.New(ERR_CVI3_REPLY)
	//		}
	//		return nil
	//	}
	//}

	var header_str string
	for i := 0; i < 6; i++ {
		header_str = c.Response.get(seq)
		if header_str != "" {
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
