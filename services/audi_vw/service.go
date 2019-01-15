package audi_vw

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/socket_listener"
	"github.com/pkg/errors"

	"encoding/xml"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"io"
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
	SN     string `json:"controller_sn"`
	Status string `json:"status"`
}

type Service struct {
	configValue   atomic.Value
	name          string
	listener      *socket_listener.SocketListener
	err           chan error
	diag          Diagnostic
	DB            *storage.Service
	WS            *wsnotify.Service
	Aiis          *aiis.Service
	Minio         *minio.Service
	wg            sync.WaitGroup
	closing       chan struct{}
	handle_buffer chan string
	Parent        *controller.Service
}

func (s *Service) Err() <-chan error {
	return s.err
}

func NewService(c Config, d Diagnostic, parent *controller.Service) *Service {
	addr := fmt.Sprintf("tcp://:%d", c.Port)

	s := &Service{
		name:    controller.AUDIPROTOCOL,
		err:     make(chan error, 1),
		diag:    d,
		wg:      sync.WaitGroup{},
		closing: make(chan struct{}, 1),
		Parent:  parent,
	}

	s.handle_buffer = make(chan string, 1024)
	lis := socket_listener.NewSocketListener(addr, s, c.ReadBufferSize*2, c.MaxConnections)
	s.listener = lis
	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (p *Service) AddNewController(cfg controller.ControllerConfig) controller.Controller {
	config := p.config()
	c := NewController(config)
	c.Srv = p //服务注入
	c.cfg = cfg

	return &c
}

func (p *Service) Write(sn string, buf []byte) error {
	if _, ok := p.Parent.Controllers[sn]; !ok {
		return fmt.Errorf("can not found controller :%s", sn)
	}
	v := p.Parent.Controllers[sn]
	c := v.(*Controller)

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

	for _, w := range p.Parent.Controllers {
		if w.Protocol() == controller.AUDIPROTOCOL {
			go w.Start() //异步启动控制器
		}
	}

	p.wg.Add(1)
	go p.HandleProcess()

	return nil
}

func (p *Service) Close() error {
	err := p.listener.Close()
	if err != nil {
		return errors.Wrapf(err, "Close Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Parent.Controllers {
		err := w.Close()
		if err != nil {
			return errors.Wrapf(err, "Close Protocol %s Writer fail", p.name)
		}
	}

	p.closing <- struct{}{}

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
	header_rest := 0

	var header_buffer string
	var header CVI3Header
	buffer := make([]byte, p.config().ReadBufferSize*2)
	for {
		n, err := c.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				p.diag.Error("Timeout in plugin AduiVW Protocol: %s", err)
				continue
			} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				p.diag.Error("using closing connection", err)
			} else if err == io.EOF {
				p.diag.Error("network Connection EOF ", err)
			}
			break
		}

		msg := string(buffer[0:n])
		//if len(msg) < HEADER_LEN {
		//	continue
		//}
		//p.diag.Debug(fmt.Sprintf("before rest: %d", rest))

		off := 0 //循环前偏移为0
		for off < n {
			if rest == 0 {
				len_msg := n - off
				if len_msg < HEADER_LEN-header_rest {
					//长度不够
					if header_rest == 0 {
						header_rest = HEADER_LEN - len_msg
					} else {
						header_rest -= len_msg
					}
					header_buffer += msg[off : off+len_msg]
					break
				} else {
					//完整
					if header_rest == 0 {
						header_buffer = msg[off : off+HEADER_LEN]
						off += HEADER_LEN
					} else {
						header_buffer += msg[off : off+header_rest]
						off += header_rest
						header_rest = 0
					}
				}
				//fmt.Printf("header rest:%d, offset:%d, n %d, header : %s\n", header_rest, off, n, header_buffer)
				header.Deserialize(header_buffer)
				header_buffer = ""
				if n-off > header.SIZ {
					//粘包
					body = msg[off : off+header.SIZ]
					p.CVIResponse(&header, c)
					p.Parse(body)
					off += header.SIZ
					rest = 0 //同样解析头
				} else {
					body = msg[off:n]
					rest = header.SIZ - (n - off)
					break
				}
			} else {
				if n-off > rest {
					//粘包
					body += string(buffer[off : off+rest]) //已经是完整的包
					p.CVIResponse(&header, c)
					p.Parse(body)
					off += rest
					rest = 0 //进入解析头
				} else {
					body += string(buffer[off:n])
					rest -= n - off
					break
				}
			}
		}

		//p.diag.Debug(fmt.Sprintf("after rest: %d", rest))

		if rest == 0 {
			p.CVIResponse(&header, c)
			p.Parse(body)
		}

		//time.Sleep(100 *time.Millisecond)
	}

}

func (p *Service) CVIResponse(header *CVI3Header, c net.Conn) {
	if header.TYP == Header_type_request_with_reply || header.TYP == Header_type_keep_alive {
		// 执行应答
		var reply CVI3Header
		reply.Init()
		reply.TYP = Header_type_reply
		reply.MID = header.MID
		replyPacket := reply.Serialize()

		_, err := c.Write([]byte(replyPacket))
		//p.diag.Debug(fmt.Sprintf("write response bytes length: %d", n))
		if err != nil {
			p.diag.Error("server reply err:%s\n", err)
		}
	}
}

func (p *Service) Parse(msg string) ([]byte, error) {

	p.handle_buffer <- msg
	return nil, nil
}

func (p *Service) HandleProcess() {
	for {
		select {
		case msg := <-p.handle_buffer:

			// 处理结果
			if strings.Contains(msg, XML_RESULT_KEY) {

				cvi3Result := CVI3Result{}
				err := xml.Unmarshal([]byte(msg), &cvi3Result)
				if err != nil {
					p.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
				} else {

					// 结果数据
					controllerResult := controller.ControllerResult{}
					XML2Result(&cvi3Result, &controllerResult)
					controllerResult.NeedPushHmi = true
					controllerResult.Raw = msg

					if strings.Contains(msg, XML_CURVE_KEY) {

						// 曲线数据
						controllerCurve := minio.ControllerCurve{}
						XML2Curve(&cvi3Result, &controllerCurve)

						p.Parent.Handle(&controllerResult, &controllerCurve)
					} else {

						p.Parent.Handle(&controllerResult, nil)
					}
				}

			}

			// 处理事件
			if strings.Contains(msg, XML_EVENT_KEY) {

				//fmt.Printf("xml evt:%s\n", msg)
				evt := Evt{}
				err := xml.Unmarshal([]byte(msg), &evt)
				if err != nil {
					p.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
				}

				// 拧紧枪状态变化
				//if strings.Contains(msg, XML_STATUS_KEY) {
				//	if evt.MSL_MSG.EVT.STS.ONC.RDY == 0 {
				//	}
				//}

				// 套同选择器事件
				//if strings.Contains(msg, XML_NUT_KEY) {
				//	// 将套筒信息推送hmi
				//	ws := wsnotify.WSSelector{}
				//	ws.SN = ""
				//	ws.Selectors = []int{}
				//	for _, v := range evt.MSL_MSG.EVT.STS.ONC.NUT.NIDs {
				//		ws.Selectors = append(ws.Selectors, nut_ids[v])
				//	}
				//	ws_str, _ := json.Marshal(ws)
				//
				//	p.WS.WSSendControllerSelectorStatus(string(ws_str))
				//}
			}

		case <-p.closing:
			p.wg.Done()
			return
		}
	}

}

// 取得控制器状态
func (p *Service) GetControllersStatus(sns []string) ([]ControllerStatus, error) {
	var status []ControllerStatus

	if len(sns) == 0 {
		for k, v := range p.Parent.Controllers {
			s := ControllerStatus{}
			s.SN = k
			s.Status = v.Status()
			status = append(status, s)
		}

		return status, nil

	} else {
		for _, v := range sns {
			c, exist := p.Parent.Controllers[v]
			if exist {
				s := ControllerStatus{}
				s.SN = v
				s.Status = c.Status()
				status = append(status, s)
			}
		}

		return status, nil
	}
}

// 设置拧接程序
func (p *Service) PSet(sn string, pset int, workorder_id int64, result_id int64, count int, user_id int64) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(ERR_CVI3_NOT_FOUND)
	}

	c := v.(*Controller)

	if c.Status() == controller.STATUS_OFFLINE {
		// 控制器离线
		return errors.New(ERR_CVI3_OFFLINE)
	}

	// 设定pset并判断控制器响应
	_, err := c.PSet(pset, workorder_id, result_id, count, user_id, c.cfg.ToolChannel)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil

}

// 拧紧抢使能控制
func (p *Service) ToolControl(sn string, enable bool) error {
	// 判断控制器是否存在
	v, exist := p.Parent.Controllers[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(ERR_CVI3_NOT_FOUND)
	}

	c := v.(*Controller)

	if c.Status() == controller.STATUS_OFFLINE {
		// 控制器离线
		return errors.New(ERR_CVI3_OFFLINE)
	}

	// 使能控制
	err := c.ToolControl(enable, c.cfg.ToolChannel)
	if err != nil {
		// 控制器请求失败
		return err
	}

	return nil

}
