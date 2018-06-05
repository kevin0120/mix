package audi_vw

import (
	"github.com/masami10/rush/socket_listener"
	"sync"
	"github.com/masami10/rush/services/controller"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync/atomic"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"encoding/xml"
	"github.com/masami10/rush/payload"
	"encoding/json"
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


func (p *Service) AddNewController(cfg controller.Config) Controller{
	config := p.config()
	c := NewController(config)
	c.Srv = p //服务注入
	c.cfg = cfg
	p.Controllers[cfg.SN] = c

	return c
}

func (p *Service) Write(serial_no string ,buf []byte)  error{
	//p.mux.Lock()
	ws := p.Controllers
	//p.mux.Unlock()

	if _, ok := ws[serial_no]; !ok{
		return fmt.Errorf("Controller %s not found ", serial_no)
	}

	//w := ws[serial_no]

	//w.Write(string(buf), w.get_sequence())

	return nil
}

func (p *Service) Open() error {

	err := p.listener.Start()
	if err != nil {
		return errors.Wrapf(err, "Open Protocol %s Listener fail", p.name)
	}

	for _, w := range p.Controllers{
		go w.Start()
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


// 服务端接受新连接
func (p *Service) NewConn(c net.Conn) {
	p.setRemoteConn(c)
}

func (p *Service) setRemoteConn(c net.Conn) (string) {
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	for k,v := range p.Controllers {
		if v.cfg.RemoteIP == ip {
			v.RemoteConn = c
			return k
		}
	}

	return ""
}

// 服务端读取
func (p *Service) Read(c net.Conn) {
	ssl := p.listener
	defer ssl.InterListener.RemoveConnection(c)
	defer c.Close()

	//buff := make([]byte, 4 * 4096) //创建一个打的buffer
	//
	//offset := 0
	//
	//scnr := bufio.NewScanner(c)
	//for {
	//	if ssl.ReadTimeout.Nanoseconds() > 0 {
	//		c.SetReadDeadline(time.Now().Add(ssl.ReadTimeout))
	//	}
	//	if !scnr.Scan() {
	//		break
	//	}
	//	copy(buff[offset:], scnr.Bytes())
	//	buf, err := p.Parse(buff)
	//	if err != nil {
	//		log.Printf("unable to parse incoming line: %s", err)
	//		//TODO rate limit
	//		continue
	//	}
	//	if buf == nil {
	//		//猜测还需要继续读取数据
	//		offset = len(buf)
	//		continue
	//	}
	//	//解析成功，进行操作
	//	offset = 0 //偏移量回到0
	//}
	//
	//if err := scnr.Err(); err != nil {
	//	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
	//		log.Printf("D! Timeout in plugin [input.socket_listener]: %s", err)
	//	} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
	//		log.Print(err)
	//	}
	//}

	buffer := make([]byte, 65535)
	for {

		n, err := c.Read(buffer)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}
		msg := string(buffer[0:n])
		if len(msg) < HEADER_LEN {
			continue
		}
		//fmt.Printf("%s\n", msg)

		header := CVI3Header{}
		header.Deserialize(msg[0: HEADER_LEN])
		var body string = msg[HEADER_LEN: n]
		var rest int = int(header.SIZ) - HEADER_LEN - n
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

		go p.Parse([]byte(body))

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

	if strings.Contains(msg, cvi3.XML_RESULT_KEY) {
		fmt.Printf("收到结果数据:%s\n", msg)

		result := cvi3.CVI3Result{}
		err := xml.Unmarshal([]byte(msg), &result)
		if err != nil {
			fmt.Printf("OnRecv err:%s\n", err.Error())
		}

		// 结果数据
		result_data := payload.XML2Result(result)

		// 波形文件
		curve_file := payload.XML2Curve(result)

		curve := payload.ControllerCurve{}
		s_curvedata, _ := json.Marshal(curve_file)
		curve.CurveData = string(s_curvedata)
		curve.Count = result_data.Count
		curve.CurveFile = result_data.CurFile
		curve.ResultID = result_data.Result_id

		//e := service.HandleCurve(curve)
		//if  e == nil {
		//	go service.HandleResult(result_data)
		//} else {
		//	fmt.Printf("OnRecv err:%s\n", e.Error())
		//}

	} else {
		fmt.Printf("recv:%s\n", msg)
	}

	return nil, nil
}

//func (service *CVI3Service) HandleResult(result payload.ControllerResult) (error) {
//	fmt.Printf("处理结果数据...\n")
//
//	var err error
//	r, err := service.DB.GetResult(result.Result_id, 0)
//	if err != nil {
//		return err
//	}
//
//	// 保存结果
//	loc, _ := time.LoadLocation("Local")
//	times := strings.Split(result.Dat, " ")
//	dt := fmt.Sprintf("%s %s", times[0], times[1])
//	r.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", dt, loc)
//	r.Result = result.Result
//	r.Count = result.Count
//	r.HasUpload = false
//	r.ControllerSN = result.Controller_SN
//	s_value, _ := json.Marshal(result.ResultValue)
//	s_pset, _ := json.Marshal(result.PSetDefine)
//
//	r.ResultValue = string(s_value)
//	r.PSetDefine = string(s_pset)
//
//	fmt.Printf("保存结果到数据库\n")
//	_, err = service.DB.UpdateResult(r)
//	if err != nil {
//		fmt.Printf("HandleResult err:%s\n", err.Error())
//		return nil
//	}
//
//	workorder, err := service.DB.GetWorkorder(result.Workorder_ID)
//	if err == nil {
//		// 结果推送hmi
//		ws_result := payload.WSResult{}
//		ws_result.Result_id = result.Result_id
//		ws_result.Count = result.Count
//		ws_result.Result = result.Result
//		ws_result.MI = result.ResultValue.Mi
//		ws_result.WI = result.ResultValue.Wi
//		ws_result.TI = result.ResultValue.Ti
//		ws_str, _ := json.Marshal(ws_result)
//
//		fmt.Printf("Websocket推送结果到HMI\n")
//		go service.APIService.WSSendResult(workorder.HMISN, string(ws_str))
//	}
//
//	if r.Count >= int(workorder.MaxRedoTimes) || r.Result == payload.RESULT_OK {
//		// 结果推送AIIS
//
//		odoo_result := payload.ODOOResult{}
//		if r.Result == payload.RESULT_OK {
//			odoo_result.Final_pass = "pass"
//			if r.Count == 1 {
//				odoo_result.One_time_pass = "pass"
//			} else {
//				odoo_result.One_time_pass = "fail"
//			}
//		} else {
//			odoo_result.Final_pass = "fail"
//			odoo_result.One_time_pass = "fail"
//		}
//
//		odoo_result.Control_date = fmt.Sprintf("%sT%s+08:00", times[0], times[1])
//
//		odoo_result.Measure_degree = result.ResultValue.Wi
//		odoo_result.Measure_result = strings.ToLower(result.Result)
//		odoo_result.Measure_t_don = result.ResultValue.Ti
//		odoo_result.Measure_torque = result.ResultValue.Mi
//		odoo_result.Op_time = result.Count
//		odoo_result.Pset_m_max = result.PSetDefine.Mp
//		odoo_result.Pset_m_min = result.PSetDefine.Mm
//		odoo_result.Pset_m_target = result.PSetDefine.Ma
//		odoo_result.Pset_m_threshold = result.PSetDefine.Ms
//		odoo_result.Pset_strategy = result.PSetDefine.Strategy
//		odoo_result.Pset_w_max = result.PSetDefine.Wp
//		odoo_result.Pset_w_min = result.PSetDefine.Wm
//		odoo_result.Pset_w_target = result.PSetDefine.Wa
//
//		curves, err := service.DB.ListCurves(result.Result_id)
//		if err != nil {
//			return err
//		}
//		for _, v := range curves {
//			curobject := payload.CURObject{}
//			curobject.OP = v.Count
//			curobject.File = v.CurveFile
//			odoo_result.CURObjects = append(odoo_result.CURObjects, curobject)
//		}
//
//		fmt.Printf("推送结果数据到AIIS\n")
//		_, err = service.ODOO.PutResult(result.Result_id, odoo_result)
//		if err == nil {
//			// 发送成功
//			r.HasUpload = true
//			fmt.Printf("推送成功，更新本地结果标识\n")
//			_, err := service.DB.UpdateResult(r)
//			if err != nil {
//				return err
//			}
//		} else {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (service *CVI3Service) HandleCurve(curve payload.ControllerCurve) (error) {
//	fmt.Printf("处理波形数据...\n")
//
//	// 保存波形到数据库
//	c := rushdb.Curves{}
//	c.ResultID = curve.ResultID
//	c.CurveData = curve.CurveData
//	c.CurveFile = curve.CurveFile
//	c.Count = curve.Count
//	c.HasUpload = false
//
//	exist, err := service.DB.CurveExist(c)
//	if err != nil {
//		return err
//	} else {
//		fmt.Printf("缓存波形数据到数据库\n")
//		if exist {
//			_, err := service.DB.UpdateCurve(c)
//			if err != nil {
//				return err
//			}
//		} else {
//			err := service.DB.InsertCurve(c)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	// 保存波形到对象存储
//	fmt.Printf("保存波形数据到对象存储\n")
//	err = service.Storage.Upload(curve.CurveFile, curve.CurveData)
//	if err != nil {
//		return err
//	} else {
//		c.HasUpload = true
//		fmt.Printf("对象存储保存成功，更新本地结果标识\n")
//		_, err = service.DB.UpdateCurve(c)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//
//}
