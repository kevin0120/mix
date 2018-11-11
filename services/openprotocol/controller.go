package openprotocol

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/aiis"
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

	REPLY_TIMEOUT   = time.Duration(100 * time.Millisecond)
	MAX_REPLY_COUNT = 20
)

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}

type Controller struct {
	w                 *socket_writer.SocketWriter
	cfg               controller.ControllerConfig
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
	TriggerStart      time.Time
	TriggerStop       time.Time
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

		wsResults := []wsnotify.WSResult{}
		wsResult := wsnotify.WSResult{}
		for _, v := range ms.Spindles {
			wsResult.Result_id = 1
			wsResult.Count = 1
			wsResult.Result = v.Result
			wsResult.MI = v.Torque
			wsResult.WI = v.Angle
			wsResult.TI = 0
			wsResult.GroupSeq = 1
			wsResults = append(wsResults, wsResult)
		}

		//ws_str, _ := json.Marshal(wsResults)
		////c.Srv.WS.WSSendResult("1122334455667788", string(ws_str))

	case MID_0052_VIN:
		// 收到条码
		ids := DeserializeIDS(pkg.Body)

		barcode := wsnotify.WSScanner{
			Barcode: ids[0],
		}

		str, _ := json.Marshal(barcode)

		c.Srv.WS.WSSend(wsnotify.WS_EVENT_SCANNER, string(str))
	}
}

//var g_group

func (c *Controller) handleResult(result_data *ResultData) {

	result_data.VIN = strings.TrimSpace(result_data.VIN)

	// raw workorder id
	result_data.ID2 = strings.TrimSpace(result_data.ID2)

	result_data.ID3 = strings.TrimSpace(result_data.ID3)

	result_data.ID4 = strings.TrimSpace(result_data.ID4)

	c.Srv.DB.UpdateTightning(c.dbController.Id, result_data.TightingID)

	controllerResult := controller.ControllerResult{}
	c.Srv.diag.Info(fmt.Sprintf("vin:%s raw workorder_id:%s user_id:%s", result_data.VIN, result_data.ID2, result_data.ID3))

	id_info := result_data.VIN + result_data.ID2 + result_data.ID3
	if id_info == "" {
		c.Srv.diag.Error(id_info, errors.New("invalid id"))
		return
	}

	kvs := strings.Split(id_info, "-")

	//controllerResult.Batch = fmt.Sprintf("%d/%d", result_data.BatchCount, result_data.BatchSize)

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

	controllerResult.ExceptionReason = result_data.TighteningErrorStatus

	if kvs[0] != "auto" {
		if c.Srv.config().SkipJob == result_data.JobID {
			return
		}

		// 手动模式 [vin-hmisn-cartype-userid]

		// vin:cartype:hmisn
		//key := result_data.VIN + ":" + result_data.ID3 + ":" + result_data.ID2
		targetID := result_data.VIN
		switch c.Srv.config().DataIndex {
		case 1:
			targetID = result_data.ID2
		case 2:
			targetID = result_data.ID3
		case 3:
			targetID = result_data.ID4
		}

		raw_id, _ := strconv.ParseInt(targetID, 10, 64)
		db_workorder, err := c.Srv.DB.GetWorkorder(raw_id, true)
		db_result, err := c.Srv.DB.FindTargetResultForJobManual(raw_id)
		if err != nil {
			c.Srv.diag.Error("FindTargetResultForJobManual failed", err)
		}

		if err == nil {
			db_result.Spent = (c.TriggerStop.UnixNano() - c.TriggerStart.UnixNano()) / 1000000
		}
		//db_result := storage.Results{}
		controllerResult.Batch = fmt.Sprintf("%d/%d", db_result.Seq, db_workorder.MaxSeq)
		db_result.Batch = controllerResult.Batch
		db_result.GunSN = result_data.ToolSerialNumber

		if controllerResult.Result == storage.RESULT_OK {
			db_result.Stage = storage.RESULT_STAGE_FINAL
		}

		// 结果推送hmi
		wsResult := wsnotify.WSResult{}
		wsResult.Result_id = controllerResult.Result_id
		wsResult.Count = controllerResult.Count
		wsResult.Result = controllerResult.Result
		wsResult.MI = controllerResult.ResultValue.Mi
		wsResult.WI = controllerResult.ResultValue.Wi
		wsResult.TI = controllerResult.ResultValue.Ti
		//wsResult.Seq = db_result.Seq
		wsResult.GroupSeq = db_result.GroupSeq

		wsResults := []wsnotify.WSResult{}
		wsResults = append(wsResults, wsResult)
		ws_str, _ := json.Marshal(wsResults)

		c.Srv.diag.Debug(fmt.Sprintf("results:%s", string(ws_str)))

		c.Srv.WS.WSSendResult(db_workorder.HMISN, string(ws_str))

		// 结果缓存数据库
		c.Srv.Parent.Handlers.SaveResult(&controllerResult, &db_result, false)

		// 结果推送aiis
		aiisResult := aiis.AIISResult{}

		aiisResult.ID = db_result.Id
		aiisResult.MO_Model = db_workorder.MO_Model
		aiisResult.Batch = controllerResult.Batch
		aiisResult.Mode = db_workorder.Mode
		aiisResult.WorkcenterCode = db_workorder.HMISN
		aiisResult.ControllerSN = db_result.ControllerSN
		aiisResult.ToolSN = db_result.GunSN
		aiisResult.TighteningId, _ = strconv.ParseInt(result_data.TightingID, 10, 64)

		if controllerResult.ExceptionReason != "" {
			aiisResult.ExceptionReason = controllerResult.ExceptionReason + aiisResult.ExceptionReason
		}

		gun, err := c.Srv.DB.GetGun(result_data.ToolSerialNumber)
		if err != nil {
			odoo_gun, err := c.Srv.Odoo.GetGun(result_data.ToolSerialNumber)
			if err == nil {
				gun.GunID = odoo_gun.ID
				gun.Serial = odoo_gun.Serial
				c.Srv.DB.Store(gun)

			}
		}

		aiisResult.GunID = gun.GunID
		aiisResult.Control_date = time.Now().Format(time.RFC3339)

		aiisResult.Vin = db_workorder.Vin
		aiisResult.Measure_degree = controllerResult.ResultValue.Wi
		aiisResult.Measure_result = strings.ToLower(controllerResult.Result)
		aiisResult.Measure_t_don = controllerResult.ResultValue.Ti
		aiisResult.Measure_torque = controllerResult.ResultValue.Mi
		aiisResult.Op_time = controllerResult.Count
		aiisResult.Pset_m_max = controllerResult.PSetDefine.Mp
		aiisResult.Pset_m_min = controllerResult.PSetDefine.Mm
		aiisResult.Pset_m_target = controllerResult.PSetDefine.Ma
		aiisResult.Pset_m_threshold = controllerResult.PSetDefine.Ms
		aiisResult.Pset_strategy = controllerResult.PSetDefine.Strategy
		aiisResult.Pset_w_max = controllerResult.PSetDefine.Wp
		aiisResult.Pset_w_min = controllerResult.PSetDefine.Wm
		aiisResult.Pset_w_target = controllerResult.PSetDefine.Wa
		aiisResult.Pset_w_threshold = 1

		aiisResult.QualityState = controller.QUALITY_STATE_PASS
		if controllerResult.ExceptionReason != "" {
			aiisResult.QualityState = controller.QUALITY_STATE_EX
		}

		aiisResult.ProductID = db_workorder.ProductID
		aiisResult.WorkcenterID = db_workorder.WorkcenterID
		aiisResult.UserID = db_workorder.UserID
		//ks := strings.Split(result_data.ID4, ":")

		c.Srv.diag.Debug("推送结果数据到AIIS ...")

		err = c.Srv.Aiis.PutResult(0, aiisResult)
		if err == nil {
			c.Srv.diag.Debug("推送AIIS成功，更新本地结果标识")
		} else {
			c.Srv.diag.Error("推送AIIS失败", err)
		}

		return
	}

	//if len(kvs) == 2 {
	//	// job模式
	//
	//	controllerResult.Workorder_ID, _ = strconv.ParseInt(kvs[0], 10, 64)
	//	controllerResult.UserID, _ = strconv.ParseInt(kvs[1], 10, 64)
	//
	//	db_result, err := c.Srv.DB.FindTargetResultForJob(controllerResult.Workorder_ID)
	//	if err != nil {
	//		c.Srv.diag.Error("FindTargetResultForJob failed", err)
	//	}
	//
	//	controllerResult.Count = db_result.Count + 1
	//	controllerResult.Result_id = db_result.ResultId
	//
	//} else if len(kvs) == 3 {
	//	// pset模式
	//
	//	// 结果id-拧接次数-用户id
	//	controllerResult.Result_id, _ = strconv.ParseInt(kvs[0], 10, 64)
	//	controllerResult.Count, _ = strconv.Atoi(kvs[1])
	//	controllerResult.UserID, _ = strconv.ParseInt(kvs[2], 10, 64)
	//
	//	db_result, err := c.Srv.DB.GetResult(controllerResult.Result_id, 0)
	//	if err != nil {
	//		c.Srv.diag.Error("GetResult failed", err)
	//	}
	//
	//	controllerResult.Workorder_ID = db_result.WorkorderID
	//} else {
	//	c.Srv.diag.Error(id_info, errors.New("invalid id"))
	//	return
	//}
	//
	//c.Srv.Parent.Handle(controllerResult, nil)
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

	//c.JobOff("1")
	c.PSetSubscribe()
	//c.CurveSubscribe()
	c.SelectorSubscribe()
	c.ResultSubcribe()
	c.JobInfoSubscribe()
	c.IOInputSubscribe()
	c.MultiSpindleResultSubscribe()
	c.VinSubscribe()
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

	defer c.Response.remove(MID_0010_PSET_LIST_REQUEST)
	c.Response.Add(MID_0010_PSET_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0010_PSET_LIST_REQUEST, "001", "", DEFAULT_MSG_END)
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

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_pset_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0012_PSET_DETAIL_REQUEST)
	c.Response.Add(MID_0012_PSET_DETAIL_REQUEST, nil)

	pset_detail := GeneratePackage(MID_0012_PSET_DETAIL_REQUEST, "002", fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)
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
	if c.Status() == controller.STATUS_OFFLINE {
		return jobs, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0030_JOB_LIST_REQUEST)
	c.Response.Add(MID_0030_JOB_LIST_REQUEST, nil)

	psets_request := GeneratePackage(MID_0030_JOB_LIST_REQUEST, "002", "", DEFAULT_MSG_END)
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

	if c.Status() == controller.STATUS_OFFLINE {
		return obj_job_detail, errors.New(controller.STATUS_OFFLINE)
	}

	defer c.Response.remove(MID_0032_JOB_DETAIL_REQUEST)
	c.Response.Add(MID_0032_JOB_DETAIL_REQUEST, nil)

	job_detail := GeneratePackage(MID_0032_JOB_DETAIL_REQUEST, "003", fmt.Sprintf("%04d", job), DEFAULT_MSG_END)
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

	c.Response.Add(MID_0018_PSET, nil)
	defer c.Response.remove(MID_0018_PSET)

	s_pset := GeneratePackage(MID_0018_PSET, "001", fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)

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

	c.Response.Add(s_cmd, nil)
	defer c.Response.remove(s_cmd)

	s_send := GeneratePackage(s_cmd, "001", "", DEFAULT_MSG_END)

	c.Write([]byte(s_send))

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

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}

// 0: set 1: reset
func (c *Controller) JobOff(off string) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0130_JOB_OFF, nil)
	defer c.Response.remove(MID_0130_JOB_OFF)

	s_off := GeneratePackage(MID_0130_JOB_OFF, "001", off, DEFAULT_MSG_END)

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
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	c.Response.Add(MID_0038_JOB_SELECT, nil)
	defer c.Response.remove(MID_0038_JOB_SELECT)

	s_job := GeneratePackage(MID_0038_JOB_SELECT, "002", fmt.Sprintf("%04d", job), DEFAULT_MSG_END)

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
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, "001", str, DEFAULT_MSG_END)

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) PSetBatchSet(pset int, batch int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s := fmt.Sprintf("%03d%02d", pset, batch)
	ide := GeneratePackage(MID_0019_PSET_BATCH_SET, "001", s, DEFAULT_MSG_END)

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) PSetBatchReset(pset int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s := fmt.Sprintf("%03d", pset)
	ide := GeneratePackage(MID_0020_PSET_BATCH_RESET, "001", s, DEFAULT_MSG_END)

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

func (c *Controller) IOInputSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0210_INPUT_SUBSCRIBE, "001", "", DEFAULT_MSG_END)

	c.Write([]byte(input))

	return nil
}

func (c *Controller) MultiSpindleResultSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0100_MULTI_SPINDLE_SUBSCRIBE, "000", "", DEFAULT_MSG_END)

	c.Write([]byte(input))

	return nil
}

func (c *Controller) VinSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	input := GeneratePackage(MID_0051_VIN_SUBSCRIBE, "002", "", DEFAULT_MSG_END)

	c.Write([]byte(input))

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

func (c *Controller) PSet(pset int, channel int, ex_info string, count int) (uint32, error) {
	// 设定结果标识

	if c.Mode.Load().(string) != MODE_PSET {
		return 0, errors.New("current mode is not pset")
	}

	// 次数控制
	c.PSetBatchSet(pset, count)

	// 结果id-拧接次数-用户id
	err := c.IdentifierSet(ex_info)
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

func (c *Controller) findIOByNo(no int, ios *[]IOStatus) (IOStatus, error) {
	for _, v := range *ios {
		if no == v.No {
			return v, nil
		}
	}

	return IOStatus{}, errors.New("not found")
}

func (c *Controller) IOSet(ios *[]IOStatus) error {
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

	s_io := GeneratePackage(MID_0200_CONTROLLER_RELAYS, "001", str_io, DEFAULT_MSG_END)

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

func (c *Controller) JobAbort() error {

	if c.Mode.Load().(string) != MODE_JOB {
		return errors.New("current mode is not job")
	}

	c.Response.Add(MID_0127_JOB_ABORT, nil)
	defer c.Response.remove(MID_0127_JOB_ABORT)

	s_job := GeneratePackage(MID_0127_JOB_ABORT, "001", "", DEFAULT_MSG_END)

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

	s_reply := reply.(string)
	if s_reply != request_errors["00"] {
		return errors.New(s_reply)
	}

	return nil
}
