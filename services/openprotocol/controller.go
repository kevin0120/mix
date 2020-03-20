package openprotocol

import (
	"encoding/binary"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type ControllerSubscribe func(string) error

type handlerPkg struct {
	SN     string
	Header OpenProtocolHeader
	Body   string
	Seq    uint32
}

type TighteningController struct {
	device.BaseDevice

	sockClients           map[string]*clientContext
	deviceConf            *tightening_device.TighteningDeviceConfig
	ProtocolService       *Service
	inputs                string
	diag                  Diagnostic
	controllerSubscribes  []ControllerSubscribe
	hasSubscribedSelector atomic.Bool
	hasSubscribedVin      atomic.Bool
	hasSubscribedIO       atomic.Bool
	dispatcherBus         Dispatcher
	dispatcherMap         map[string]dispatcherbus.DispatcherMap
	isGlobalConn          bool
	mtxStatus             sync.Mutex

	instance IOpenProtocolController
	ioNotify io.IONotify
	opened   atomic.Bool

	internalIO tightening_device.ITighteningIO
}

func (c *TighteningController) CreateIO() tightening_device.ITighteningIO {
	c.internalIO = NewTighteningIO(c)
	return c.internalIO
}

func (c *TighteningController) createToolsByConfig() error {
	conf := c.deviceConf
	d := c.diag
	if conf == nil {
		return errors.New("Device Config Is Empty")
	}
	for _, v := range conf.Tools {
		tool := NewTool(c, v, d)
		c.dispatcherMap[tool.SerialNumber()] = dispatcherbus.DispatcherMap{
			tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherResult): utils.CreateDispatchHandlerStruct(tool.onResult),
			tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherCurve):  utils.CreateDispatchHandlerStruct(tool.onCurve),
		}
		c.AddChildren(v.SN, tool)
	}
	return nil
}

func (c *TighteningController) initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher) {

	c.opened.Store(false)
	c.dispatcherMap = map[string]dispatcherbus.DispatcherMap{}
	c.sockClients = map[string]*clientContext{}
	c.isGlobalConn = false
	c.hasSubscribedIO.Store(false)
	c.hasSubscribedSelector.Store(false)
	c.hasSubscribedVin.Store(false)
	c.BaseDevice = device.CreateBaseDevice(deviceConfig.Model, d, service)
	c.diag = d
	c.deviceConf = deviceConfig
	c.ProtocolService = service
	c.dispatcherBus = dp
	c.mtxStatus = sync.Mutex{}

	c.BaseDevice.Cfg = c.getInstance().GetVendorModel()[IoModel]

	c.initSubscribeInfos()

	if err := c.createToolsByConfig(); err != nil {
		d.Error("newController createToolsByConfig Error", err)
	}

	c.initClients(deviceConfig, d)
}

func (c *TighteningController) initClients(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic) {

	for _, v := range deviceConfig.Tools {
		endpoint := v.Endpoint
		sn := v.SN
		if deviceConfig.Endpoint != "" {
			// 全局链接
			c.isGlobalConn = true
			endpoint = deviceConfig.Endpoint
			sn = deviceConfig.SN
		} else {
			// 每个工具独立链接
			c.isGlobalConn = false
		}

		client := newClientContext(endpoint, d, c.getInstance().(IClientHandler), sn, c.getInstance().OpenProtocolParams())
		c.sockClients[sn] = client

		if c.isGlobalConn {
			break
		}
	}
}

func (c *TighteningController) getClient(sn string) *clientContext {
	if c.isGlobalConn {
		return c.getDefaultTransportClient()
	}

	return c.getTransportClientBySymbol(sn)
}

func (c *TighteningController) UpdateToolStatus(sn string, status string) {
	tool, err := c.getToolViaSerialNumber(sn)
	if err != nil {
		return
	}

	tool.(*TighteningTool).UpdateStatus(status)
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

func (c *TighteningController) initSubscribeInfos() {
	c.controllerSubscribes = []ControllerSubscribe{
		c.ResultSubcribe,
		c.SelectorSubscribe,
		c.JobInfoSubscribe,
		c.IOInputSubscribe,
		c.VinSubscribe,
		c.AlarmSubcribe,
		c.CurveSubscribe,
	}
}

func (c *TighteningController) processSubscribeControllerInfo(sn string) {
	for _, subscribe := range c.controllerSubscribes {
		err := subscribe(sn)
		if err != nil {
			c.diag.Debug(fmt.Sprintf("SN: %s OpenProtocol SubscribeControllerInfo Failed: %s", sn, err.Error()))
		}
	}
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

func (c *TighteningController) handleMsg(pkg *handlerPkg, context *clientContext) error {
	c.ProtocolService.diag.Debug(fmt.Sprintf("OpenProtocol Recv %s: %s%s\n", pkg.SN, pkg.Header.Serialize(), pkg.Body))

	handler, err := GetMidHandler(pkg.Header.MID)
	if err != nil {
		return err
	}

	if utils.ArrayContains(midsWithSeq, pkg.Header.MID) {
		pkg.Seq = <-context.requestChannel
	}

	return handler(c, pkg)
}

func (c *TighteningController) handleResult(result tightening_device.TighteningResult) error {
	result.ControllerSN = c.deviceConf.SN
	tool, err := c.getInstance().GetToolViaChannel(result.ChannelID)
	if err != nil {
		return err
	}

	toolSerialNumber := tool.SerialNumber()

	result.ToolSN = toolSerialNumber

	// 分发结果到工具进行处理
	c.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherResult), result)

	return nil
}

// seq, count
func (c *TighteningController) calBatch(workorderID int64) (int, int) {
	result, err := c.ProtocolService.storageService.FindTargetResultForJobManual(workorderID)
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
	if c.opened.Load() {
		return nil
	}
	c.opened.Store(true)

	for _, v := range c.deviceConf.Tools {
		_ = c.ProtocolService.storageService.UpdateTool(&storage.Tools{
			Serial: v.SN,
			Mode:   "pset",
		})
	}
	for _, dd := range c.dispatcherMap {
		c.dispatcherBus.LaunchDispatchersByHandlerMap(dd)
	}

	c.clearToolsResultAndCurve()

	// 启动客户端
	c.startupClients()

	return nil
}

func (c *TighteningController) startupClients() {

	for _, v := range c.sockClients {
		v.start()
	}
}

func (c *TighteningController) shutdownClients() {

	for _, v := range c.sockClients {
		v.stop()
	}
}

func (c *TighteningController) Stop() error {
	if !c.opened.Load() {
		return nil
	}
	c.opened.Store(false)

	// 停止客户端
	c.shutdownClients()

	for _, dd := range c.dispatcherMap {
		c.dispatcherBus.ReleaseDispatchersByHandlerMap(dd)
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

func (c *TighteningController) IOWrite(index uint16, status uint16) error {

	strIo := ""
	for i := 0; i < 10; i++ {
		if i == int(index) {
			switch status {
			case io.OutputStatusOff:
				strIo += "0"

			case io.OutputStatusOn:
				strIo += "1"

			case io.OutputStatusFlash:
				strIo += "2"
			}
		} else {
			strIo += "3"
		}
	}

	reply, err := c.getDefaultTransportClient().ProcessRequest(MID_0200_CONTROLLER_RELAYS, "", "", "", strIo)
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(reply.(string))
	}

	return nil
}

func (c *TighteningController) IORead() (string, string, error) {
	return "", c.inputs, nil
}

func (c *TighteningController) SetIONotify(notify io.IONotify) {
	c.ioNotify = notify
}

func (c *TighteningController) Protocol() string {
	return tightening_device.TIGHTENING_OPENPROTOCOL
}

func (c *TighteningController) clearToolsResultAndCurve() {
	for _, tool := range c.deviceConf.Tools {
		err := c.ProtocolService.storageService.ClearToolResultAndCurve(tool.SN)
		if err != nil {
			c.diag.Error(fmt.Sprintf("Clear Tool: %s Result And Curve Failed", tool.SN), err)
		}
	}
}

func (c *TighteningController) handlerOldResults() error {
	return nil
}

func (c *TighteningController) HandleStatus(sn string, status string) {
	c.mtxStatus.Lock()
	defer c.mtxStatus.Unlock()

	if status == c.Status() {
		return
	}
	c.UpdateStatus(status)

	c.NotifyIOStatus(status)
	c.diag.Info(fmt.Sprintf("OpenProtocol HandleStatus Model:%s SN:%s %s\n", c.Model(), sn, status))
	ss := []device.Status{
		{
			Type:   tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER,
			SN:     c.deviceConf.SN,
			Status: status,
			Config: c.Config(),
		},
	}

	// 分发控制器状态
	c.doDispatch(dispatcherbus.DispatcherDeviceStatus, ss)
}

func (c *TighteningController) getDefaultTransportClient() *clientContext {

	for _, sw := range c.sockClients {
		return sw
	}
	return nil
}

func (c *TighteningController) getTransportClientBySymbol(symbol string) *clientContext {

	if sw, ok := c.sockClients[symbol]; !ok {
		err := errors.Errorf("Can Not Found TransportService For %s", symbol)
		c.diag.Error("getTransportClientBySymbol", err)
		return nil
	} else {
		return sw
	}
}

func (c *TighteningController) PSetSubscribe(sn string) error {

	reply, err := c.getClient(sn).ProcessRequest(MID_0014_PSET_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0060_LAST_RESULT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) SelectorSubscribe(sn string) error {
	if c.hasSubscribedSelector.Load() {
		return nil
	}
	c.hasSubscribedSelector.Store(true)

	reply, err := c.getDefaultTransportClient().ProcessRequest(MID_0250_SELECTOR_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0250_SELECTOR_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) JobInfoSubscribe(sn string) error {

	reply, err := c.getClient(sn).ProcessRequest(MID_0034_JOB_INFO_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0034_JOB_INFO_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) IOInputSubscribe(sn string) error {
	if c.hasSubscribedIO.Load() {
		return nil
	}
	c.hasSubscribedIO.Store(true)

	reply, err := c.getDefaultTransportClient().ProcessRequest(MID_0210_INPUT_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0210_INPUT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) MultiSpindleResultSubscribe(sn string) error {

	reply, err := c.getClient(sn).ProcessRequest(MID_0100_MULTI_SPINDLE_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0100_MULTI_SPINDLE_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) VinSubscribe(sn string) error {
	if c.hasSubscribedVin.Load() {
		return nil
	}
	c.hasSubscribedVin.Store(true)

	reply, err := c.getDefaultTransportClient().ProcessRequest(MID_0051_VIN_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0051_VIN_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) ResultSubcribe(sn string) error {

	reply, err := c.getClient(sn).ProcessRequest(MID_0060_LAST_RESULT_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0060_LAST_RESULT_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) AlarmSubcribe(sn string) error {

	reply, err := c.getClient(sn).ProcessRequest(MID_0070_ALARM_SUBSCRIBE, "1", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] {
		return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0070_ALARM_SUBSCRIBE, reply.(string)))
	}

	return nil
}

func (c *TighteningController) CurveSubscribe(sn string) error {
	reply, err := c.getClient(sn).ProcessRequest(MID_7408_LAST_CURVE_SUBSCRIBE, "1", "", "", "")
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

func (c *TighteningController) GetVendorModel() map[string]interface{} {
	return nil
}

func (c *TighteningController) GetVendorMid(mid string) (string, error) {
	rev, exist := c.getInstance().GetVendorModel()[mid]
	if !exist {
		return "", errors.New(fmt.Sprintf("MID %s Not Supported", mid))
	}

	return rev.(string), nil
}

func (c *TighteningController) New() IOpenProtocolController {
	return &TighteningController{}
}

func (c *TighteningController) getInstance() IOpenProtocolController {
	if c.instance == nil {
		panic("Controller Instance Is Nil")
	}

	return c.instance
}

func (c *TighteningController) SetInstance(instance IOpenProtocolController) {
	c.instance = instance
}

func (c *TighteningController) doDispatch(name string, data interface{}) {
	if err := c.dispatcherBus.Dispatch(name, data); err != nil {
		c.diag.Error("Dispatch Failed", err)
	}
}

func (c *TighteningController) OpenProtocolParams() *OpenProtocolParams {
	return &OpenProtocolParams{
		MaxKeepAliveCheck: 3,
		MaxReplyTime:      3 * time.Second,
		KeepAlivePeriod:   8 * time.Second,
	}
}

func (c *TighteningController) ioSerialNumber() string {
	return fmt.Sprintf(tightening_device.TIGHTENING_CONTROLLER_IO_SN_FORMAT, c.SerialNumber())
}

func (c *TighteningController) NotifyIOStatus(status string) {
	if c.ioNotify == nil {
		return
	}

	c.internalIO.(*TighteningIO).UpdateStatus(status)
	c.ioNotify.OnStatus(c.ioSerialNumber(), status)
}

func (c *TighteningController) NotifyIOContact(t string, status string) {
	if c.ioNotify == nil {
		return
	}

	c.ioNotify.OnChangeIOStatus(c.ioSerialNumber(), t, status)
}
