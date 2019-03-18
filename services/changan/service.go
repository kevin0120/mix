package changan

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/wsnotify"
	"github.com/masami10/aiis/socket_writer"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	ReciveNewTask(msg string)
}

type Service struct {
	WS                *wsnotify.Service
	HTTPDService      *httpd.Service
	Conn              net.Conn
	Opened            bool
	ReadTimeout       time.Duration
	diag              Diagnostic
	Seq               int
	mtxSeq            sync.Mutex
	responses         ResponseQueue
	configValue       atomic.Value
	Msgs              chan AndonMsg
	requestTaskInfo   chan string
	stop              chan chan struct{}
	AndonDB           *AndonDB
	w                 *socket_writer.SocketWriter
	writeBuffer       chan []byte
	StatusValue       atomic.Value
	keepAliveCount    int32
	keep_period       time.Duration
	keepaliveDeadLine atomic.Value
}

func (c *Service) GetSequenceNum() int {
	defer c.mtxSeq.Unlock()
	c.mtxSeq.Lock()

	x := c.Seq
	if x >= 9999 {
		c.Seq = 1
	} else {
		c.Seq += 1
	}
	return x
}

func NewService(d Diagnostic, c Config, h *httpd.Service, ws *wsnotify.Service) *Service {

	s := &Service{
		diag:            d,
		HTTPDService:    h,
		WS:              ws,
		Msgs:            make(chan AndonMsg, 100),
		requestTaskInfo: make(chan string, 1),
		Seq:             1,
		Opened:          false,
		ReadTimeout:     time.Duration(c.ReadTimeout),
		AndonDB: &AndonDB{
			cfg: &c.DB,
			eng: nil,
		},
		mtxSeq:      sync.Mutex{},
		responses:   ResponseQueue{},
		writeBuffer: make(chan []byte, 1024),
		keep_period: time.Duration(3 * time.Second),
	}

	s.responses.Init()
	s.configValue.Store(c)
	return s

}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Write(buf []byte) {
	s.writeBuffer <- buf
}

func (s *Service) setKeepAlive(con net.Conn) error {

	c := s.configValue.Load().(Config)

	tcpc, ok := con.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set keep alive on a %s socket", c.AndonAddr)
	}
	KeepAlivePeriod := time.Duration(c.KeepAlivePeriod)
	if KeepAlivePeriod.Nanoseconds() == 0 {
		return tcpc.SetKeepAlive(false)
	}
	if err := tcpc.SetKeepAlive(true); err != nil {
		return err
	}
	return tcpc.SetKeepAlivePeriod(KeepAlivePeriod)
}

func (s *Service) Open() error {

	//开始配置
	c := s.configValue.Load().(Config)

	if !c.Enable {
		return nil
	}

	s.AndonDB.StartService()

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/operation/{workcenter:string}",
		HandlerFunc: s.andonGetTaskbyworkCenter,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/andon-test",
		HandlerFunc: s.andonTest,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/vehicle",
		HandlerFunc: s.getVehicle,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	spl := strings.SplitN(c.AndonAddr, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid address: %s", c.AndonAddr)
	}

	s.StatusValue.Store(ANDON_STATUS_OFFLINE)

	s.w = socket_writer.NewSocketWriter(s.Config().AndonAddr, nil)
	s.w.OnRead = s.OnRead

	go s.Connect()

	return nil
}

func (s *Service) Close() error {
	if !s.Opened {
		return nil
	}

	s.AndonDB.StopService()
	stopping := make(chan struct{})
	s.stop <- stopping

	<-stopping
	s.Conn = nil
	s.Opened = false

	return nil
}

func (s *Service) Status() string {

	return s.StatusValue.Load().(string)
}

func (s *Service) updateStatus(status string) {

	if status != s.Status() {

		s.StatusValue.Store(status)

		if status == ANDON_STATUS_OFFLINE {
			s.Close()

			// 断线重连
			go s.Connect()
		}

	}
}

func (s *Service) Connect() {
	s.Seq = 1
	s.StatusValue.Store(ANDON_STATUS_OFFLINE)
	s.updateKeepAliveCount(0)
	s.diag.Debug("andon connecting ...")
	//var con net.Conn
	var err error
	for {
		err = s.w.Connect(time.Duration(1 * time.Second))
		if err == nil {
			s.diag.Debug("andon connected")
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	s.updateStatus(ANDON_STATUS_ONLINE)
	s.updateKeepAliveDeadLine()
	//con = s.w.Conn
	// 设置keep alive
	//s.setKeepAlive(con)

	// 注册aiis 服务到andon系统中
	s.Write(PakcageMsg(MSG_REGIST, s.GetSequenceNum(), AndonGUID{GUID: s.Config().GUID}))

	//go s.readHandler(con)
	go s.manage()

	//s.Opened = true
}

//func (s *Service) write(buf []byte) error {
//	//if s.w.Conn == nil {
//	//	return nil
//	//}
//
//	//s.Conn.SetWriteDeadline(time.Now().Add(s.ReadTimeout))
//
//	err := s.w.Write(buf)
//	if err != nil {
//		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
//			s.diag.Error(fmt.Sprintf("D! Timeout in Write: %s", err.Error()), err)
//			return err
//		} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
//			s.diag.Error(fmt.Sprintf("D!: %s", err.Error()), err)
//			return err
//		}
//	}
//
//	return nil
//}

func (s *Service) sendKeepalive() {
	if s.Status() == ANDON_STATUS_OFFLINE {
		return
	}

	s.Write(PakcageMsg(MSG_HEART, s.GetSequenceNum(), nil))
}

func (s *Service) manage() {

	c := s.configValue.Load().(Config)

	for {
		select {
		//case <-time.After(time.Duration(c.KeepAlivePeriod)):
		//	s.Write(PakcageMsg(MSG_HEART, s.GetSequenceNum(), nil))

		case <-time.After(s.keep_period):
			if s.Status() == ANDON_STATUS_OFFLINE {
				continue
			}
			if s.KeepAliveCount() >= MAX_KEEP_ALIVE_CHECK {
				go s.updateStatus(ANDON_STATUS_OFFLINE)
				s.updateKeepAliveCount(0)
				continue
			}
			if s.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				s.sendKeepalive()
				s.updateKeepAliveDeadLine() //更新keepalivedeadline
				s.addKeepAliveCount()
			}

		case v := <-s.writeBuffer:
			//for nextWriteThreshold.After(time.Now()) {
			//	time.Sleep(time.Microsecond * 100)
			//}
			err := s.w.Write([]byte(v))
			if err != nil {
				s.diag.Debug("andon write fail")
				//c.Srv.diag.Error("Write data fail", err)
			} else {
				s.updateKeepAliveDeadLine()
				//c.updateKeepAliveDeadLine()
			}
			//nextWriteThreshold = time.Now().Add(c.req_timeout)
		case msg := <-s.requestTaskInfo:
			//为了不会永远将其阻塞,超过timeout后的应答将会跑到此处
			fmt.Printf("should never go here : %s", msg)
		case msg := <-s.Msgs:

			switch msg.MsgType {
			case MSG_HEART_ACK:
				s.diag.Debug(fmt.Sprintf("heart beat seq : %d\n", msg.Seq))
			case MSG_TASK:
				strData, _ := json.Marshal(msg.Data)
				s.diag.ReciveNewTask(string(strData))
				var tasks []AndonTask
				err := json.Unmarshal(strData, &tasks)
				if err != nil {
					break
				}

				for _, v := range tasks {
					t, _ := json.Marshal(v)
					s.WS.WSSendTask(v.Workcenter, string(t))
				}
				s.Write(PakcageMsg(MSG_TASK_ACK, msg.Seq, nil))
			case MSG_GET_TASK_ACK:
				strData, err := json.Marshal(msg.Data)
				if err != nil {
					s.requestTaskInfo <- fmt.Sprintf("error: %s", err.Error())
				}

				s.requestTaskInfo <- string(strData[:])
			case MSG_GUID_REQ:
				d := AndonGUID{GUID: c.GUID}

				s.Write(PakcageMsg(MSG_GUID_REQ_ACK, msg.Seq, d))

			case MSG_REGIST_ACK:

			case MSG_VEHICLE_REQ_ACK:
				if msg.Data != nil {
					workcenterCode := s.responses.getAndRemove(msg.Seq)
					strData, _ := json.Marshal(msg.Data)
					v := AndonTask{Workcenter: workcenterCode.(string)}
					err := json.Unmarshal(strData, &v)
					if err == nil {
						t, _ := json.Marshal(v)
						s.WS.WSSendTask(v.Workcenter, string(t))
					}
				}
			default:
				str, _ := json.Marshal(msg)
				s.diag.Debug(fmt.Sprintf("not support msg type:", string(str)))
			}

		case stop := <-s.stop:
			if s.Conn != nil {
				s.Conn.Close()
			}
			close(stop)
			return
		}
	}

}

func (s *Service) OnRead(c net.Conn) {

	defer c.Close()

	buffer := make([]byte, 65535)

	for {
		n, err := c.Read(buffer)
		if err != nil {
			break
		}

		s.updateKeepAliveCount(0)

		str := string(buffer[0:n])
		s.diag.Debug(fmt.Sprintf("andon recv:%s\n", str))

		var msg AndonMsg

		err = json.Unmarshal(buffer[0:n], &msg)
		if err == nil {
			s.Msgs <- msg
		}
	}
}

func (s *Service) andonGetTaskbyworkCenter(ctx iris.Context) {
	c := s.configValue.Load().(Config)
	workcenter := ctx.Params().Get("workcenter")

	workcenters := []string{workcenter}
	payload := PakcageMsg(MSG_GET_TASK, s.GetSequenceNum(), workcenters)

	s.Write(payload)

	select {
	case <-time.After(time.Duration(3 * c.ReadTimeout)):
		ctx.Writef(fmt.Sprintf("Try to get workcenter: %s task fail, Timeout!", workcenter))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	case msg := <-s.requestTaskInfo:
		if strings.HasPrefix(msg, "error:") {
			// error happen
			ctx.Writef(msg)
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		var task AndonTask

		if err := json.Unmarshal([]byte(msg), &task); err != nil {
			//错误的返回值
			ctx.Writef(fmt.Sprintf("not corrent task struct: %s", msg))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		s.WS.WSSendTask(workcenter, msg) //发送到指定工位,相应的作业内容

		ctx.StatusCode(iris.StatusOK)
		return
	}
}

func (s *Service) andonTest(ctx iris.Context) {
	andon_msg := AndonMsg{}
	err := ctx.ReadJSON(&andon_msg)

	if err != nil {
		ctx.Writef(fmt.Sprintf("param error: %s", err.Error()))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	str_data, _ := json.Marshal(andon_msg.Data)

	switch andon_msg.MsgType {
	case MSG_TASK:
		tasks := []AndonTask{}
		err = json.Unmarshal(str_data, &tasks)
		if err != nil {
			ctx.Writef(fmt.Sprintf("tasks error: %s", err.Error()))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		for _, v := range tasks {
			t, _ := json.Marshal(v)
			s.WS.WSSendTask(v.Workcenter, string(t))

			//fmt.Printf("send task -- workcenter:%s payload:%s\n", v.Workcenter, string(t))
		}
	}
}

func (s *Service) getVehicle(ctx iris.Context) {
	vin := ctx.URLParam("vin")
	if vin == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("vin is required")
		return
	}

	workcenterCode := ctx.URLParam("workcenter_code")
	if workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workcenter_code is required")
		return
	}

	// 和andon通信，根据vin查询对应车辆, 收到响应后通过websocket推送车辆信息给hmi
	seq := s.GetSequenceNum()
	s.responses.update(seq, workcenterCode)
	payload := PakcageMsg(MSG_VEHICLE_REQ, seq, AndonVehicle{Vin: vin})
	s.Write(payload)

	ctx.StatusCode(iris.StatusNoContent)

}

func (s *Service) KeepAliveCount() int32 {
	return atomic.LoadInt32(&s.keepAliveCount)
}

func (s *Service) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&s.keepAliveCount, i)
}

func (s *Service) addKeepAliveCount() {
	atomic.AddInt32(&s.keepAliveCount, 1)
}

func (s *Service) updateKeepAliveDeadLine() {
	s.keepaliveDeadLine.Store(time.Now().Add(s.keep_period))
}

func (s *Service) KeepAliveDeadLine() time.Time {
	return s.keepaliveDeadLine.Load().(time.Time)
}
