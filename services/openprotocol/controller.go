package openprotocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DAIL_TIMEOUT         = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK = 3

	REPLY_TIMEOUT  = time.Duration(100 * time.Millisecond)
	MAX_REPLY_TIME = time.Duration(2000 * time.Millisecond)

	DEFAULT_TCP_KEEPALIVE = time.Duration(5 * time.Second)
)

type ToolDispatch struct {
	resultDispatch *utils.Dispatcher
	curveDispatch  *utils.Dispatcher
}

type ControllerSubscribe func() error

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}

type handlerPkg_curve struct {
	Header *OpenProtocolHeader
	Body   []byte
}

type TighteningController struct {
	w                    *socket_writer.SocketWriter
	cfg                  *tightening_device.TighteningDeviceConfig
	keepAliveCount       int32
	keepPeriod           time.Duration
	reqTimeout           time.Duration
	getToolInfoPeriod    time.Duration
	Response             ResponseQueue
	Srv                  *Service
	dbController         *storage.Controllers
	buffer               chan []byte
	closing              chan chan struct{}
	handlerBuf           chan handlerPkg
	keepaliveDeadLine    atomic.Value
	protocol             string
	inputs               string
	diag                 Diagnostic
	tempResultCurve      map[int]*tightening_device.TighteningCurve
	mtxResult            sync.Mutex
	model                string
	receiveBuf           chan []byte
	controllerSubscribes []ControllerSubscribe

	toolDispatches     map[string]*ToolDispatch
	externalDispatches map[string]*utils.Dispatcher

	requestChannel chan uint32
	sequence       *utils.Sequence

	device.BaseDevice
}

// TODO: 如果工具序列号没有配置，则通过探测加入设备列表。
func NewController(protocolConfig *Config, deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic) *TighteningController {

	c := TighteningController{
		diag:              d,
		buffer:            make(chan []byte, 1024),
		closing:           make(chan chan struct{}),
		keepPeriod:        time.Duration(protocolConfig.KeepAlivePeriod),
		reqTimeout:        time.Duration(protocolConfig.ReqTimeout),
		getToolInfoPeriod: time.Duration(protocolConfig.GetToolInfoPeriod),
		protocol:          controller.OPENPROTOCOL,
		mtxResult:         sync.Mutex{},
		cfg:               deviceConfig,
		BaseDevice:        device.CreateBaseDevice(),
		tempResultCurve:   map[int]*tightening_device.TighteningCurve{},

		toolDispatches: map[string]*ToolDispatch{},
	}

	c.controllerSubscribes = []ControllerSubscribe{
		//c.PSetSubscribe,
		c.ResultSubcribe,
		c.SelectorSubscribe,
		c.JobInfoSubscribe,
		c.IOInputSubscribe,
		c.MultiSpindleResultSubscribe,
		c.VinSubscribe,
		c.AlarmSubcribe,
		c.CurveSubscribe,
	}

	for _, v := range deviceConfig.Tools {
		tool := NewTool(&c, v, d)
		c.toolDispatches[v.SN] = &ToolDispatch{
			resultDispatch: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			curveDispatch:  utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		}
		c.toolDispatches[v.SN].resultDispatch.Regist(tool.OnResult)
		c.toolDispatches[v.SN].curveDispatch.Regist(tool.OnCurve)

		c.AddChildren(v.SN, tool)
	}

	c.externalDispatches = map[string]*utils.Dispatcher{
		tightening_device.DISPATCH_RESULT:            utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		tightening_device.DISPATCH_IO:                utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		tightening_device.DISPATCH_CONTROLLER_STATUS: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		tightening_device.DISPATCH_TOOL_STATUS:       utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}

	return &c
}

//func (c *TighteningController) UpdateToolStatus(status string) {
//	s := c.toolStatus.Load().(string)
//	if s != status {
//		c.toolStatus.Store(status)
//
//		// 推送工具状态
//		//ts := wsnotify.WSToolStatus{
//		//	ToolSN: c.cfg.Tools[0].SerialNO,
//		//	Status: status,
//		//}
//		//
//		//str, _ := json.Marshal(ts)
//		//c.Srv.WS.WSSend(wsnotify.WS_EVETN_TOOL, string(str))
//	}
//}

func (c *TighteningController) findToolSNByChannel(channel int) (string, error) {
	for _, v := range c.cfg.Tools {
		if v.Channel == channel {
			return v.SN, nil
		}
	}

	return "", errors.New("Tool Not Found")
}

func (c *TighteningController) LoadController(controller *storage.Controllers) {
	c.dbController = controller
}

func (c *TighteningController) Inputs() string {
	return c.inputs
}

func (c *TighteningController) Subscribe() {
	for _, subscribe := range c.controllerSubscribes {
		err := subscribe()
		if err != nil {
			c.diag.Debug(fmt.Sprintf("OpenProtocol Subscribe Failed: %s", err.Error()))
		}
	}
}

func (c *TighteningController) ProcessRequest(mid string, noack string, station string, spindle string, data string) (interface{}, error) {
	rev, err := GetVendorMid(c.Model(), mid)
	if err != nil {
		return nil, err
	}

	if c.Status() == device.STATUS_OFFLINE {
		return nil, errors.New(device.STATUS_OFFLINE)
	}

	pkg := GeneratePackage(mid, rev, noack, station, spindle, data)

	seq := c.sequence.GetSequence()
	c.requestChannel <- seq
	c.Response.Add(seq, nil)

	c.Write([]byte(pkg))
	ctx, _ := context.WithTimeout(context.Background(), MAX_REPLY_TIME)
	reply := c.Response.Get(seq, ctx)

	if reply == nil {
		return nil, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	return reply, nil
}

func DataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float64, Angle []float64) {
	lenO := len(original)
	data := make([]byte, lenO, lenO) // 最大只会这些数据
	writeOffset := 0
	step := 1
	for i := 0; i < lenO; i += step {
		step = 1 // 初始化step
		if original[i] != 0xff {
			data[writeOffset] = original[i]
			writeOffset += 1
			continue
		}
		//检测到0xff
		if i+1 == lenO {
			//下一个字节大于整体长度, 最后一个字节了
			data[writeOffset] = original[i]
			writeOffset += 1
			break
		}
		switch original[i+1] {
		case 0xff:
			data[writeOffset] = 0xff
			writeOffset += 1
			step = 2 //跳过这个字节
		case 0xfe:
			data[writeOffset] = 0x00
			writeOffset += 1
			step = 2 //跳过这个字节
		default:
			e := errors.New("Desoutter Protocol Curve Raw Data 0xff不能单独出现")
			d.Error("DataDecoding", e)
			// do nothing
		}
	}
	if writeOffset%6 != 0 {
		e := errors.New("Desoutter Protocol Curve Raw Data Convert Fail")
		d.Error("DataDecoding Fail", e)
		return
	}
	// 所有位减1
	for i := 0; i < writeOffset; i++ {
		if data[i] == 0x00 {
			data[i] = 0xff
		} else {
			data[i] = data[i] - 1
		}
	}

	for i := 0; i < writeOffset; i += 6 {
		a := binary.LittleEndian.Uint16(data[i : i+2])
		b := binary.LittleEndian.Uint32(data[i+2 : i+6])
		Torque = append(Torque, float64(a)*torqueCoefficient)
		Angle = append(Angle, float64(b)*angleCoefficient)
	}
	return
}

func (c *TighteningController) HandleMsg(pkg *handlerPkg) error {
	c.Srv.diag.Debug(fmt.Sprintf("OpenProtocol Recv %s: %s%s\n", c.cfg.SN, pkg.Header.Serialize(), pkg.Body))

	handler, err := GetMidHandler(pkg.Header.MID)
	if err != nil {
		return err
	}

	return handler(c, pkg)
}

func (c *TighteningController) handleResult(result *tightening_device.TighteningResult) error {

	//if utils.ArrayContains(c.Srv.config().SkipJobs, result_data.JobID) {
	//	return nil
	//}
	//
	//if c.dbController != nil {
	//	c.Srv.DB.UpdateTightning(c.dbController.Id, result_data.TightingID)
	//}
	//
	//controllerResult := controller.ControllerResult{}
	//controllerResult.NeedPushHmi = true
	//
	//if c.Model() == tightening_device.ModelDesoutterDeltaWrench {
	//	controllerResult.ToolSN = c.cfg.SN
	//} else {
	//	controllerResult.ToolSN = result_data.ToolSerialNumber
	//}
	//
	//gun, err := c.Srv.DB.GetGun(controllerResult.ToolSN)
	//if err != nil {
	//	c.Srv.diag.Error("get gun failed", err)
	//	return err
	//}
	//
	//psetTrace := tightening_device.PSetSet{}
	//_ = json.Unmarshal([]byte(gun.Trace), &psetTrace)
	//
	//controllerResult.TighteningID = result_data.TightingID
	//controllerResult.Count = psetTrace.Count
	//controllerResult.Batch = fmt.Sprintf("%d/%d", psetTrace.Sequence, psetTrace.Total)
	//
	//dat_kvs := strings.Split(result_data.TimeStamp, ":")
	//controllerResult.Dat = fmt.Sprintf("%s %s:%s:%s", dat_kvs[0], dat_kvs[1], dat_kvs[2], dat_kvs[3])
	//
	//controllerResult.PSet = result_data.PSetID
	//controllerResult.Controller_SN = c.cfg.SN
	//if result_data.TighteningStatus == "0" {
	//	controllerResult.Result = storage.RESULT_NOK
	//} else {
	//	controllerResult.Result = storage.RESULT_OK
	//}
	//
	//controllerResult.ResultValue.Mi = result_data.Torque / 100
	//controllerResult.ResultValue.Wi = result_data.Angle
	////controllerResult.ResultValue.Ti = result_data.
	//
	//switch result_data.Strategy {
	//case "01":
	//	controllerResult.PSetDefine.Strategy = controller.STRATEGY_AW
	//
	//case "02":
	//	controllerResult.PSetDefine.Strategy = controller.STRATEGY_AW
	//
	//case "03":
	//	controllerResult.PSetDefine.Strategy = controller.STRATEGY_ADW
	//
	//case "04":
	//	controllerResult.PSetDefine.Strategy = controller.STRATEGY_AD
	//}
	//
	//if result_data.ResultType == "02" {
	//	controllerResult.Result = storage.RESULT_LSN
	//	controllerResult.NeedPushHmi = true
	//	controllerResult.PSetDefine.Strategy = controller.STRATEGY_LN
	//}
	//
	//controllerResult.PSetDefine.Mp = result_data.TorqueMax / 100
	//controllerResult.PSetDefine.Mm = result_data.TorqueMin / 100
	//controllerResult.PSetDefine.Ma = result_data.TorqueFinalTarget / 100
	//
	//controllerResult.PSetDefine.Wp = result_data.AngleMax
	//controllerResult.PSetDefine.Wm = result_data.AngleMin
	//controllerResult.PSetDefine.Wa = result_data.FinalAngleTarget
	//
	//controllerResult.ExceptionReason = result_data.TighteningErrorStatus

	//targetID := result_data.VIN
	//switch c.Srv.config().DataIndex {
	//case 1:
	//	targetID = result_data.ID2
	//case 2:
	//	targetID = result_data.ID3
	//case 3:
	//	targetID = result_data.ID4
	//}

	//controllerResult.Workorder_ID = psetTrace.WorkorderID
	//controllerResult.NeedPushAiis = true

	//controllerResult.Seq, controllerResult.Count = c.calBatch(controllerResult.Workorder_ID)

	//c.result = &controllerResult

	//c.Srv.Parent.Handlers.Handle(&controllerResult, c.result_CURVE)

	//c.updateResult(&controllerResult, nil, result_data.ChannelID)
	//	c.handleResultandClear(result_data.ChannelID)

	//tighteningResult := result_data.ToTighteningResult()
	result.ControllerSN = c.cfg.SN
	toolSN, err := c.findToolSNByChannel(result.ChannelID)
	if err != nil {
		return err
	}

	result.ToolSN = toolSN

	// 分发结果到工具进行处理
	c.toolDispatches[toolSN].resultDispatch.Dispatch(result)

	return nil
}

// seq, count
func (c *TighteningController) calBatch(workorderID int64) (int, int) {
	result, err := c.Srv.DB.FindTargetResultForJobManual(workorderID)
	if err != nil {
		return 1, 1
	}

	if result.Result == storage.RESULT_OK {
		return result.GroupSeq + 1, 1
	} else {
		return result.GroupSeq, result.Count + 1
	}
}

func (c *TighteningController) Start() error {
	for _, v := range c.cfg.Tools {
		_ = c.Srv.DB.UpdateTool(&storage.Guns{
			Serial: v.SN,
			Mode:   "pset",
		})
	}

	for _, dispatch := range c.toolDispatches {
		dispatch.resultDispatch.Start()
		dispatch.curveDispatch.Start()
	}

	for _, dispatch := range c.externalDispatches {
		dispatch.Start()
	}

	c.w = socket_writer.NewSocketWriter(c.cfg.Endpoint, c)

	// 启动处理
	go c.manage()

	go c.Connect()

	return nil
}

func (c *TighteningController) Stop() error {
	for _, dispatch := range c.toolDispatches {
		dispatch.resultDispatch.Release()
		dispatch.curveDispatch.Release()
	}

	return nil
}

func (c *TighteningController) GetTool(toolSN string) (tightening_device.ITighteningTool, error) {
	tool, exist := c.Children()[toolSN]
	if !exist {
		return nil, errors.New("Not Found")
	}

	return tool.(tightening_device.ITighteningTool), nil
}

func (c *TighteningController) SetOutput(outputs []tightening_device.ControllerOutput) error {
	if len(outputs) == 0 {
		return errors.New("Output List Is Required")
	}

	strIo := ""
	for i := 0; i < 10; i++ {
		io, err := c.findIOByNo(i, outputs)
		if err != nil {
			strIo += "3"
		} else {
			switch io.Status {
			case tightening_device.IO_STATUS_OFF:
				strIo += "0"

			case tightening_device.IO_STATUS_ON:
				strIo += "1"

			case tightening_device.IO_STATUS_FLASHING:
				strIo += "2"
			}
		}
	}

	reply, err := c.ProcessRequest(MID_0200_CONTROLLER_RELAYS, "", "", "", strIo)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

func (c *TighteningController) Protocol() string {
	return c.protocol
}

func (c *TighteningController) clearToolsResultAndCurve() {
	for _, tool := range c.cfg.Tools {
		err := c.Srv.DB.ClearToolResultAndCurve(tool.SN)
		if err != nil {
			c.diag.Error(fmt.Sprintf("Clear Tool: %s Result And Curve Failed", tool.SN), err)
		}
	}
}

func (c *TighteningController) Connect() error {
	c.UpdateStatus(device.STATUS_OFFLINE)
	c.handlerBuf = make(chan handlerPkg, 1024)
	c.receiveBuf = make(chan []byte, 65535)
	c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence()
	c.Response = ResponseQueue{
		Results: map[interface{}]interface{}{},
		mtx:     sync.Mutex{},
	}

	for {
		err := c.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.Srv.diag.Error("Connect Err", err)
		} else {
			break
		}

		time.Sleep(time.Duration(c.reqTimeout))
	}

	// 处理不完整的结果和曲线
	c.clearToolsResultAndCurve()

	c.handleStatus(device.STATUS_ONLINE)

	c.startComm()

	return nil
}

func (c *TighteningController) getTighteningCount() {
	for {
		select {
		case <-time.After(c.getToolInfoPeriod):
			rev, err := GetVendorMid(c.Model(), MID_0040_TOOL_INFO_REQUEST)
			if err == nil {
				continue
			}

			if c.Status() == device.STATUS_OFFLINE {
				continue
			}
			req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", "", "", "")
			c.Write([]byte(req))
		case stopDone := <-c.closing:
			close(stopDone)
			return
		}
	}
}

func (c *TighteningController) ToolInfoReq() error {
	rev, err := GetVendorMid(c.Model(), MID_0040_TOOL_INFO_REQUEST)
	if err != nil {
		return err
	}

	req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", "", "", "")
	c.Write([]byte(req))

	//var reply interface{} = nil

	//for i := 0; i < MAX_REPLY_COUNT; i++ {
	//	reply = c.Response.get(MID_0040_TOOL_INFO_REQUEST)
	//	if reply != nil {
	//		break
	//	}
	//
	//	time.Sleep(REPLY_TIMEOUT)
	//}
	//
	//if reply != nil {
	//	ti := reply.(ToolInfo)
	//
	//	c.tighteningDevice.AddDevice(ti.ControllerSN, c)
	//	c.tighteningDevice.AddDevice(ti.ToolSN, c)
	//	c.handleStatus(controller.STATUS_ONLINE)
	//} else {
	//	c.handleStatus(controller.STATUS_OFFLINE)
	//}

	return nil
}

func (c *TighteningController) SolveOldResults() {

	if c.dbController == nil || c.dbController.LastID == "0" {
		return
	}

	lastResult, err := c.getOldResult(0)
	if err != nil {
		return
	}

	if lastResult.TighteningID != c.dbController.LastID {
		startId, _ := strconv.ParseInt(c.dbController.LastID, 10, 64)
		endId, _ := strconv.ParseInt(lastResult.TighteningID, 10, 64)

		for i := startId + 1; i <= endId; i++ {
			result, err := c.getOldResult(i)
			if err != nil {
				c.diag.Error(fmt.Sprintf("Get Old Result Failed: %d", i), err)
			} else {
				c.handleResult(result)
			}
		}
	}
}

func (c *TighteningController) KeepAliveCount() int32 {
	return atomic.LoadInt32(&c.keepAliveCount)
}

func (c *TighteningController) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&c.keepAliveCount, i)
}

func (c *TighteningController) addKeepAliveCount() {
	atomic.AddInt32(&c.keepAliveCount, 1)
}

func (c *TighteningController) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keepPeriod))
}

func (c *TighteningController) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *TighteningController) sendKeepalive() {
	if c.Status() == device.STATUS_OFFLINE {
		return
	}

	keepAlive := GeneratePackage(MID_9999_ALIVE, DEFAULT_REV, "1", "", "", "")
	c.Write([]byte(keepAlive))
}

func (c *TighteningController) startComm() error {
	reply, err := c.ProcessRequest(MID_0001_START, "", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

func (c *TighteningController) Write(buf []byte) {
	c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.cfg.SN, string(buf)))
	c.buffer <- buf
}

func (c *TighteningController) handleStatus(status string) {

	if status != c.Status() {
		c.diag.Debug(fmt.Sprintf("OpenProtocol handleStatus %s:%s %s\n", c.Model(), c.cfg.SN, status))

		c.UpdateStatus(status)

		if status == device.STATUS_OFFLINE {

			// 断线重连
			go c.Connect()
		}

		// 分发控制器状态
		c.GetDispatch(tightening_device.DISPATCH_CONTROLLER_STATUS).Dispatch(&tightening_device.TighteningControllerStatus{
			ControllerSN: c.cfg.SN,
			Status:       status,
		})
	}
}

func (c *TighteningController) Read(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		conn.SetReadDeadline(time.Now().Add(c.keepPeriod * MAX_KEEP_ALIVE_CHECK).Add(1 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			c.Srv.diag.Error("read failed", err)
			c.handleStatus(device.STATUS_OFFLINE)
			break
		}

		c.updateKeepAliveCount(0)
		c.receiveBuf <- buffer[0:n]
	}
}

func (c *TighteningController) Parse(msg string) {
	header := msg[0:LEN_HEADER]
	headerObj := OpenProtocolHeader{}
	headerObj.Deserialize(header)
	body := msg[LEN_HEADER:]

	pkg := handlerPkg{
		Header: headerObj,
		Body:   body,
	}

	c.handlerBuf <- pkg
}

func (c *TighteningController) Close() error {

	for i := 0; i < 2; i++ {
		//两个协程需要关闭
		closed := make(chan struct{})
		c.closing <- closed

		<-closed
	}

	return c.w.Close()
}

func (c *TighteningController) handlePackageOPPayload(src []byte, data []byte) error {
	msg := append(src, data...)

	//c.diag.Debug(fmt.Sprintf("%s op target buf: %s", c.cfg.SN, string(msg)))

	lenMsg := len(msg)

	// 如果头的长度不够
	if lenMsg < LEN_HEADER {
		return errors.New("Head Is Error")
	}

	header := OpenProtocolHeader{}
	header.Deserialize(string(msg[0:LEN_HEADER]))

	// 如果body的长度匹配
	if header.LEN == lenMsg-LEN_HEADER {
		pkg := handlerPkg{
			Header: header,
			Body:   string(msg[LEN_HEADER : LEN_HEADER+header.LEN]),
		}

		c.handlerBuf <- pkg
	} else {
		return errors.New("body len err")
	}

	return nil
}

func (c *TighteningController) manage() {

	lenBuf := 65535
	handleBuf := make([]byte, lenBuf)
	writeOffset := 0

	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(c.keepPeriod):
			if c.Status() == device.STATUS_OFFLINE {
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
				c.Srv.diag.Error("Write Data Fail", err)
			} else {
				c.updateKeepAliveDeadLine()
			}
			nextWriteThreshold = time.Now().Add(c.reqTimeout)
		case stopDone := <-c.closing:
			c.diag.Debug("manage exit")
			close(stopDone)
			return //退出manage协程

		case buf := <-c.receiveBuf:
			// 处理接收缓冲
			var readOffset = 0

			for {
				if readOffset >= len(buf) {
					break
				}
				index := bytes.IndexByte(buf[readOffset:], OP_TERMINAL)
				if index == -1 {
					// 没有结束字符,放入缓冲等待后续处理
					//c.diag.Debug("Index Is Empty")
					restBuf := buf[readOffset:]
					if writeOffset+len(restBuf) > lenBuf {
						c.diag.Error("full", errors.New("full"))
						break
					}

					copy(handleBuf[writeOffset:writeOffset+len(restBuf)], restBuf)
					writeOffset += len(restBuf)
					break
				} else {
					// 找到结束字符，结合缓冲进行处理
					err := c.handlePackageOPPayload(handleBuf[0:writeOffset], buf[readOffset:readOffset+index])
					if err != nil {
						//数据需要丢弃
						c.diag.Error("msg", err)
					}

					writeOffset = 0
					readOffset += index + 1
				}
			}

		case pkg := <-c.handlerBuf:
			err := c.HandleMsg(&pkg)
			if err != nil {
				c.diag.Error("Open Protocol HandleMsg Fail", err)
			}
		}
	}
}

func (c *TighteningController) getOldResult(last_id int64) (*tightening_device.TighteningResult, error) {
	reply, err := c.ProcessRequest(MID_0064_OLD_SUBSCRIBE, "", "", "", fmt.Sprintf("%010d", last_id))
	if err != nil {
		return nil, err
	}

	return reply.(*tightening_device.TighteningResult), nil
}

func (c *TighteningController) PSetSubscribe() error {

	reply, err := c.ProcessRequest(MID_0014_PSET_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0060_LAST_RESULT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) SelectorSubscribe() error {
	reply, err := c.ProcessRequest(MID_0250_SELECTOR_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0250_SELECTOR_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) JobInfoSubscribe() error {

	reply, err := c.ProcessRequest(MID_0034_JOB_INFO_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0034_JOB_INFO_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) IOInputSubscribe() error {
	reply, err := c.ProcessRequest(MID_0210_INPUT_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0210_INPUT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) MultiSpindleResultSubscribe() error {

	reply, err := c.ProcessRequest(MID_0100_MULTI_SPINDLE_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0100_MULTI_SPINDLE_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) VinSubscribe() error {
	reply, err := c.ProcessRequest(MID_0051_VIN_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0051_VIN_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) ResultSubcribe() error {

	reply, err := c.ProcessRequest(MID_0060_LAST_RESULT_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0060_LAST_RESULT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) AlarmSubcribe() error {

	reply, err := c.ProcessRequest(MID_0070_ALARM_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0070_ALARM_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) CurveSubscribe() error {
	reply, err := c.ProcessRequest(MID_7408_LAST_CURVE_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_7408_LAST_CURVE_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) findIOByNo(no int, outputs []tightening_device.ControllerOutput) (tightening_device.ControllerOutput, error) {
	for _, v := range outputs {
		if no == v.OutputNo {
			return v, nil
		}
	}

	return tightening_device.ControllerOutput{}, errors.New("Not Found")
}

func (c *TighteningController) Model() string {
	return c.cfg.Model
}

func (c *TighteningController) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER
}

func (c *TighteningController) GetDispatch(name string) *utils.Dispatcher {
	return c.externalDispatches[name]
}
