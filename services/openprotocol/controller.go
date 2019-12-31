package openprotocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"net"
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

func NewController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher) (tightening_device.ITighteningController, error) {

	var c IOpenProtocolController
	switch deviceConfig.Model {
	case tightening_device.ModelDesoutterCvi3:
		c = &CVI3Controller{}
	case tightening_device.ModelDesoutterCvi2:
		c = &CVI2Controller{}
	case tightening_device.ModelDesoutterDeltaWrench:
		c = &WrenchController{}
	default:
		return nil, errors.New(fmt.Sprintf("Controller Model:%s Not Support", deviceConfig.Model))
	}

	controllerInstance := c.DefaultControllerGet()
	controllerInstance.initController(deviceConfig, d, service, dp)
	return controllerInstance, nil
}

type ControllerSubscribe func() error

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}

type TighteningController struct {
	sockClients          map[string]*socket_writer.SocketWriter
	deviceConf           *tightening_device.TighteningDeviceConfig
	keepAliveCount       int32
	keepPeriod           time.Duration
	reqTimeout           time.Duration
	getToolInfoPeriod    time.Duration
	Response             ResponseQueue
	ProtocolService      *Service
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
	dispatcherBus        Dispatcher
	dispatcherMap        map[string]dispatcherbus.DispatcherMap
	externalDispatches   map[string]*utils.Dispatcher

	requestChannel chan uint32
	sequence       *utils.Sequence

	handleRecvBuf []byte
	writeOffset   int
	device.BaseDevice
}

func (c *TighteningController) DefaultControllerGet() IOpenProtocolController {

	c.buffer = make(chan []byte, 1024)
	c.closing = make(chan chan struct{})
	c.keepPeriod = time.Duration(OpenProtocolDefaultKeepAlivePeriod)
	c.reqTimeout = time.Duration(OpenProtocolDefaultKeepAlivePeriod)
	c.getToolInfoPeriod = time.Duration(OpenProtocolDefaultGetTollInfoPeriod)
	c.protocol = tightening_device.TIGHTENING_OPENPROTOCOL
	c.dispatcherMap = map[string]dispatcherbus.DispatcherMap{}
	c.tempResultCurve = map[int]*tightening_device.TighteningCurve{}
	c.sockClients = map[string]*socket_writer.SocketWriter{}

	return c
}
func (c *TighteningController) createToolsByConfig() error {
	conf := c.deviceConf
	d := c.diag
	if conf == nil {
		return errors.New("Device Config Is Empty")
	}
	for _, v := range conf.Tools {
		tool := CreateTool(c, v, d)
		c.dispatcherMap[tool.SerialNumber()] = dispatcherbus.DispatcherMap{
			tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DISPATCH_RESULT): utils.CreateDispatchHandlerStruct(tool.OnResult),
			tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DISPATCH_CURVE):  utils.CreateDispatchHandlerStruct(tool.OnCurve),
		}
		c.AddChildren(v.SN, tool)
	}
	return nil
}

func (c *TighteningController) initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher) {
	c.BaseDevice = device.CreateBaseDevice(deviceConfig.Model, d, service)
	c.diag = d
	c.deviceConf = deviceConfig
	c.ProtocolService = service
	c.dispatcherBus = dp

	c.BaseDevice.Cfg = c.GetVendorModel()[IO_MODEL]

	c.initSubscribeInfos()

	if err := c.createToolsByConfig(); err != nil {
		d.Error("NewController createToolsByConfig Error", err)
	}
}

func (c *TighteningController) UpdateToolStatus(status string) {
	var ss []device.DeviceStatus
	for sn, v := range c.Children() {
		tool := v.(*TighteningTool)
		tool.UpdateStatus(status)
		ss = append(ss, device.DeviceStatus{
			Type:   tightening_device.TIGHTENING_DEVICE_TYPE_TOOL,
			SN:     sn,
			Status: status,
		})
	}
	if data, err := json.Marshal(ss); err == nil {
		c.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_DEVICE_STATUS, data)
	}
}

func (c *TighteningController) GetToolViaSerialNumber(toolSN string) (tightening_device.ITighteningTool, error) {
	return c.getToolViaSerialNumber(toolSN)
}

func (c *TighteningController) GetToolViaChannel(channel int) (tightening_device.ITighteningTool, error) {
	var serialNumber = ""
	for _, v := range c.deviceConf.Tools {
		if v.Channel == channel {
			serialNumber = v.SN
			break
		}
	}
	if serialNumber != "" {
		if tool, err := c.getToolViaSerialNumber(serialNumber); err != nil {
			return nil, errors.New("GetToolViaChannel Error, Tool Not Found")
		} else {
			return tool, nil
		}
	}

	return nil, errors.New("GetToolViaChannel Tool Not Found")
}

func (c *TighteningController) LoadController(controller *storage.Controllers) {
	c.dbController = controller
}

func (c *TighteningController) Inputs() string {
	return c.inputs
}

func (c *TighteningController) initSubscribeInfos() {
	c.controllerSubscribes = []ControllerSubscribe{
		//c.PSetSubscribe,
		c.ResultSubcribe,
		c.SelectorSubscribe,
		c.JobInfoSubscribe,
		c.IOInputSubscribe,
		//c.MultiSpindleResultSubscribe,
		c.VinSubscribe,
		c.AlarmSubcribe,
		c.CurveSubscribe,
	}
}

func (c *TighteningController) ProcessSubscribeControllerInfo() {
	for _, subscribe := range c.controllerSubscribes {
		err := subscribe()
		if err != nil {
			c.diag.Debug(fmt.Sprintf("OpenProtocol SubscribeControllerInfo Failed: %s", err.Error()))
		}
	}
}

func (c *TighteningController) ProcessRequest(mid string, noack string, station string, spindle string, data string) (interface{}, error) {
	rev, err := c.GetVendorMid(mid)
	if err != nil {
		return nil, err
	}

	if c.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	pkg := GeneratePackage(mid, rev, noack, station, spindle, data)

	seq := c.sequence.GetSequence()
	c.requestChannel <- seq
	c.Response.Add(seq, nil)

	c.Write([]byte(pkg))
	ctx, _ := context.WithTimeout(context.Background(), MAX_REPLY_TIME)
	reply := c.Response.Get(seq, ctx)

	if reply == nil {
		return nil, errors.New(tightening_device.TIGHTENING_ERR_TIMEOUT)
	}

	return reply, nil
}

func (c *TighteningController) CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float64, Angle []float64) {
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
			e := errors.New("Desoutter IProtocol Curve Raw Data 0xff不能单独出现")
			d.Error("CurveDataDecoding", e)
			// do nothing
		}
	}
	if writeOffset%6 != 0 {
		e := errors.New("Desoutter IProtocol Curve Raw Data Convert Fail")
		d.Error("CurveDataDecoding Fail", e)
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
	c.ProtocolService.diag.Debug(fmt.Sprintf("OpenProtocol Recv %s: %s%s\n", c.deviceConf.SN, pkg.Header.Serialize(), pkg.Body))

	handler, err := GetMidHandler(pkg.Header.MID)
	if err != nil {
		return err
	}

	return handler(c, pkg)
}

func (c *TighteningController) handleResult(result *tightening_device.TighteningResult) error {
	result.ControllerSN = c.deviceConf.SN
	tool, err := c.GetToolViaChannel(result.ChannelID)
	if err != nil {
		return err
	}

	toolSerialNumber := tool.SerialNumber()

	result.ToolSN = toolSerialNumber

	// 分发结果到工具进行处理
	c.dispatcherBus.Dispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DISPATCH_RESULT), result)

	return nil
}

// seq, count
func (c *TighteningController) calBatch(workorderID int64) (int, int) {
	result, err := c.ProtocolService.DB.FindTargetResultForJobManual(workorderID)
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
	for _, v := range c.deviceConf.Tools {
		_ = c.ProtocolService.DB.UpdateTool(&storage.Guns{
			Serial: v.SN,
			Mode:   "pset",
		})
	}
	for _, dd := range c.dispatcherMap {
		c.dispatcherBus.LaunchDispatchersByHandlerMap(dd)
	}

	// 启动处理
	go c.manage()
	go c.handleRecv()

	go c.Connect()

	return nil
}

func (c *TighteningController) Stop() error {
	for _, dd := range c.dispatcherMap {
		for name, v := range dd {
			c.dispatcherBus.Release(name, v.ID)
		}
	}

	return nil
}

func (c *TighteningController) getToolViaSerialNumber(toolSN string) (tightening_device.ITighteningTool, error) {
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
	for _, tool := range c.deviceConf.Tools {
		err := c.ProtocolService.DB.ClearToolResultAndCurve(tool.SN)
		if err != nil {
			c.diag.Error(fmt.Sprintf("Clear Tool: %s Result And Curve Failed", tool.SN), err)
		}
	}
}

func (c *TighteningController) initToolConnection() {
	for _, tool := range c.deviceConf.Tools {
		c.sockClients[tool.SN] = socket_writer.NewSocketWriter(c.deviceConf.Endpoint, c)
	}
}

func (c *TighteningController) CreateTransports() {
	for _, child := range c.Children() {
		if child == nil {
			continue
		}
		if writer := socket_writer.NewSocketWriter(c.deviceConf.Endpoint, c); writer != nil {
			c.sockClients[child.SerialNumber()] = writer
		}
	}
}

func (c *TighteningController) doConnect(doConnectToolSymbol string, forceCreate bool) error {
	if forceCreate {
		//重新创建transport
		if writer := socket_writer.NewSocketWriter(c.deviceConf.Endpoint, c); writer != nil {
			c.sockClients[doConnectToolSymbol] = writer
		} else {
			err := errors.Errorf("Create Transport For %s Fail", doConnectToolSymbol)
			c.diag.Error("doConnect", err)
			return err
		}
	}
	if writer, ok := c.sockClients[doConnectToolSymbol]; !ok {
		err := errors.Errorf("Can Not Found Transport For %s", doConnectToolSymbol)
		c.diag.Error("doConnect", err)
		return err
	} else {
		for {
			err := writer.Connect(DAIL_TIMEOUT)
			if err != nil {
				e := errors.Wrapf(err, "Connect To Tool %s", doConnectToolSymbol)
				c.diag.Error("doConnect", e)
			} else {
				c.diag.Debug(fmt.Sprintf("doConnect Connect To %s Success", doConnectToolSymbol))
				break
			}

			time.Sleep(time.Duration(c.reqTimeout))
		}
	}
	return nil
}

//todo: 连接创建后做一下系统必要操作
func (c *TighteningController) doConnectPart2(doConnectToolSymbol string) error {
	// 处理不完整的结果和曲线
	c.clearToolsResultAndCurve()

	c.handleStatus(device.BaseDeviceStatusOnline)

	return c.startComm()
}

func (c *TighteningController) CloseTransport() {
	var isCloseTransportSymbols []string
	for symbol, writer := range c.sockClients {
		if err := writer.Close(); err != nil {
			e := errors.Wrapf(err, "Close Transport For %s Error", symbol)
			c.diag.Error("CloseTransport", e)
		} else {
			isCloseTransportSymbols = append(isCloseTransportSymbols, symbol)
		}
	}
	for _, isCloseSymbol := range isCloseTransportSymbols {
		delete(c.sockClients, isCloseSymbol)
	}
}

func (c *TighteningController) Connect() error {
	c.UpdateStatus(device.BaseDeviceStatusOffline)
	c.handlerBuf = make(chan handlerPkg, 1024)
	c.writeOffset = 0
	c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence()
	c.Response = ResponseQueue{
		Results: map[interface{}]interface{}{},
		mtx:     sync.Mutex{},
	}

	c.CreateTransports() // always success

	for _, tool := range c.Children() {
		sn := tool.SerialNumber()
		if err := c.doConnect(sn, false); err == nil {
			c.doConnectPart2(sn)
		}
	}

	return nil
}

func (c *TighteningController) getTighteningCount() {
	for {
		select {
		case <-time.After(c.getToolInfoPeriod):
			rev, err := c.GetVendorMid(MID_0040_TOOL_INFO_REQUEST)
			if err == nil {
				continue
			}

			if c.Status() == device.BaseDeviceStatusOffline {
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
	rev, err := c.GetVendorMid(MID_0040_TOOL_INFO_REQUEST)
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
	//	c.handleStatus(controller.BaseDeviceStatusOnline)
	//} else {
	//	c.handleStatus(controller.BaseDeviceStatusOffline)
	//}

	return nil
}

func (c *TighteningController) handlerOldResults() error {
	return nil
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
	if c.Status() == device.BaseDeviceStatusOffline {
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
	c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.deviceConf.SN, string(buf)))
	c.buffer <- buf
}

func (c *TighteningController) handleStatus(status string) {

	if status != c.Status() {
		c.diag.Debug(fmt.Sprintf("OpenProtocol handleStatus %s:%s %s\n", c.Model(), c.deviceConf.SN, status))

		c.UpdateStatus(status)

		if status == device.BaseDeviceStatusOffline {

			// 断线重连
			go c.Connect()
		}
		ss := []device.DeviceStatus{
			{
				Type:   tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER,
				SN:     c.deviceConf.SN,
				Status: status,
				Config: c.Config(),
			},
		}
		// 分发控制器状态 -> tightening device
		c.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_DEVICE_STATUS, ss)
	}
}

func (c *TighteningController) Read(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.diag.Error("Controller Close Error", err)
		}
	}()

	buffer := make([]byte, c.ProtocolService.config().ReadBufferSize)

	for {
		conn.SetReadDeadline(time.Now().Add(c.keepPeriod * MAX_KEEP_ALIVE_CHECK).Add(1 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			c.ProtocolService.diag.Error("read failed", err)
			c.handleStatus(device.BaseDeviceStatusOffline)
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
	c.CloseTransport()

	return nil
}

func (c *TighteningController) handlePackageOPPayload(src []byte, data []byte) error {
	msg := append(src, data...)

	//c.diag.Debug(fmt.Sprintf("%s op target buf: %s", c.deviceConf.SerialNumber, string(msg)))

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
		return errors.New(fmt.Sprintf("Body Len Err: %s", string(msg)))
	}

	return nil
}

func (c *TighteningController) getDefaultTransportClient() *socket_writer.SocketWriter {
	for _, sw := range c.sockClients {
		return sw
	}
	return nil
}

func (c *TighteningController) getTransportClientBySymbol(symbol string) *socket_writer.SocketWriter {
	if sw, ok := c.sockClients[symbol]; !ok {
		err := errors.Errorf("Can Not Found Transport For %s", symbol)
		c.diag.Error("getTransportClientBySymbol", err)
		return nil
	} else {
		return sw
	}
}

func (c *TighteningController) handleRecv() {
	c.handleRecvBuf = make([]byte, c.ProtocolService.config().ReadBufferSize)
	c.receiveBuf = make(chan []byte, 65535)
	lenBuf := len(c.handleRecvBuf)

	for {
		select {
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
					restBuf := buf[readOffset:]
					if c.writeOffset+len(restBuf) > lenBuf {
						c.diag.Error("full", errors.New("full"))
						break
					}

					copy(c.handleRecvBuf[c.writeOffset:c.writeOffset+len(restBuf)], restBuf)
					c.writeOffset += len(restBuf)
					break
				} else {
					// 找到结束字符，结合缓冲进行处理
					err := c.handlePackageOPPayload(c.handleRecvBuf[0:c.writeOffset], buf[readOffset:readOffset+index])
					if err != nil {
						//数据需要丢弃
						c.diag.Error("msg", err)
					}

					c.writeOffset = 0
					readOffset += index + 1
				}
			}
		}
	}
}

func (c *TighteningController) manage() {

	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(c.keepPeriod):
			if c.Status() == device.BaseDeviceStatusOffline {
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

			//fixme 获取默认的transport，未来要根据发送的结构体获取真实的transport
			sw := c.getDefaultTransportClient()

			if sw == nil {
				continue
			}

			err := sw.Write([]byte(v))
			if err != nil {
				c.ProtocolService.diag.Error("Write Data Fail", err)
			} else {
				c.updateKeepAliveDeadLine()
			}
			nextWriteThreshold = time.Now().Add(c.reqTimeout)

		case stopDone := <-c.closing:
			c.diag.Debug("manage exit")
			close(stopDone)
			return //退出manage协程

		case pkg := <-c.handlerBuf:
			err := c.HandleMsg(&pkg)
			if err != nil {
				c.diag.Error("Open IProtocol HandleMsg Fail", err)
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
	return c.deviceConf.Model
}

func (c *TighteningController) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER
}

func (c *TighteningController) GetDispatch(name string) *utils.Dispatcher {
	return c.externalDispatches[name]
}

func (c *TighteningController) GetVendorModel() map[string]interface{} {
	return nil
}

func (c *TighteningController) GetVendorMid(mid string) (string, error) {
	rev, exist := c.GetVendorModel()[mid]
	if !exist {
		return "", errors.New(fmt.Sprintf("MID %s Not Supported", mid))
	}

	return rev.(string), nil
}
