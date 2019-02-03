package changan

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/wsnotify"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	ReciveNewTask(msg string)
}

type Service struct {
	WS              *wsnotify.Service
	HTTPDService    *httpd.Service
	Conn            net.Conn
	Opened          bool
	ReadTimeout     time.Duration
	diag            Diagnostic
	Seq             int
	configValue     atomic.Value
	Msgs            chan AndonMsg
	requestTaskInfo chan string
	stop            chan chan struct{}
	AndonDB         *AndonDB
}

func (c *Service) GetSequenceNum() int {
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
		},
	}

	s.configValue.Store(c)
	return s

}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
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

	spl := strings.SplitN(c.AndonAddr, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid address: %s", c.AndonAddr)
	}

	go s.Connect(spl)

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

func (s *Service) Connect(spl []string) {

	var con net.Conn
	var err error
	for {
		con, err = net.DialTimeout(spl[0], spl[1], 3*time.Second)
		if err == nil {
			break
		}
	}

	s.Conn = con

	// 设置keep alive
	s.setKeepAlive(con)

	// 注册aiis 服务到andon系统中
	if err = s.write(PakcageMsg(MSG_REGIST, s.GetSequenceNum(), nil)); err != nil {
		s.diag.Error("Registry to Andon fail", err)
		//return err
	}

	go s.manage()

	go s.readHandler(s.Conn)

	s.Opened = true
}

func (s *Service) write(buf []byte) error {
	if s.Conn == nil {
		return nil
	}

	s.Conn.SetWriteDeadline(time.Now().Add(s.ReadTimeout))

	_, err := s.Conn.Write(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			s.diag.Error(fmt.Sprintf("D! Timeout in Write: %s",err.Error()), err)
			return err
		} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
			s.diag.Error(fmt.Sprintf("D!: %s",err.Error()), err)
			return err
		}
	}

	return nil
}

func (s *Service) manage() {

	c := s.configValue.Load().(Config)

	for {
		select {
		case <-time.After(time.Duration(c.KeepAlivePeriod)):
			err := s.write(PakcageMsg(MSG_HEART, s.GetSequenceNum(), nil))
			if err != nil {
				s.diag.Error("send keep Alive msg fail", err)
			}
		case msg := <-s.requestTaskInfo:
			//为了不会永远将其阻塞,超过timeout后的应答将会跑到此处
			fmt.Printf("should never go here : %s", msg)
		case msg := <-s.Msgs:
			switch msg.MsgType {
			case MSG_HEART_ACK:
				fmt.Printf("heart beat seq : %d\n", msg.Seq)
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
				s.write(PakcageMsg(MSG_TASK_ACK, s.GetSequenceNum(), nil))
			case MSG_GET_TASK_ACK:
				strData, err := json.Marshal(msg.Data)
				if err != nil {
					s.requestTaskInfo <- fmt.Sprintf("error: %s", err.Error())
				}

				s.requestTaskInfo <- string(strData[:])
			case MSG_GUID_REQ:
				d := AndonGUID{GUID: c.GUID}

				s.write(PakcageMsg(MSG_GUID_REQ_ACK, s.GetSequenceNum(), d))
			default:
				fmt.Println("not support msg type")
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

func (s *Service) readHandler(c net.Conn) {
	defer s.Open()
	defer s.Close()

	conf := s.configValue.Load().(Config)

	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		fmt.Errorf("cannot set keep alive on a %s socket", conf.AndonAddr)
	}

	tcpc.SetReadBuffer(2 * conf.ReadBufferSize) // 设定读取的buffer 大小
	d := json.NewDecoder(tcpc)

	for {
		var msg AndonMsg
		err := d.Decode(&msg)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.diag.Error("Timeout in changan aiis andon Protocol: %s", err)
				continue
			} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				s.diag.Error("using closing connection", err)
				return
			} else if err == io.EOF {
				s.diag.Error("network Connection EOF ", err)
				return
			}
			break
		}
		// 发送到通道中进行处理
		s.Msgs <- msg
		////打印测试
		//fmt.Println(msg)

	}
}

func (s *Service) andonGetTaskbyworkCenter(ctx iris.Context) {
	c := s.configValue.Load().(Config)
	workcenter := ctx.Params().Get("workcenter")

	payload := PakcageMsg(MSG_GET_TASK, s.GetSequenceNum(), AndonWorkCenter{Workcenter: workcenter})

	if err := s.write(payload); err != nil {
		ctx.Writef(fmt.Sprintf("Try to get workcenter: %s task fail", workcenter))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

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
