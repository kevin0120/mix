package masterplc

import (
	"sync/atomic"
	"github.com/masami10/aiis/services/rush"
	"github.com/masami10/aiis/socket_listener"
	"fmt"
	"net"
	"strings"
	"io"
	"encoding/xml"
	"sync"
	"github.com/pkg/errors"
	"github.com/masami10/aiis/services/minio"
		"github.com/masami10/aiis/services/storage"
	"encoding/json"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag           Diagnostic
	configValue    atomic.Value
	Rush			*rush.Service
	Minio			*minio.Service
	listener      *socket_listener.SocketListener
	handleBuffer chan string
	wg            sync.WaitGroup
	closing       chan struct{}
}

func NewService(d Diagnostic, c Config) *Service {
	if c.Enable {
		s := &Service{
			diag:      d,
			wg:      sync.WaitGroup{},
			closing: make(chan struct{}, 1),
		}

		s.configValue.Store(c)

		addr := fmt.Sprintf("tcp://:%d", c.Port)
		s.handleBuffer = make(chan string, 1024)
		lis := socket_listener.NewSocketListener(addr, s, c.ReadBufferSize*2, c.MaxConnections)
		s.listener = lis

		return s
	}


	return nil
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	err := s.listener.Start()
	if err != nil {
		return errors.Wrapf(err, "Open %s Listener fail", "masterplc")
	}

	s.wg.Add(1)

	for i := 0; i < s.Config().Workers; i++ {
		s.diag.Debug(fmt.Sprintf("init masterplc worker %d\n", i + 1))
		go s.HandleProcess()
	}

	return nil
}

func (s *Service) Close() error {
	err := s.listener.Close()
	if err != nil {
		return errors.Wrapf(err, "Close %s Listener fail", "masterplc")
	}

	s.closing <- struct{}{}

	s.wg.Wait() //阻塞 等待全部handler关闭

	return nil
}

func (s *Service) NewConn(c net.Conn) {
	s.diag.Debug(fmt.Sprintf("new conn %s:\n", c.RemoteAddr()))
}

func (s *Service) Parse(msg string) ([]byte, error) {

	s.handleBuffer <- msg
	return nil, nil
}

// 服务端读取
func (s *Service) Read(c net.Conn) {
	ssl := s.listener
	defer ssl.InterListener.RemoveConnection(c)
	defer c.Close()

	rest := 0
	body := ""
	header_rest := 0

	var header_buffer string
	var header CVI3Header
	buffer := make([]byte, s.Config().ReadBufferSize*2)
	for {
		n, err := c.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.diag.Error("Timeout in plugin AduiVW Protocol: %s", err)
				continue
			} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				s.diag.Error("using closing connection", err)
			} else if err == io.EOF {
				s.diag.Error("network Connection EOF ", err)
			}
			break
		}

		msg := string(buffer[0:n])

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
					s.CVIResponse(&header, c)
					s.Parse(body)
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
					s.CVIResponse(&header, c)
					s.Parse(body)
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
			s.CVIResponse(&header, c)
			s.Parse(body)
		}

		//time.Sleep(100 *time.Millisecond)
	}

}

func (s *Service) CVIResponse(header *CVI3Header, c net.Conn) {
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
			s.diag.Error("server reply err:%s\n", err)
		}
	}
}

func (s *Service) HandleProcess() {
	for {
		select {
		case msg := <-s.handleBuffer:

			// 处理结果
			//fmt.Printf("%s\n", msg)
			s.diag.Debug(fmt.Sprintf("masterplc result: %s\n", msg))
			cvi3Result := CVI3Result{}
			err := xml.Unmarshal([]byte(msg), &cvi3Result)
			if err != nil {
				s.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
				return
			}

			// 结果数据
			opResult := storage.OperationResult{}
			XML2Result(&cvi3Result, &opResult)

			cr := rush.CResult{
				Result:      &opResult,
				ID:     0,
				Stream: nil,
			}

			s.Rush.AddResultTask(cr)

			// 处理曲线
			if strings.Contains(msg, XML_RESULT_KEY) {
				curveFile := fmt.Sprintf("masterplc_%d.json", opResult.TighteningId)
				cur := storage.CURObject {
					File: curveFile,
					OP: 1,
				}

				opResult.CurObjects = append(opResult.CurObjects, cur)

				// 波形文件
				controllerCurveFile := ControllerCurveFile{}
				XML2Curve(&cvi3Result, &controllerCurveFile)

				// 保存波形到对象存储
				s.diag.Debug("保存波形数据到对象存储 ...")
				strCurveData, _ := json.Marshal(controllerCurveFile)

				err = s.Minio.Upload(curveFile, string(strCurveData))
				if err != nil {
					s.diag.Error("对象存储保存失败", err)
				} else {
					s.diag.Debug("对象存储保存成功")
				}
			}



			// 处理事件
			if strings.Contains(msg, XML_EVENT_KEY) {

				//fmt.Printf("xml evt:%s\n", msg)
				evt := Evt{}
				err := xml.Unmarshal([]byte(msg), &evt)
				if err != nil {
					s.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
					return
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

		case <-s.closing:
			s.wg.Done()
			return
		}
	}

}