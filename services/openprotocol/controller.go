package openprotocol

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
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
type Controller struct {
	w                 *socket_writer.SocketWriter
	cfg               controller.ControllerConfig
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
	TriggerStart      time.Time
	TriggerStop       time.Time
	inputs            string
	diag              Diagnostic
	MID_7410_CURVE    handlerPkg_curve
	result_CURVE      *minio.ControllerCurve
	toolStatus        atomic.Value
	model             string
	tighteningDevice  *tightening_device.Service

	tightening_device.TighteningDevice
}

func NewController(c Config, d Diagnostic) Controller {

	cont := Controller{
		diag:              d,
		buffer:            make(chan []byte, 1024),
		closing:           make(chan chan struct{}),
		keep_period:       time.Duration(c.KeepAlivePeriod),
		req_timeout:       time.Duration(c.ReqTimeout),
		getToolInfoPeriod: time.Duration(c.GetToolInfoPeriod),
		Response:          ResponseQueue{},
		handlerBuf:        make(chan handlerPkg, 1024),
		protocol:          controller.OPENPROTOCOL,
		result_CURVE:      nil,
	}

	cont.StatusValue.Store(controller.STATUS_OFFLINE)
	cont.toolStatus.Store(controller.EVT_TOOL_DISCONNECTED)

	return cont
}

func (c *Controller) Tools() map[string]string {
	rt := map[string]string{}

	toolStatus := c.toolStatus.Load().(string)
	for _, v := range c.cfg.Tools {
		rt[v.SerialNO] = toolStatus
	}

	return rt
}

func (c *Controller) UpdateToolStatus(status string) {
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

func (c *Controller) LoadController(controller *storage.Controllers) {
	c.dbController = controller
}

func (c *Controller) Inputs() string {
	return c.inputs
}

func (c *Controller) handlerProcess() {
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

func (c *Controller) Data_decoding(original []byte, Torque_Coefficient float64, Angle_Coefficient float64) (Torque []float64, Angle []float64) {
	var byte_num bool
	var data []byte
	for i, _ := range original {
		if original[i] == 0xff && !byte_num {
			byte_num = true
			continue
		}
		if byte_num {
			if original[i] == 0xff {
				data = append(data, 0xfe)
			} else if original[i] == 0xfe {
				data = append(data, 0xff)
			}
			byte_num = false
		} else {
			data = append(data, original[i]-1)
		}
	}
	for i := 0; i < len(data)/6; i++ {
		_ = data[i*6+1]
		_ = data[i*6+5]
		a := uint16(data[i*6]) | uint16(data[i*6+1])<<8
		b := uint32(data[i*6+2]) | uint32(data[i*6+3])<<8 | uint32(data[i*6+4])<<16 | uint32(data[i*6+5])<<24
		Torque = append(Torque, float64(a)*Torque_Coefficient)
		Angle = append(Angle, float64(b)*Angle_Coefficient)
	}
	return
}

func (c *Controller) HandleMsg(pkg *handlerPkg) error {
	c.Srv.diag.Debug(fmt.Sprintf("%s%s\n", pkg.Header.Serialize(), pkg.Body))

	switch pkg.Header.MID {
	case MID_7410_LAST_CURVE:
		//结果曲线
		if c.result_CURVE == nil {
			c.result_CURVE = &minio.ControllerCurve{}
		}

		Torque_Coefficient, _ := strconv.ParseFloat(strings.TrimSpace(pkg.Body[27:41]), 64)
		Angle_Coefficient, _ := strconv.ParseFloat(strings.TrimSpace(pkg.Body[43:57]), 64)
		if pkg.Body[69:71] == "01" {
			c.MID_7410_CURVE.Header = &pkg.Header
			c.MID_7410_CURVE.Body = []byte(pkg.Body)
		}
		if pkg.Body[69:71] == pkg.Body[65:67] {
			Torque, Angle := c.Data_decoding(c.MID_7410_CURVE.Body[71:], Torque_Coefficient, Angle_Coefficient)
			//fmt.Println(Torque, Angle)
			//fmt.Println(c.MID_7410_CURVE.Body[71:])

			c.result_CURVE.CurveContent = minio.ControllerCurveFile{}
			c.result_CURVE.CurveContent.CUR_M = Torque
			c.result_CURVE.CurveContent.CUR_W = Angle

			//c.result_CURVE.CurveContent.Result=c.TightingID
			//c.Srv.Parent.Handlers.HandleCurve(c.result_CURVE)

			c.MID_7410_CURVE.Header = nil
			c.MID_7410_CURVE.Body = nil

		} else {
			c.MID_7410_CURVE.Body = append(c.MID_7410_CURVE.Body, []byte(pkg.Body[71:len(pkg.Body)])...)
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

		if inputs.Inputs[c.Srv.config().IOTrigger-1] == '1' {
			// 开始trigger计时
			c.TriggerStart = time.Now()
			go c.Srv.DB.UpdateResultTriggerTime("trigger_start", c.TriggerStart, c.cfg.SN)
		} else if inputs.Inputs[c.Srv.config().IOTrigger-1] == '0' {
			// 释放trigger
			c.TriggerStop = time.Now()
			go c.Srv.DB.UpdateResultTriggerTime("trigger_stop", c.TriggerStop, c.cfg.SN)
		}

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

func (c *Controller) handleResult(result_data *ResultData, carve *minio.ControllerCurve) error {

	if ArrayContains(c.Srv.config().SkipJobs, result_data.JobID) {
		return nil
	}

	if c.dbController != nil {
		c.Srv.DB.UpdateTightning(c.dbController.Id, result_data.TightingID)
	}

	controllerResult := controller.ControllerResult{}
	controllerResult.NeedPushHmi = false
	controllerResult.TighteningID = result_data.TightingID

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

	targetID := result_data.VIN
	switch c.Srv.config().DataIndex {
	case 1:
		targetID = result_data.ID2
	case 2:
		targetID = result_data.ID3
	case 3:
		targetID = result_data.ID4
	}

	controllerResult.Workorder_ID, _ = strconv.ParseInt(targetID, 10, 64)
	controllerResult.NeedPushAiis = true

	controllerResult.GunSN = result_data.ToolSerialNumber
	controllerResult.Seq, controllerResult.Count = c.calBatch(controllerResult.Workorder_ID)

	if c.result_CURVE != nil {
		c.result_CURVE.CurveContent.Result = controllerResult.Result
		c.result_CURVE.CurveFile = result_data.ControllerName + result_data.ToolSerialNumber + result_data.TightingID
	}

	c.Srv.Parent.Handlers.Handle(&controllerResult, c.result_CURVE)

	return nil

}

// seq, count
func (c *Controller) calBatch(workorderID int64) (int, int) {
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

func (c *Controller) Start() {

	c.w = socket_writer.NewSocketWriter(c.cfg.RemoteIP, c)

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

func (c *Controller) getTighteningCount() {
	for {
		select {
		case <-time.After(c.getToolInfoPeriod):
			rev := GetVendorMid(c.Model(), MID_0040_TOOL_INFO_REQUEST)
			if rev == "" {
				continue
			}

			if c.Status() == controller.STATUS_OFFLINE {
				continue
			}
			req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", DEFAULT_MSG_END)
			c.Write([]byte(req))
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *Controller) ToolInfoReq() {
	rev := GetVendorMid(c.Model(), MID_0040_TOOL_INFO_REQUEST)
	if rev == "" {
		return
	}

	//defer c.Response.remove(MID_0040_TOOL_INFO_REQUEST)
	//c.Response.Add(MID_0040_TOOL_INFO_REQUEST, nil)

	req := GeneratePackage(MID_0040_TOOL_INFO_REQUEST, rev, "", DEFAULT_MSG_END)
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
}

func (c *Controller) GetPSetList() ([]int, error) {
	var psets []int

	rev := GetVendorMid(c.Model(), MID_0010_PSET_LIST_REQUEST)
	if rev == "" {
		return psets, errors.New("not supported")
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return psets, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0010_PSET_LIST_REQUEST)
	c.Response.Add(MID_0010_PSET_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0010_PSET_LIST_REQUEST, rev, "", DEFAULT_MSG_END)
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

func (c *Controller) GetPSetDetail(pset int) (PSetDetail, error) {
	var obj_pset_detail PSetDetail

	rev := GetVendorMid(c.Model(), MID_0012_PSET_DETAIL_REQUEST)
	if rev == "" {
		return obj_pset_detail, errors.New("not supported")
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_pset_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0012_PSET_DETAIL_REQUEST)
	c.Response.Add(MID_0012_PSET_DETAIL_REQUEST, nil)

	pset_detail := GeneratePackage(MID_0012_PSET_DETAIL_REQUEST, rev, fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)
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

func (c *Controller) GetJobList() ([]int, error) {
	var jobs []int
	rev := GetVendorMid(c.Model(), MID_0030_JOB_LIST_REQUEST)
	if rev == "" {
		return jobs, errors.New("not supported")
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return jobs, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0030_JOB_LIST_REQUEST)
	c.Response.Add(MID_0030_JOB_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0030_JOB_LIST_REQUEST, rev, "", DEFAULT_MSG_END)
	c.Write([]byte(psets_request))

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

	job_list := reply.(JobList)

	return job_list.jobs, nil
}

func (c *Controller) GetJobDetail(job int) (JobDetail, error) {
	var obj_job_detail JobDetail
	rev := GetVendorMid(c.Model(), MID_0032_JOB_DETAIL_REQUEST)
	if rev == "" {
		return obj_job_detail, errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_job_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0032_JOB_DETAIL_REQUEST)
	c.Response.Add(MID_0032_JOB_DETAIL_REQUEST, nil)

	job_detail := GeneratePackage(MID_0032_JOB_DETAIL_REQUEST, rev, fmt.Sprintf("%04d", job), DEFAULT_MSG_END)
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
		return obj_job_detail, errors.New(controller.ERR_CONTROLER_TIMEOUT)
	}

	switch v := reply.(type) {
	case string:
		return obj_job_detail, errors.New(v)

	case JobDetail:
		return reply.(JobDetail), nil

	default:
		return obj_job_detail, errors.New(controller.ERR_KNOWN)
	}

}

func (c *Controller) SolveOldResults() {

	if c.dbController == nil || c.dbController.LastID == "0" {
		return
	}

	c.Response.Add(MID_0064_OLD_SUBSCRIBE, MID_0064_OLD_SUBSCRIBE)
	if c.getOldResult(0) != nil {
		return
	}

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
	rev := GetVendorMid(c.Model(), MID_0001_START)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	//if c.Status() == controller.STATUS_OFFLINE {
	//	return errors.New("status offline")
	//}

	start := GeneratePackage(MID_0001_START, rev, "", DEFAULT_MSG_END)

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

	lenHeader := LEN_HEADER
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
				if len_msg < lenHeader-header_rest {
					//长度不够
					if header_rest == 0 {
						header_rest = lenHeader - len_msg
					} else {
						header_rest -= len_msg
					}
					header_buffer += msg[off : off+len_msg]
					break
				} else {
					//完整
					if header_rest == 0 {
						header_buffer = msg[off : off+lenHeader]
						off += lenHeader
					} else {
						//fmt.Printf("off:%d rest:%d msg:%s\n", off, off+header_rest, msg)
						if off < (off + header_rest) {
							header_buffer += msg[off : off+header_rest]
							off += header_rest
							header_rest = 0
						}
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

func (c *Controller) Parse(msg string) {
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

func (c *Controller) Close() error {

	for i := 0; i < 2; i++ {
		//两个协程需要关闭
		closed := make(chan struct{})
		c.closing <- closed

		<-closed
	}

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
	rev := GetVendorMid(c.Model(), MID_0064_OLD_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_last_result := GeneratePackage(MID_0064_OLD_SUBSCRIBE, rev, fmt.Sprintf("%010d", last_id), DEFAULT_MSG_END)

	c.Write([]byte(s_last_result))

	return nil
}

func (c *Controller) pset(pset int) error {
	rev := GetVendorMid(c.Model(), MID_0018_PSET)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0018_PSET, nil)
	defer c.Response.remove(MID_0018_PSET)

	s_pset := GeneratePackage(MID_0018_PSET, rev, fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)

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

func (c *Controller) ToolControl(enable bool) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_cmd := MID_0042_TOOL_DISABLE
	if enable {
		s_cmd = MID_0043_TOOL_ENABLE
	}

	rev := GetVendorMid(c.Model(), s_cmd)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	c.Response.Add(s_cmd, nil)
	defer c.Response.remove(s_cmd)

	sSend := GeneratePackage(s_cmd, rev, "", DEFAULT_MSG_END)

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
func (c *Controller) JobOff(off string) error {
	rev := GetVendorMid(c.Model(), MID_0130_JOB_OFF)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0130_JOB_OFF, nil)
	defer c.Response.remove(MID_0130_JOB_OFF)

	s_off := GeneratePackage(MID_0130_JOB_OFF, rev, off, DEFAULT_MSG_END)

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

func (c *Controller) jobSelect(job int) error {
	rev := GetVendorMid(c.Model(), MID_0038_JOB_SELECT)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0038_JOB_SELECT, nil)
	defer c.Response.remove(MID_0038_JOB_SELECT)

	s_job := GeneratePackage(MID_0038_JOB_SELECT, rev, fmt.Sprintf("%04d", job), DEFAULT_MSG_END)

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

func (c *Controller) IdentifierSet(str string) error {
	rev := GetVendorMid(c.Model(), MID_0150_IDENTIFIER_SET)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, rev, str, DEFAULT_MSG_END)

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) PSetBatchSet(pset int, batch int) error {
	rev := GetVendorMid(c.Model(), MID_0019_PSET_BATCH_SET)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s := fmt.Sprintf("%03d%02d", pset, batch)
	ide := GeneratePackage(MID_0019_PSET_BATCH_SET, rev, s, DEFAULT_MSG_END)

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) PSetBatchReset(pset int) error {
	rev := GetVendorMid(c.Model(), MID_0020_PSET_BATCH_RESET)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s := fmt.Sprintf("%03d", pset)
	ide := GeneratePackage(MID_0020_PSET_BATCH_RESET, rev, s, DEFAULT_MSG_END)

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) PSetSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0014_PSET_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0014_PSET_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) SelectorSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0250_SELECTOR_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0250_SELECTOR_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) JobInfoSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0034_JOB_INFO_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0034_JOB_INFO_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) IOInputSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0210_INPUT_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0210_INPUT_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(input))

	return nil
}

func (c *Controller) MultiSpindleResultSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0100_MULTI_SPINDLE_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0100_MULTI_SPINDLE_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(input))

	return nil
}

func (c *Controller) VinSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_0051_VIN_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0051_VIN_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(input))

	return nil
}

func (c *Controller) ResultSubcribe() error {
	rev := GetVendorMid(c.Model(), MID_0060_LAST_RESULT_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0060_LAST_RESULT_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) AlarmSubcribe() error {
	rev := GetVendorMid(c.Model(), MID_0070_ALARM_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	payload := GeneratePackage(MID_0070_ALARM_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(payload))

	return nil
}

func (c *Controller) CurveSubscribe() error {
	rev := GetVendorMid(c.Model(), MID_7408_LAST_CURVE_SUBSCRIBE)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_7408_LAST_CURVE_SUBSCRIBE, rev, "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) PSet(pset int, channel int, ex_info string, count int) (uint32, error) {

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

func (c *Controller) findIOByNo(no int, ios *[]IOStatus) (IOStatus, error) {
	for _, v := range *ios {
		if no == v.No {
			return v, nil
		}
	}

	return IOStatus{}, errors.New("not found")
}

func (c *Controller) IOSet(ios *[]IOStatus) error {
	rev := GetVendorMid(c.Model(), MID_0200_CONTROLLER_RELAYS)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
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

	s_io := GeneratePackage(MID_0200_CONTROLLER_RELAYS, rev, str_io, DEFAULT_MSG_END)

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

func (c *Controller) JobSet(id_info string, job int) error {

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

func (c *Controller) JobAbort() error {
	rev := GetVendorMid(c.Model(), MID_0127_JOB_ABORT)
	if rev == "" {
		return errors.New(controller.ERR_NOT_SUPPORTED)
	}

	if c.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	c.Response.Add(MID_0127_JOB_ABORT, nil)
	defer c.Response.remove(MID_0127_JOB_ABORT)

	s_job := GeneratePackage(MID_0127_JOB_ABORT, rev, "", DEFAULT_MSG_END)

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

func (c *Controller) Model() string {
	return c.model
}

func (c *Controller) SetModel(model string) {
	c.model = model
}

func (c *Controller) SetJob(r *tightening_device.JobSet) tightening_device.Reply {
	return tightening_device.Reply{}
}

func (c *Controller) SetPSet(r *tightening_device.PSetSet) tightening_device.Reply {
	rt := tightening_device.Reply{
		Result: 0,
		Msg:    "",
	}

	err := c.pset(r.PSet)
	if err != nil {
		rt.Result = -1
		rt.Msg = err.Error()
	}

	return rt
}

func (c *Controller) Enable(r *tightening_device.ToolEnable) tightening_device.Reply {
	return tightening_device.Reply{}
}

func (c *Controller) DeviceType() string {
	return "controller"
}

func (c *Controller) Children() []string {
	tools := []string{}
	for k, _ := range c.Tools() {
		tools = append(tools, k)
	}

	return tools
}

func (s *Controller) Data() interface{} {
	return nil
}

func (s *Controller) Config() interface{} {
	return nil
}
