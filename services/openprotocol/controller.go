package openprotocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
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

	REPLY_TIMEOUT   = time.Duration(100 * time.Millisecond)
	MAX_REPLY_COUNT = 20
)

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}
type handlerPkg_curve struct {
	Header *OpenProtocolHeader
	Body   []byte
}
type TighteningController struct {
	w                 *socket_writer.SocketWriter
	cfg               *tightening_device.TighteningDeviceConfig
	StatusValue       atomic.Value
	keepAliveCount    int32
	keep_period       time.Duration
	req_timeout       time.Duration
	getToolInfoPeriod time.Duration
	Response          ResponseQueue
	Srv               *Service
	dbController      *storage.Controllers
	buffer            chan []byte
	closing           chan chan struct{}
	handlerBuf        chan handlerPkg
	keepaliveDeadLine atomic.Value
	protocol          string
	Mode              atomic.Value
	inputs            string
	diag              Diagnostic
	MID_7410_CURVE    handlerPkg_curve
	result_CURVE      *minio.ControllerCurve
	result            *controller.ControllerResult
	mtxResult         sync.Mutex
	toolStatus        atomic.Value
	model             string
	receiveBuf        chan []byte

	tightening_device.ITighteningController
	device.BaseDevice
}

// TODO: 如果工具序列号没有配置，则通过探测加入设备列表。
func NewController(protocolConfig *Config, deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic) TighteningController {

	c := TighteningController{
		diag:              d,
		buffer:            make(chan []byte, 1024),
		closing:           make(chan chan struct{}),
		keep_period:       time.Duration(protocolConfig.KeepAlivePeriod),
		req_timeout:       time.Duration(protocolConfig.ReqTimeout),
		getToolInfoPeriod: time.Duration(protocolConfig.GetToolInfoPeriod),
		Response:          ResponseQueue{},
		handlerBuf:        make(chan handlerPkg, 1024),
		protocol:          controller.OPENPROTOCOL,
		result_CURVE:      nil,
		mtxResult:         sync.Mutex{},
		receiveBuf:        make(chan []byte, 65535),
		cfg:               deviceConfig,
		BaseDevice:        device.CreateBaseDevice(),
	}

	c.StatusValue.Store(controller.STATUS_OFFLINE)
	c.toolStatus.Store(controller.EVT_TOOL_DISCONNECTED)

	for _, v := range deviceConfig.Tools {
		tool := NewTool(&c, v, d)
		c.AddChildren(v.SN, &tool)
	}

	return c
}

func (c *TighteningController) Tools() map[string]tightening_device.ITighteningTool {
	//rt := map[string]string{}

	//toolStatus := c.toolStatus.Load().(string)
	//for _, v := range c.cfg.Tools {
	//	rt[v.SerialNO] = toolStatus
	//}

	return nil
}

func (c *TighteningController) UpdateToolStatus(status string) {
	s := c.toolStatus.Load().(string)
	if s != status {
		c.toolStatus.Store(status)

		// 推送工具状态
		//ts := wsnotify.WSToolStatus{
		//	ToolSN: c.cfg.Tools[0].SerialNO,
		//	Status: status,
		//}
		//
		//str, _ := json.Marshal(ts)
		//c.Srv.WS.WSSend(wsnotify.WS_EVETN_TOOL, string(str))
	}
}

func (c *TighteningController) LoadController(controller *storage.Controllers) {
	c.dbController = controller
}

func (c *TighteningController) Inputs() string {
	return c.inputs
}

func (c *TighteningController) handlerProcess() {
	for {
		select {
		case pkg := <-c.handlerBuf:
			err := c.HandleMsg(&pkg)
			if err != nil {
				c.diag.Error("OP protocol HandleMsg fail", err)
			}
		}
	}
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

// TODO: case封装成函数
func (c *TighteningController) HandleMsg(pkg *handlerPkg) error {
	c.Srv.diag.Debug(fmt.Sprintf("%s: %s%s\n", c.cfg.SN, pkg.Header.Serialize(), pkg.Body))

	switch pkg.Header.MID {
	case MID_7410_LAST_CURVE:
		//结果曲线
		if c.result_CURVE == nil {
			c.result_CURVE = &minio.ControllerCurve{}
		}

		torqueCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(pkg.Body[27:41]), 64)
		angleCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(pkg.Body[43:57]), 64)
		if pkg.Body[69:71] == "01" {
			c.MID_7410_CURVE.Header = &pkg.Header
			c.MID_7410_CURVE.Body = []byte{}
		}
		c.MID_7410_CURVE.Body = append(c.MID_7410_CURVE.Body, []byte(pkg.Body[71:len(pkg.Body)])...)
		if pkg.Body[69:71] == pkg.Body[65:67] {
			Torque, Angle := DataDecoding(c.MID_7410_CURVE.Body, torqueCoefficient, angleCoefficient, c.diag)
			c.diag.Debug(fmt.Sprintf("Torque: %+v", Torque))
			c.diag.Debug(fmt.Sprintf("Angle: %+v", Angle))

			var curve minio.ControllerCurve
			curve.CurveContent.CUR_M = Torque
			curve.CurveContent.CUR_W = Angle

			//c.result_CURVE.CurveContent.Result=c.TightingID
			//c.Srv.Parent.Handlers.HandleCurve(c.result_CURVE)

			c.MID_7410_CURVE.Header = nil
			c.MID_7410_CURVE.Body = nil

			c.updateResult(nil, &curve)
			c.handleResultandClear()
		}
		/*
			fmt.Println([]byte(pkg.Body[71:]))
			fmt.Println(pkg.Header)
			fmt.Println(pkg.Body[65:67], "***********", pkg.Body[69:71])
		*/
	case MID_0061_LAST_RESULT:
		// 结果数据

		result_data := ResultData{}
		result_data.Deserialize(pkg.Body)

		return c.handleResult(&result_data, nil)

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

			return c.handleResult(&result_data, nil)
		}

	case MID_0013_PSET_DETAIL_REPLY:
		// pset详细数据
		pset_detail := PSetDetail{}
		pset_detail.Deserialize(pkg.Body)

		c.Response.update(MID_0012_PSET_DETAIL_REQUEST, pset_detail)

	case MID_0011_PSET_LIST_REPLY:
		// pset列表
		pset_list := PSetList{}
		pset_list.Deserialize(pkg.Body)

		c.Response.update(MID_0010_PSET_LIST_REQUEST, pset_list)

		//case MID_7410_LAST_CURVE:
		// 处理波形

	case MID_0031_JOB_LIST_REPLY:
		job_list := JobList{}
		job_list.Deserialize(pkg.Body)

		c.Response.update(MID_0030_JOB_LIST_REQUEST, job_list)

	case MID_0033_JOB_DETAIL_REPLY:
		// job详细数据
		job_detail := JobDetail{}
		job_detail.Deserialize(pkg.Body)

		c.Response.update(MID_0032_JOB_DETAIL_REQUEST, job_detail)

	case MID_0004_CMD_ERR:
		// 请求错误

		err_code := pkg.Body[4:6]
		c.Response.update(pkg.Body[0:4], request_errors[err_code])

	case MID_0005_CMD_OK:
		// 请求正确

		c.Response.update(pkg.Body, request_errors["00"])

	case MID_0035_JOB_INFO:
		job_info := JobInfo{}
		job_info.Deserialize(pkg.Body)

		if job_info.JobCurrentStep == 0 {
			// 推送job选择信息

			job_select := wsnotify.WSJobSelect{
				JobID: job_info.JobID,
			}

			ws_str, _ := json.Marshal(job_select)
			c.Srv.WS.WSSendJob(string(ws_str))
		}

	case MID_0211_INPUT_MONITOR:
		inputs := IOMonitor{}
		inputs.Deserialize(pkg.Body)

		inputs.ControllerSN = c.cfg.SN

		c.inputs = inputs.Inputs

		str, _ := json.Marshal(inputs)
		go c.Srv.WS.WSSendIOInput(string(str))

	case MID_0101_MULTI_SPINDLE_RESULT:
		ms := MultiSpindleResult{}
		ms.Deserialize(pkg.Body)

		wsResults := make([]wsnotify.WSResult, len(ms.Spindles), len(ms.Spindles))
		//wsResult := wsnotify.WSResult{}
		for idx, v := range ms.Spindles {
			wsResults[idx].Result = v.Result
			wsResults[idx].MI = v.Torque
			wsResults[idx].WI = v.Angle
			//wsResults = append(wsResults, wsResult)
		}

		wsStrs, err := json.Marshal(wsResults)
		if err == nil {
			c.Srv.WS.WSSend(wsnotify.WS_EVENT_RESULT, string(wsStrs))
		}

		return err

	case MID_0052_VIN:
		// 收到条码
		ids := DeserializeIDS(pkg.Body)

		bc := ""
		for _, v := range c.Srv.config().VinIndex {
			if v < 0 || v > (MAX_IDS_NUM-1) {
				continue
			}

			bc += ids[v]
		}

		barcode := wsnotify.WSScanner{
			Barcode: bc,
		}

		str, _ := json.Marshal(barcode)

		c.Srv.WS.WSSend(wsnotify.WS_EVENT_SCANNER, string(str))
	case MID_0071_ALARM:
		// 收到报警信息
		var ai AlarmInfo
		err := ai.Deserialize(pkg.Body)
		if err != nil {
			c.diag.Error("alarm info deserialize fail", err)
			return err
		} else {
			// 参见 项目管理,长安项目中文件:http://116.62.21.97/web#id=325&view_type=form&model=ir.attachment&active_id=3&menu_id=90
			// 第11页,错误代码:Tool calibration required:E305
			//if ai.ErrorCode == "E305" {
			//	// do nothing,当前未确认是否为这个错误代码
			//}

			switch ai.ErrorCode {
			case EVT_CONTROLLER_TOOL_CONNECT:
				c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

			case EVT_CONTROLLER_TOOL_DISCONNECT:
				c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
			}
		}
		return nil

	case MID_0076_ALARM_STATUS:
		var as AlarmStatus
		err := as.Deserialize(pkg.Body)
		if err != nil {
			c.diag.Error("alarm status deserialize fail", err)
		} else {
			switch as.ErrorCode {
			case EVT_CONTROLLER_NO_ERR:
				c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

			case EVT_CONTROLLER_TOOL_DISCONNECT:
				c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
			}
		}

	case MID_0041_TOOL_INFO_REPLY:
		// 收到工具信息
		var ti ToolInfo
		err := ti.Deserialize(pkg.Body)

		if c.Status() == controller.STATUS_OFFLINE {
			c.Response.update(MID_0040_TOOL_INFO_REQUEST, ti)
		} else {
			if err != nil {
				c.diag.Error("tool info deserialize fail", err)
			} else {

				// 将数据通过api传给odoo
				if ti.ToolSN == "" {
					return errors.New("Tool Serial Number is empty string")
				}

				if ti.TotalTighteningCount == 0 || ti.CountSinLastService == 0 {
					//不需要尝试创建维修/标定单据
					return nil
				}

				go c.Srv.TryCreateMaintenance(ti) // 协程处理
			}
		}

	}
	return nil
}

func ArrayContains(s []int, e int) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}

func (c *TighteningController) handleResult(result_data *ResultData, carve *minio.ControllerCurve) error {

	if ArrayContains(c.Srv.config().SkipJobs, result_data.JobID) {
		return nil
	}

	if c.dbController != nil {
		c.Srv.DB.UpdateTightning(c.dbController.Id, result_data.TightingID)
	}

	controllerResult := controller.ControllerResult{}
	controllerResult.NeedPushHmi = true

	if c.Model() == tightening_device.ModelDesoutterDeltaWrench {
		controllerResult.GunSN = c.cfg.SN
	} else {
		controllerResult.GunSN = result_data.ToolSerialNumber
	}

	gun, err := c.Srv.DB.GetGun(controllerResult.GunSN)
	if err != nil {
		c.Srv.diag.Error("get gun failed", err)
		return err
	}

	psetTrace := tightening_device.PSetSet{}
	_ = json.Unmarshal([]byte(gun.Trace), &psetTrace)

	controllerResult.TighteningID = result_data.TightingID
	controllerResult.Count = psetTrace.Count
	controllerResult.Batch = fmt.Sprintf("%d/%d", psetTrace.Sequence, psetTrace.Total)

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

	if result_data.ResultType == "02" {
		controllerResult.Result = storage.RESULT_LSN
		controllerResult.NeedPushHmi = true
		controllerResult.PSetDefine.Strategy = controller.STRATEGY_LN
	}

	controllerResult.PSetDefine.Mp = result_data.TorqueMax / 100
	controllerResult.PSetDefine.Mm = result_data.TorqueMin / 100
	controllerResult.PSetDefine.Ma = result_data.TorqueFinalTarget / 100

	controllerResult.PSetDefine.Wp = result_data.AngleMax
	controllerResult.PSetDefine.Wm = result_data.AngleMin
	controllerResult.PSetDefine.Wa = result_data.FinalAngleTarget

	controllerResult.ExceptionReason = result_data.TighteningErrorStatus

	//targetID := result_data.VIN
	//switch c.Srv.config().DataIndex {
	//case 1:
	//	targetID = result_data.ID2
	//case 2:
	//	targetID = result_data.ID3
	//case 3:
	//	targetID = result_data.ID4
	//}

	controllerResult.Workorder_ID = psetTrace.WorkorderID
	controllerResult.NeedPushAiis = true

	//controllerResult.Seq, controllerResult.Count = c.calBatch(controllerResult.Workorder_ID)

	//c.result = &controllerResult

	//c.Srv.Parent.Handlers.Handle(&controllerResult, c.result_CURVE)
	c.updateResult(&controllerResult, nil)
	c.handleResultandClear()

	return nil
}

func (c *TighteningController) updateResult(result *controller.ControllerResult, curve *minio.ControllerCurve) {
	defer c.mtxResult.Unlock()
	c.mtxResult.Lock()

	if result != nil {
		c.result = result
	}

	if curve != nil {
		c.result_CURVE = curve
	}
}

func (c *TighteningController) handleResultandClear() {
	defer c.mtxResult.Unlock()
	c.mtxResult.Lock()

	if c.result != nil && c.result_CURVE != nil {

		if c.result_CURVE != nil {
			c.result_CURVE.CurveContent.Result = c.result.Result
			c.result_CURVE.CurveFile = fmt.Sprintf("%s-%s.json", c.cfg.SN, c.result.TighteningID)
		}

		c.Srv.Parent.Handlers.Handle(*c.result, *c.result_CURVE)
		c.result = nil
		c.result_CURVE = nil
	}
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

	_ = c.Srv.DB.UpdateGun(&storage.Guns{
		Serial: c.cfg.SN,
		Mode:   "pset",
	})

	for _, v := range c.cfg.Tools {
		_ = c.Srv.DB.UpdateGun(&storage.Guns{
			Serial: v.SN,
			Mode:   "pset",
		})
	}

	c.w = socket_writer.NewSocketWriter(c.cfg.Endpoint, c)

	go c.handlerProcess()

	c.Connect()

	return nil
}

func (c *TighteningController) Stop() error {
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
	return nil
}

func (c *TighteningController) Protocol() string {
	return c.protocol
}

func (c *TighteningController) Connect() error {
	c.StatusValue.Store(controller.STATUS_OFFLINE)
	c.Response = ResponseQueue{
		Results: map[string]interface{}{},
		mtx:     sync.Mutex{},
	}

	c.Mode.Store(MODE_JOB)

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
	c.ToolInfoReq()

	//c.JobOff("1")
	c.PSetSubscribe()

	c.ResultSubcribe()
	c.CurveSubscribe()
	c.SelectorSubscribe()

	c.JobInfoSubscribe()
	c.IOInputSubscribe()
	//c.MultiSpindleResultSubscribe()
	c.VinSubscribe()

	c.AlarmSubcribe()

	// 启动发送
	go c.manage()

	go c.SolveOldResults()

	go c.getTighteningCount()

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

			if c.Status() == controller.STATUS_OFFLINE {
				continue
			}
			req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", "", "", "")
			c.Write([]byte(req))
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *TighteningController) ToolInfoReq() error {
	rev, err := GetVendorMid(c.Model(), MID_0040_TOOL_INFO_REQUEST)
	if err != nil {
		return err
	}

	//defer c.Response.remove(MID_0040_TOOL_INFO_REQUEST)
	//c.Response.Add(MID_0040_TOOL_INFO_REQUEST, nil)

	req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", "", "", "")
	c.Write([]byte(req))

	//var reply interface{} = nil
	//
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
	//	c.updateStatus(controller.STATUS_ONLINE)
	//} else {
	//	c.updateStatus(controller.STATUS_OFFLINE)
	//}

	return nil
}

func (c *TighteningController) GetPSetList() ([]int, error) {
	var psets []int

	rev, err := GetVendorMid(c.Model(), MID_0010_PSET_LIST_REQUEST)
	if err != nil {
		return psets, err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return psets, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0010_PSET_LIST_REQUEST)
	c.Response.Add(MID_0010_PSET_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0010_PSET_LIST_REQUEST, rev, "", "", "", "")
	c.Write([]byte(psets_request))

	var reply interface{} = nil

	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0010_PSET_LIST_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return psets, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	pset_list := reply.(PSetList)

	return pset_list.psets, nil
}

func (c *TighteningController) GetPSetDetail(pset int) (PSetDetail, error) {
	var obj_pset_detail PSetDetail

	rev, err := GetVendorMid(c.Model(), MID_0012_PSET_DETAIL_REQUEST)
	if err != nil {
		return obj_pset_detail, err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_pset_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0012_PSET_DETAIL_REQUEST)
	c.Response.Add(MID_0012_PSET_DETAIL_REQUEST, nil)

	pset_detail := GeneratePackage(MID_0012_PSET_DETAIL_REQUEST, rev, "", "", "", fmt.Sprintf("%03d", pset))
	c.Write([]byte(pset_detail))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0012_PSET_DETAIL_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return obj_pset_detail, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	switch v := reply.(type) {
	case string:
		return obj_pset_detail, errors.New(v)

	case PSetDetail:
		return reply.(PSetDetail), nil

	default:
		return obj_pset_detail, errors.New(controller.ERR_KNOWN)
	}

}

func (c *TighteningController) GetJobList() ([]int, error) {
	var jobs []int
	rev, err := GetVendorMid(c.Model(), MID_0030_JOB_LIST_REQUEST)
	if err != nil {
		return jobs, err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return jobs, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0030_JOB_LIST_REQUEST)
	c.Response.Add(MID_0030_JOB_LIST_REQUEST, nil)

	psetsRequest := GeneratePackage(MID_0030_JOB_LIST_REQUEST, rev, "", "", "", "")
	c.Write([]byte(psetsRequest))

	var reply interface{} = nil

	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0030_JOB_LIST_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return jobs, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	jobList := reply.(JobList)

	return jobList.jobs, nil
}

func (c *TighteningController) GetJobDetail(job int) (JobDetail, error) {
	var objJobDetail JobDetail
	rev, err := GetVendorMid(c.Model(), MID_0032_JOB_DETAIL_REQUEST)
	if err != nil {
		return objJobDetail, err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return objJobDetail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0032_JOB_DETAIL_REQUEST)
	c.Response.Add(MID_0032_JOB_DETAIL_REQUEST, nil)

	job_detail := GeneratePackage(MID_0032_JOB_DETAIL_REQUEST, rev, "", "", "", fmt.Sprintf("%04d", job))
	c.Write([]byte(job_detail))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0032_JOB_DETAIL_REQUEST)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return objJobDetail, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	switch v := reply.(type) {
	case string:
		return objJobDetail, errors.New(v)

	case JobDetail:
		return reply.(JobDetail), nil

	default:
		return objJobDetail, errors.New(controller.ERR_KNOWN)
	}

}

func (c *TighteningController) SolveOldResults() {

	if c.dbController == nil || c.dbController.LastID == "0" {
		return
	}

	c.Response.Add(MID_0064_OLD_SUBSCRIBE, MID_0064_OLD_SUBSCRIBE)
	if c.getOldResult(0) != nil {
		return
	}

	var lastResult interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		lastResult = c.Response.get(MID_0065_OLD_DATA)
		if lastResult != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if lastResult == nil {
		return
	}

	objLastResult := lastResult.(ResultData)

	if objLastResult.TightingID != c.dbController.LastID {
		startId, _ := strconv.ParseInt(c.dbController.LastID, 10, 64)
		endId, _ := strconv.ParseInt(objLastResult.TightingID, 10, 64)

		for i := startId + 1; i <= endId; i++ {
			c.getOldResult(i)
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
	c.keepaliveDeadLine.Store(time.Now().Add(c.keep_period))
}

func (c *TighteningController) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *TighteningController) Status() string {

	return c.StatusValue.Load().(string)
}

func (c *TighteningController) sendKeepalive() {
	if c.Status() == controller.STATUS_OFFLINE {
		return
	}

	keepAlive := GeneratePackage(MID_9999_ALIVE, DEFAULT_REV, "1", "", "", "")
	c.Write([]byte(keepAlive))
}

func (c *TighteningController) startComm() error {
	rev, err := GetVendorMid(c.Model(), MID_0001_START)
	if err != nil {
		return err
	}

	//if c.Status() == controller.STATUS_OFFLINE {
	//	return errors.New("status offline")
	//}

	start := GeneratePackage(MID_0001_START, rev, "", "", "", "")

	//c.Response.Add(MID_0001_START, "")

	c.Write([]byte(start))

	return nil
}

func (c *TighteningController) Write(buf []byte) {
	c.buffer <- buf
}

func (c *TighteningController) updateStatus(status string) {

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

		c.Srv.diag.Debug(fmt.Sprintf("%s:%s %s\n", c.Model(), c.cfg.SN, status))

	}
}

func (c *TighteningController) Read(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			c.Srv.diag.Error("read failed", err)
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

	c.diag.Debug(fmt.Sprintf("%s op target buf: %s", c.cfg.SN, string(msg)))

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
					c.diag.Debug("Index Is Empty")
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
		}
	}
}

func (c *TighteningController) getOldResult(last_id int64) error {
	rev, err := GetVendorMid(c.Model(), MID_0064_OLD_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_last_result := GeneratePackage(MID_0064_OLD_SUBSCRIBE, rev, "", "", "", fmt.Sprintf("%010d", last_id))

	c.Write([]byte(s_last_result))

	return nil
}

func (c *TighteningController) pset(pset int) error {
	rev, err := GetVendorMid(c.Model(), MID_0018_PSET)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0018_PSET, nil)
	defer c.Response.remove(MID_0018_PSET)

	s_pset := GeneratePackage(MID_0018_PSET, rev, "", "", "", fmt.Sprintf("%03d", pset))

	c.Write([]byte(s_pset))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0018_PSET)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

func (c *TighteningController) ToolControl(enable bool) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_cmd := MID_0042_TOOL_DISABLE
	if enable {
		s_cmd = MID_0043_TOOL_ENABLE
	}

	rev, err := GetVendorMid(c.Model(), s_cmd)
	if err != nil {
		return err
	}

	c.Response.Add(s_cmd, nil)
	defer c.Response.remove(s_cmd)

	sSend := GeneratePackage(s_cmd, rev, "", "", "", "")

	c.Write([]byte(sSend))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(s_cmd)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	sReply := reply.(string)
	if sReply != request_errors["00"] {
		return errors.New(sReply)
	}

	return nil
}

// 0: set 1: reset
func (c *TighteningController) JobOff(off string) error {
	rev, err := GetVendorMid(c.Model(), MID_0130_JOB_OFF)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0130_JOB_OFF, nil)
	defer c.Response.remove(MID_0130_JOB_OFF)

	s_off := GeneratePackage(MID_0130_JOB_OFF, rev, "", "", "", off)

	c.Write([]byte(s_off))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0130_JOB_OFF)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	if off == "0" {
		c.Mode.Store(MODE_PSET)
	} else {
		c.Mode.Store(MODE_JOB)
	}

	return nil
}

func (c *TighteningController) jobSelect(job int) error {
	rev, err := GetVendorMid(c.Model(), MID_0038_JOB_SELECT)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0038_JOB_SELECT, nil)
	defer c.Response.remove(MID_0038_JOB_SELECT)

	s_job := GeneratePackage(MID_0038_JOB_SELECT, rev, "", "", "", fmt.Sprintf("%04d", job))

	c.Write([]byte(s_job))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0038_JOB_SELECT)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

func (c *TighteningController) IdentifierSet(str string) error {
	rev, err := GetVendorMid(c.Model(), MID_0150_IDENTIFIER_SET)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, rev, "", "", "", str)

	c.Write([]byte(ide))

	return nil
}

func (c *TighteningController) PSetBatchSet(pset int, batch int) error {
	rev, err := GetVendorMid(c.Model(), MID_0019_PSET_BATCH_SET)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0019_PSET_BATCH_SET, nil)
	defer c.Response.remove(MID_0019_PSET_BATCH_SET)

	s := fmt.Sprintf("%03d%02d", pset, batch)
	ide := GeneratePackage(MID_0019_PSET_BATCH_SET, rev, "", "", "", s)

	c.Write([]byte(ide))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0019_PSET_BATCH_SET)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

func (c *TighteningController) PSetBatchReset(pset int) error {
	rev, err := GetVendorMid(c.Model(), MID_0020_PSET_BATCH_RESET)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s := fmt.Sprintf("%03d", pset)
	ide := GeneratePackage(MID_0020_PSET_BATCH_RESET, rev, "", "", "", s)

	c.Write([]byte(ide))

	return nil
}

func (c *TighteningController) PSetSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0014_PSET_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0014_PSET_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(pset))

	return nil
}

func (c *TighteningController) SelectorSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0250_SELECTOR_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0250_SELECTOR_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(pset))

	return nil
}

func (c *TighteningController) JobInfoSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0034_JOB_INFO_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0034_JOB_INFO_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(pset))

	return nil
}

func (c *TighteningController) IOInputSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0210_INPUT_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0210_INPUT_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(input))

	return nil
}

func (c *TighteningController) MultiSpindleResultSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0100_MULTI_SPINDLE_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0100_MULTI_SPINDLE_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(input))

	return nil
}

func (c *TighteningController) VinSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0051_VIN_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0051_VIN_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(input))

	return nil
}

func (c *TighteningController) ResultSubcribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0060_LAST_RESULT_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0060_LAST_RESULT_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(pset))

	return nil
}

func (c *TighteningController) AlarmSubcribe() error {
	rev, err := GetVendorMid(c.Model(), MID_0070_ALARM_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	payload := GeneratePackage(MID_0070_ALARM_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(payload))

	return nil
}

func (c *TighteningController) CurveSubscribe() error {
	rev, err := GetVendorMid(c.Model(), MID_7408_LAST_CURVE_SUBSCRIBE)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_7408_LAST_CURVE_SUBSCRIBE, rev, "1", "", "", "")

	c.Write([]byte(pset))

	return nil
}

func (c *TighteningController) PSet(pset int, channel int, ex_info string, count int) (uint32, error) {

	if c.Mode.Load().(string) != MODE_PSET {
		return 0, errors.New("current mode is not pset")
	}

	// 次数控制
	//c.PSetBatchSet(pset, count)

	// 结果id-拧接次数-用户id
	c.IdentifierSet(ex_info)

	// 设定pset
	err := c.pset(pset)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (c *TighteningController) findIOByNo(no int, ios *[]IOStatus) (IOStatus, error) {
	for _, v := range *ios {
		if no == v.No {
			return v, nil
		}
	}

	return IOStatus{}, errors.New("not found")
}

func (c *TighteningController) IOSet(ios *[]IOStatus) error {
	rev, err := GetVendorMid(c.Model(), MID_0200_CONTROLLER_RELAYS)
	if err != nil {
		return err
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	if len(*ios) == 0 {
		return errors.New("io status list is required")
	}

	str_io := ""
	for i := 0; i < 10; i++ {
		io, err := c.findIOByNo(i, ios)
		if err != nil {
			str_io += "3"
		} else {
			switch io.Status {
			case IO_STATUS_OFF:
				str_io += "0"

			case IO_STATUS_ON:
				str_io += "1"

			case IO_STATUS_FLASHING:
				str_io += "2"
			}
		}
	}

	c.Response.Add(MID_0200_CONTROLLER_RELAYS, nil)
	defer c.Response.remove(MID_0200_CONTROLLER_RELAYS)

	s_io := GeneratePackage(MID_0200_CONTROLLER_RELAYS, rev, "", "", "", str_io)

	c.Write([]byte(s_io))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0200_CONTROLLER_RELAYS)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

func (c *TighteningController) JobSet(id_info string, job int) error {

	if c.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	_ = c.IdentifierSet(id_info)

	err := c.jobSelect(job)
	if err != nil {
		return err
	}

	return nil
}

func (c *TighteningController) JobAbort() error {
	rev, err := GetVendorMid(c.Model(), MID_0127_JOB_ABORT)
	if err != nil {
		return err
	}

	if c.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	c.Response.Add(MID_0127_JOB_ABORT, nil)
	defer c.Response.remove(MID_0127_JOB_ABORT)

	s_job := GeneratePackage(MID_0127_JOB_ABORT, rev, "", "", "", "")

	c.Write([]byte(s_job))

	var reply interface{} = nil
	for i := 0; i < MAX_REPLY_COUNT; i++ {
		reply = c.Response.get(MID_0127_JOB_ABORT)
		if reply != nil {
			break
		}

		time.Sleep(REPLY_TIMEOUT)
	}

	if reply == nil {
		return errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	sReply := reply.(string)
	if sReply != request_errors["00"] {
		return errors.New(sReply)
	}

	return nil
}

func (c *TighteningController) Model() string {
	return c.model
}

func (c *TighteningController) SetModel(model string) {
	c.model = model
}

func (c *TighteningController) SetJob(r *tightening_device.JobSet) tightening_device.Reply {
	return tightening_device.Reply{}
}

func (c *TighteningController) SetPSet(r *tightening_device.PSetSet) tightening_device.Reply {
	rt := tightening_device.Reply{
		Result: 0,
		Msg:    "",
	}

	_ = c.PSetBatchSet(r.PSet, 1)

	err := c.pset(r.PSet)
	if err != nil {
		rt.Result = -1
		rt.Msg = err.Error()
	}

	return rt
}

func (c *TighteningController) Enable(r *tightening_device.ToolEnable) tightening_device.Reply {
	rt := tightening_device.Reply{
		Result: 0,
		Msg:    "",
	}

	err := c.ToolControl(r.Enable)
	if err != nil {
		rt.Result = -1
		rt.Msg = err.Error()
	}

	return rt
}

func (c *TighteningController) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER
}

func (c *TighteningController) Children() map[string]device.IDevice {

	return map[string]device.IDevice{}
}

func (s *TighteningController) Data() interface{} {
	return nil
}

func (s *TighteningController) Config() interface{} {
	return nil
}
