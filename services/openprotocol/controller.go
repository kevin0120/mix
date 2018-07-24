package openprotocol

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/socket_writer"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DAIL_TIMEOUT         = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK = 3

	REPLY_TIMEOUT   = time.Duration(300 * time.Millisecond)
	MAX_REPLY_COUNT = 10
)

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}

type Controller struct {
	w                 *socket_writer.SocketWriter
	cfg               controller.Config
	StatusValue       atomic.Value
	keepAliveCount    int32
	keep_period       time.Duration
	req_timeout       time.Duration
	Response          ResponseQueue
	Srv               *Service
	dbController      *storage.Controllers
	buffer            chan []byte
	closing           chan chan struct{}
	handlerBuf        chan handlerPkg
	keepaliveDeadLine atomic.Value
	protocol          string
	Mode              atomic.Value
}

func NewController(c Config) Controller {

	cont := Controller{
		buffer:      make(chan []byte, 1024),
		closing:     make(chan chan struct{}),
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
		Response:    ResponseQueue{},
		handlerBuf:  make(chan handlerPkg, 1024),
		protocol:    controller.OPENPROTOCOL,
	}

	cont.StatusValue.Store(controller.STATUS_OFFLINE)

	return cont
}

func (c *Controller) LoadController(controller *storage.Controllers) {
	c.dbController = controller
}

func (c *Controller) handlerProcess() {
	for {
		select {
		case pkg := <-c.handlerBuf:
			c.HandleMsg(&pkg)
		}
	}
}

func (c *Controller) HandleMsg(pkg *handlerPkg) {
	c.Srv.diag.Debug(fmt.Sprintf("%s%s\n", pkg.Header.Serialize(), pkg.Body))

	switch pkg.Header.MID {
	case MID_0061_LAST_RESULT:
		// 结果数据

		result_data := ResultData{}
		result_data.Deserialize(pkg.Body)
		c.handleResult(&result_data)

	case MID_0065_OLD_DATA:
		// 历史结果数据

		result_data := ResultData{}
		result_data.DeserializeOld(pkg.Body)

		flag := c.Response.get(MID_0064_OLD_SUBSCRIBE)

		if flag != nil {
			defer c.Response.remove(MID_0064_OLD_SUBSCRIBE)
			c.Response.Add(MID_0065_OLD_DATA, result_data)
		} else {
			// 处理历史数据
			pset_detail, err := c.GetPSetDetail(result_data.PSetID)
			if err == nil {
				result_data.TorqueMin = pset_detail.TorqueMin
				result_data.TorqueMax = pset_detail.TorqueMax
				result_data.TorqueFinalTarget = pset_detail.TorqueTarget
				result_data.AngleMax = pset_detail.AngleMax
				result_data.AngleMin = pset_detail.AngleMin
				result_data.FinalAngleTarget = pset_detail.AngleTarget
			}

			c.handleResult(&result_data)
		}

	case MID_0013_PSET_DETAIL_REPLY:
		// pset详细数据
		pset_detail := PSetDetail{}
		pset_detail.Deserialize(pkg.Body)

		c.Response.update(MID_0013_PSET_DETAIL_REPLY, pset_detail)

	case MID_0011_PSET_LIST_REPLY:
		// pset列表
		pset_list := PSetList{}
		pset_list.Deserialize(pkg.Body)

		c.Response.update(MID_0011_PSET_LIST_REPLY, pset_list)

		//case MID_7410_LAST_CURVE:
		// 处理波形
	}
}

func (c *Controller) handleResult(result_data *ResultData) {

	c.Srv.DB.UpdateTightning(c.dbController.Id, result_data.TightingID)

	controllerResult := controller.ControllerResult{}
	id_info := result_data.VIN + result_data.ID2 + result_data.ID3 + result_data.ID4
	kvs := strings.Split(id_info, "-")

	if len(kvs) == 2 {
		// job模式

		controllerResult.Workorder_ID, _ = strconv.ParseInt(kvs[0], 10, 64)
		controllerResult.UserID, _ = strconv.ParseInt(kvs[1], 10, 64)

		db_result, err := c.Srv.DB.FindTargetResultForJob(controllerResult.Workorder_ID)
		if err != nil {
			c.Srv.diag.Error("FindTargetResultForJob failed", err)
		}

		controllerResult.Count = db_result.Count + 1
		controllerResult.Result_id = db_result.ResultId

	} else {
		// pset模式

		// 结果id-拧接次数-用户id
		controllerResult.Result_id, _ = strconv.ParseInt(kvs[0], 10, 64)
		controllerResult.Count, _ = strconv.Atoi(kvs[1])
		controllerResult.UserID, _ = strconv.ParseInt(kvs[2], 10, 64)

		db_result, err := c.Srv.DB.GetResult(controllerResult.Result_id, 0)
		if err != nil {
			c.Srv.diag.Error("GetResult failed", err)
		}

		controllerResult.Workorder_ID = db_result.WorkorderID
	}

	dat_kvs := strings.Split(result_data.TimeStamp, ":")
	controllerResult.Dat = fmt.Sprintf("%s %s:%s:%s", dat_kvs[0], dat_kvs[1], dat_kvs[2], dat_kvs[3])

	controllerResult.PSet = result_data.PSetID
	controllerResult.Controller_SN = c.cfg.SN
	if result_data.TighteningStatus == "0" {
		controllerResult.Result = storage.RESULT_NOK
	} else {
		controllerResult.Result = storage.RESULT_OK
	}

	controllerResult.ResultValue.Mi = result_data.Torque / 100
	controllerResult.ResultValue.Wi = result_data.Angle
	//controllerResult.ResultValue.Ti = result_data.

	switch result_data.Strategy {
	case "01":
		controllerResult.PSetDefine.Strategy = controller.STRATEGY_AW

	case "02":
		controllerResult.PSetDefine.Strategy = controller.STRATEGY_AW

	case "03":
		controllerResult.PSetDefine.Strategy = controller.STRATEGY_ADW

	case "04":
		controllerResult.PSetDefine.Strategy = controller.STRATEGY_AD
	}

	controllerResult.PSetDefine.Mp = result_data.TorqueMax / 100
	controllerResult.PSetDefine.Mm = result_data.TorqueMin / 100
	controllerResult.PSetDefine.Ma = result_data.TorqueFinalTarget / 100

	controllerResult.PSetDefine.Wp = result_data.AngleMax
	controllerResult.PSetDefine.Wm = result_data.AngleMin
	controllerResult.PSetDefine.Wa = result_data.FinalAngleTarget

	controllerResult.ExceptionReason = result_data.TighteningStatus

	c.Srv.Parent.Handle(controllerResult, nil)
}

func (c *Controller) Start() {

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	go c.handlerProcess()

	c.Connect()
}

func (c *Controller) Protocol() string {
	return c.protocol
}

func (c *Controller) Connect() error {
	c.StatusValue.Store(controller.STATUS_OFFLINE)
	c.Response = ResponseQueue{
		Results: map[string]interface{}{},
		mtx:     sync.Mutex{},
	}

	for {
		err := c.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.Srv.diag.Error("connect err", err)
		} else {
			break
		}

		time.Sleep(time.Duration(c.req_timeout))
	}

	c.updateStatus(controller.STATUS_ONLINE)

	c.startComm()

	c.JobOff("0")
	c.PSetSubscribe()
	//c.CurveSubscribe()
	c.SelectorSubscribe()
	c.ResultSubcribe()
	c.JobInfoSubscribe()
	//c.DataSubscribeCurve()
	// 启动发送
	go c.manage()

	go c.SolveOldResults()

	return nil
}

func (c *Controller) GetPSetList() ([]int, error) {
	var psets []int
	if c.Status() == controller.STATUS_OFFLINE {
		return psets, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0011_PSET_LIST_REPLY)
	c.Response.Add(MID_0011_PSET_LIST_REPLY, nil)

	psets_request := GeneratePackage(MID_0010_PSET_LIST_REQUEST, "001", "", DEFAULT_MSG_END)
	c.Write([]byte(psets_request))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0011_PSET_LIST_REPLY)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return psets, errors.New(controller.ERR_NOT_FOUND)
	}

	pset_list := reply.(PSetList)

	return pset_list.psets, nil
}

func (c *Controller) GetPSetDetail(pset int) (PSetDetail, error) {
	var obj_pset_detail PSetDetail

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_pset_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0013_PSET_DETAIL_REPLY)
	c.Response.Add(MID_0013_PSET_DETAIL_REPLY, nil)

	pset_detail := GeneratePackage(MID_0012_PSET_DETAIL_REQUEST, "002", fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)
	c.Write([]byte(pset_detail))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0013_PSET_DETAIL_REPLY)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return obj_pset_detail, errors.New(controller.ERR_NOT_FOUND)
	}

	return reply.(PSetDetail), nil
}

func (c *Controller) SolveOldResults() {
	if c.dbController.LastID == "0" {
		return
	}

	c.Response.Add(MID_0064_OLD_SUBSCRIBE, MID_0064_OLD_SUBSCRIBE)
	c.getOldResult(0)

	var last_result interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		last_result = c.Response.get(MID_0065_OLD_DATA)
		if last_result != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if last_result == nil {
		return
	}

	obj_last_result := last_result.(ResultData)

	if obj_last_result.TightingID != c.dbController.LastID {
		start_id, _ := strconv.ParseInt(c.dbController.LastID, 10, 64)
		end_id, _ := strconv.ParseInt(obj_last_result.TightingID, 10, 64)

		for i := start_id + 1; i <= end_id; i++ {
			c.getOldResult(i)
		}

	}
}

func (c *Controller) KeepAliveCount() int32 {
	return atomic.LoadInt32(&c.keepAliveCount)
}

func (c *Controller) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&c.keepAliveCount, i)
}

func (c *Controller) addKeepAliveCount() {
	atomic.AddInt32(&c.keepAliveCount, 1)
}

func (c *Controller) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keep_period))
}

func (c *Controller) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *Controller) Status() string {

	return c.StatusValue.Load().(string)
}

func (c *Controller) sendKeepalive() {
	if c.Status() == controller.STATUS_OFFLINE {
		return
	}

	keep_alive := GeneratePackage(MID_9999_ALIVE, DEFAULT_REV, "", DEFAULT_MSG_END)
	c.Write([]byte(keep_alive))
}

func (c *Controller) startComm() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	start := GeneratePackage(MID_0001_START, "003", "", DEFAULT_MSG_END)

	//c.Response.Add(MID_0001_START, "")

	c.Write([]byte(start))

	return nil
}

func (c *Controller) Write(buf []byte) {
	c.buffer <- buf
}

func (c *Controller) updateStatus(status string) {

	if status != c.Status() {

		c.StatusValue.Store(status)

		if status == controller.STATUS_OFFLINE {
			c.Close()

			// 断线重连
			go c.Connect()
		}

		// 将最新状态推送给hmi
		s := wsnotify.WSStatus{
			SN:     c.cfg.SN,
			Status: string(status),
		}

		msg, _ := json.Marshal(s)
		c.Srv.WS.WSSendControllerStatus(string(msg))

		c.Srv.diag.Debug(fmt.Sprintf("CVI3:%s %s\n", c.cfg.SN, status))

	}
}

func (c *Controller) Read(conn net.Conn) {
	defer conn.Close()

	len_header := LEN_HEADER
	rest := 0
	body := ""
	header_rest := 0

	var header_buffer string
	var header OpenProtocolHeader

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		c.updateKeepAliveCount(0)

		msg := string(buffer[0:n])

		off := 0 //循环前偏移为0
		for off < n {
			if rest == 0 {
				len_msg := n - off
				if len_msg < len_header-header_rest {
					//长度不够
					if header_rest == 0 {
						header_rest = len_header - len_msg
					} else {
						header_rest -= len_msg
					}
					header_buffer += msg[off : off+len_msg]
					break
				} else {
					//完整
					if header_rest == 0 {
						header_buffer = msg[off : off+len_header]
						off += len_header
					} else {
						header_buffer += msg[off : off+header_rest]
						off += header_rest
						header_rest = 0
					}
				}
				//fmt.Printf("header rest:%d, offset:%d, n %d, header : %s\n", header_rest, off, n, header_buffer)
				header.Deserialize(header_buffer)
				header_buffer = ""
				if n-off > header.LEN {
					//粘包
					body = msg[off : off+header.LEN]

					//c.Parse(body)
					pkg := handlerPkg{
						Header: header,
						Body:   body,
					}
					c.handlerBuf <- pkg
					off += header.LEN + 1
					rest = 0 //同样解析头

				} else {
					body = msg[off:n]
					rest = header.LEN - (n - off)
					break
				}
			} else {
				if n-off > rest {
					//粘包
					body += string(buffer[off : off+rest]) //已经是完整的包
					//p.Parse(body)
					pkg := handlerPkg{
						Header: header,
						Body:   body,
					}
					c.handlerBuf <- pkg
					off += rest + 1
					rest = 0 //进入解析头
				} else {
					body += string(buffer[off:n])
					rest -= n - off
					break
				}
			}
		}

	}
}

func (c *Controller) Close() error {

	closed := make(chan struct{})
	c.closing <- closed

	<-closed
	return c.w.Close()
}

func (c *Controller) manage() {
	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(c.keep_period):
			if c.Status() == controller.STATUS_OFFLINE {
				continue
			}
			if c.KeepAliveCount() >= MAX_KEEP_ALIVE_CHECK {
				go c.updateStatus(controller.STATUS_OFFLINE)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendKeepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}
		case v := <-c.buffer:
			for nextWriteThreshold.After(time.Now()) {
				time.Sleep(time.Microsecond * 100)
			}
			err := c.w.Write([]byte(v))
			if err != nil {
				c.Srv.diag.Error("Write data fail", err)
			} else {
				c.updateKeepAliveDeadLine()
			}
			nextWriteThreshold = time.Now().Add(c.req_timeout)
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *Controller) getOldResult(last_id int64) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_last_result := GeneratePackage(MID_0064_OLD_SUBSCRIBE, "006", fmt.Sprintf("%010d", last_id), DEFAULT_MSG_END)

	c.Write([]byte(s_last_result))

	return nil
}

func (c *Controller) pset(pset int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_pset := GeneratePackage(MID_0018_PSET, "001", fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)

	c.Write([]byte(s_pset))

	return nil
}

// 0: set 1: reset
func (c *Controller) JobOff(off string) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_off := GeneratePackage(MID_0130_JOB_OFF, "001", off, DEFAULT_MSG_END)

	c.Write([]byte(s_off))

	if off == "0" {
		c.Mode.Store(MODE_PSET)
	} else {
		c.Mode.Store(MODE_JOB)
	}

	return nil
}

func (c *Controller) jobSelect(job int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_job := GeneratePackage(MID_0038_JOB_SELECT, "002", fmt.Sprintf("%04d", job), DEFAULT_MSG_END)

	c.Write([]byte(s_job))

	return nil
}

func (c *Controller) IdentifierSet(str string) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, "001", str, DEFAULT_MSG_END)

	//c.Response.Add(MID_0018_PSET, "")

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) DataSubscribeCurve() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	cs := GeneratePackage(MID_0008_DATA_SUB, "001", "0900001350                             01001", DEFAULT_MSG_END)

	c.Write([]byte(cs))

	return nil
}

func (c *Controller) PSetSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0014_PSET_SUBSCRIBE, "000", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) SelectorSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0250_SELECTOR_SUBSCRIBE, "001", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) JobInfoSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0034_JOB_INFO_SUBSCRIBE, "003", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) ResultSubcribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0060_LAST_RESULT_SUBSCRIBE, "998", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) CurveSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_7408_LAST_CURVE_SUBSCRIBE, "000", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) PSet(pset int, workorder_id int64, result_id int64, count int, user_id int64, channel int) (uint32, error) {
	// 设定结果标识

	if c.Mode.Load().(string) != MODE_PSET {
		return 0, errors.New("current mode is not pset")
	}

	// 结果id-拧接次数-用户id
	err := c.IdentifierSet(fmt.Sprintf("%d-%d-%d", result_id, count, user_id))
	if err != nil {
		return 0, err
	}

	// 设定pset
	err = c.pset(pset)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (c *Controller) JobSet(id_info string, job int) error {

	if c.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	err := c.IdentifierSet(id_info)
	if err != nil {
		return err
	}

	err = c.jobSelect(job)
	if err != nil {
		return err
	}

	return nil
}
