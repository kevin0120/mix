package openprotocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/ascii"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
	"strings"
	"time"
)

func GetMidHandler(mid string) (MidHandler, error) {
	h, exist := MidHandlers[mid]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Handler Not Found, Mid: %s", mid))
	}

	return h, nil
}

var MidHandlers = map[string]MidHandler{
	MID_9999_ALIVE:                handleMID_9999_ALIVE,
	MID_0002_START_ACK:            handleMID_0002_START_ACK,
	MID_0005_CMD_OK:               handleMID_0005_CMD_OK,
	MID_0004_CMD_ERR:              handleMID_0004_CMD_ERR,
	MID_7410_LAST_CURVE:           handleMID_7410_LAST_CURVE,
	MID_0061_LAST_RESULT:          handleMID_0061_LAST_RESULT,
	MID_0065_OLD_DATA:             handleMID_0065_OLD_DATA,
	MID_0013_PSET_DETAIL_REPLY:    handleMID_0013_PSET_DETAIL_REPLY,
	MID_0011_PSET_LIST_REPLY:      handleMID_0011_PSET_LIST_REPLY,
	MID_0031_JOB_LIST_REPLY:       handleMID_0031_JOB_LIST_REPLY,
	MID_0033_JOB_DETAIL_REPLY:     handleMID_0033_JOB_DETAIL_REPLY,
	MID_0035_JOB_INFO:             handleMID_0035_JOB_INFO,
	MID_0211_INPUT_MONITOR:        handleMID_0211_INPUT_MONITOR,
	MID_0101_MULTI_SPINDLE_RESULT: handleMID_0101_MULTI_SPINDLE_RESULT,
	MID_0052_VIN:                  handleMID_0052_VIN,
	MID_0071_ALARM:                handleMID_0071_ALARM,
	MID_0076_ALARM_STATUS:         handleMID_0076_ALARM_STATUS,
	MID_0041_TOOL_INFO_REPLY:      handleMID_0041_TOOL_INFO_REPLY,
}

type MidHandler func(controller *TighteningController, pkg *handlerPkg) error

func handleMID_9999_ALIVE(c *TighteningController, pkg *handlerPkg) error {
	return nil
}

func handleMID_0002_START_ACK(c *TighteningController, pkg *handlerPkg) error {
	seq := <-c.requestChannel
	c.Response.update(seq, request_errors["00"])

	// TODO
	go c.Subscribe()
	//go c.SolveOldResults()
	//go c.getTighteningCount()

	return nil
}

// 处理曲线
func handleMID_7410_LAST_CURVE(c *TighteningController, pkg *handlerPkg) error {
	//讲收到的曲线先做反序列化处理
	var curve = CurveBody{}
	err := ascii.Unmarshal(pkg.Body, &curve)
	if err != nil {
		c.diag.Error("ascii.Unmarshal", err)
	}
	if curve.ToolNumber == 0 {
		e := errors.New("收到的结果曲线数据不合法，未指定工具号")
		c.diag.Error("handleMID_7410_LAST_CURVE", e)
		return e
	}

	//_, ok := c.tempResultCurve[curve.ToolNumber]
	//结果曲线 判断本toolNumber是否收到过数据
	if _, ok := c.tempResultCurve[curve.ToolNumber]; !ok {
		c.tempResultCurve[curve.ToolNumber] = &tightening_device.TighteningCurve{}
	}
	//收到的数据进行解析并将结果加到临时的切片中，等待整条曲线接收完毕。
	torqueCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TorqueString), 64)
	angleCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.AngleString), 64)
	Torque, Angle := DataDecoding([]byte(curve.Data), torqueCoefficient, angleCoefficient, c.diag)

	c.tempResultCurve[curve.ToolNumber].CUR_M = append(c.tempResultCurve[curve.ToolNumber].CUR_M, Torque...)
	c.tempResultCurve[curve.ToolNumber].CUR_W = append(c.tempResultCurve[curve.ToolNumber].CUR_W, Angle...)
	//当本次数据为本次拧紧曲线的最后一次数据时
	if curve.Num == curve.Id {
		//若取到的点的数量大于协议解析出来该曲线的点数，多出的部分删掉，否则有多少发多少.
		if curve.MeasurePoints < len(c.tempResultCurve[curve.ToolNumber].CUR_M) {
			c.tempResultCurve[curve.ToolNumber].CUR_M = c.tempResultCurve[curve.ToolNumber].CUR_M[0:curve.MeasurePoints]
			c.tempResultCurve[curve.ToolNumber].CUR_W = c.tempResultCurve[curve.ToolNumber].CUR_W[0:curve.MeasurePoints]
		}

		//本次曲线全部解析完毕后,降临时存储的数据清空
		toolSN, err := c.findToolSNByChannel(curve.ToolNumber)
		if err != nil {
			return err
		}

		defer delete(c.tempResultCurve, curve.ToolNumber)
		c.tempResultCurve[curve.ToolNumber].ToolSN = toolSN
		c.tempResultCurve[curve.ToolNumber].UpdateTime = time.Now()
		c.toolDispatches[toolSN].curveDispatch.Dispatch(c.tempResultCurve[curve.ToolNumber])
	}
	return nil
}

// 处理结果
func handleMID_0061_LAST_RESULT(c *TighteningController, pkg *handlerPkg) error {
	result_data := ResultData{}
	err := ascii.Unmarshal(pkg.Body, &result_data)
	if err != nil {
		return err
	}

	return c.handleResult(&result_data, nil)
}

// 处理历史结果
func handleMID_0065_OLD_DATA(c *TighteningController, pkg *handlerPkg) error {
	result_data := ResultDataOld{}
	err := ascii.Unmarshal(pkg.Body, &result_data)
	if err != nil {
		return err
	}

	//flag := c.Response.get(MID_0064_OLD_SUBSCRIBE)

	//if flag != nil {
	//	defer c.Response.remove(MID_0064_OLD_SUBSCRIBE)
	//	c.Response.Add(MID_0065_OLD_DATA, result_data)
	//} else {
	//	// 处理历史数据
	//	// TODO: 通过工具获取psetdetail
	//	//pset_detail, err := c.GetPSetDetail(result_data.PSetID)
	//	//if err == nil {
	//	//	result_data.TorqueMin = pset_detail.TorqueMin
	//	//	result_data.TorqueMax = pset_detail.TorqueMax
	//	//	result_data.TorqueFinalTarget = pset_detail.TorqueTarget
	//	//	result_data.AngleMax = pset_detail.AngleMax
	//	//	result_data.AngleMin = pset_detail.AngleMin
	//	//	result_data.FinalAngleTarget = pset_detail.AngleTarget
	//	//}
	//
	//	return c.handleResult(&result_data, nil)
	//}

	return nil
}

// pset详情
func handleMID_0013_PSET_DETAIL_REPLY(c *TighteningController, pkg *handlerPkg) error {
	pset_detail, err := DeserializePSetDetail(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-c.requestChannel
	c.Response.update(seq, pset_detail)

	return nil
}

// pset列表
func handleMID_0011_PSET_LIST_REPLY(c *TighteningController, pkg *handlerPkg) error {
	pset_list := PSetList{}
	err := pset_list.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-c.requestChannel
	c.Response.update(seq, pset_list)

	return nil
}

// job列表
func handleMID_0031_JOB_LIST_REPLY(c *TighteningController, pkg *handlerPkg) error {
	job_list := JobList{}
	err := job_list.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-c.requestChannel
	c.Response.update(seq, job_list)

	return nil
}

// job详情
func handleMID_0033_JOB_DETAIL_REPLY(c *TighteningController, pkg *handlerPkg) error {
	jobDetaill, err := DeserializeJobDetail(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-c.requestChannel
	c.Response.update(seq, jobDetaill)

	return nil
}

// 请求错误
func handleMID_0004_CMD_ERR(c *TighteningController, pkg *handlerPkg) error {
	errCode := pkg.Body[4:6]

	seq := <-c.requestChannel
	c.Response.update(seq, errCode)

	return nil
}

// 请求成功
func handleMID_0005_CMD_OK(c *TighteningController, pkg *handlerPkg) error {
	seq := <-c.requestChannel
	c.Response.update(seq, request_errors["00"])

	return nil
}

// job推送信息
func handleMID_0035_JOB_INFO(c *TighteningController, pkg *handlerPkg) error {
	job_info := JobInfo{}
	err := ascii.Unmarshal(pkg.Body, &job_info)
	if err != nil {
		return err
	}

	// 加入判断，防止重复推送
	if job_info.JobStatus == JOB_INFO_NOT_COMPLETED &&
		job_info.JobBatchCounter == 0 &&
		job_info.JobCurrentStep > 0 &&
		job_info.JobTotalStep > 0 &&
		job_info.JobTighteningStatus == 0 {
		// 推送job选择信息

		job_select := wsnotify.WSJobSelect{
			JobID: job_info.JobID,
		}

		ws_str, _ := json.Marshal(job_select)
		c.Srv.diag.Debug(fmt.Sprintf("push job to hmi: %s", string(ws_str)))
		c.Srv.WS.WSSendJob(string(ws_str))
	}

	return nil

}

// 控制器输入变化
func handleMID_0211_INPUT_MONITOR(c *TighteningController, pkg *handlerPkg) error {
	inputs := IOMonitor{}
	err := inputs.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	inputs.ControllerSN = c.cfg.SN

	c.inputs = inputs.Inputs

	// 分发控制器输入状态
	c.GetDispatch(tightening_device.DISPATCH_IO).Dispatch(inputs.ToTighteningControllerInput())

	return nil
}

// 多轴结果
func handleMID_0101_MULTI_SPINDLE_RESULT(c *TighteningController, pkg *handlerPkg) error {
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
}

// 收到条码推送
func handleMID_0052_VIN(c *TighteningController, pkg *handlerPkg) error {
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

	return nil
}

// 报警信息
func handleMID_0071_ALARM(c *TighteningController, pkg *handlerPkg) error {
	var ai AlarmInfo
	err := ascii.Unmarshal(pkg.Body, &ai)
	if err != nil {
		return err
	}

	// 参见 项目管理,长安项目中文件:http://116.62.21.97/web#id=325&view_type=form&model=ir.attachment&active_id=3&menu_id=90
	// 第11页,错误代码:Tool calibration required:E305
	//if ai.ErrorCode == "E305" {
	//	// do nothing,当前未确认是否为这个错误代码
	//}

	switch ai.ErrorCode {
	case EVT_CONTROLLER_TOOL_CONNECT:
		//c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

	case EVT_CONTROLLER_TOOL_DISCONNECT:
		//c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
	}

	return nil
}

// 报警状态
func handleMID_0076_ALARM_STATUS(c *TighteningController, pkg *handlerPkg) error {
	var as = AlarmStatus{}
	err := ascii.Unmarshal(pkg.Body, &as)
	if err != nil {
		return err
	}

	switch as.ErrorCode {
	case EVT_CONTROLLER_NO_ERR:
		//c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

	case EVT_CONTROLLER_TOOL_DISCONNECT:
		//c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
	}

	return nil
}

// 工具状态(维护)
func handleMID_0041_TOOL_INFO_REPLY(c *TighteningController, pkg *handlerPkg) error {
	var ti ToolInfo
	err := ti.Deserialize(pkg.Body)

	if err != nil {
		return err
	}

	// 将数据通过api传给odoo
	if ti.ToolSN == "" {
		return errors.New("Tool Serial Number Is Empty String")
	}

	if ti.TotalTighteningCount == 0 || ti.CountSinLastService == 0 {
		//不需要尝试创建维修/标定单据
		return nil
	}

	go c.Srv.TryCreateMaintenance(ti) // 协程处理

	return nil
}
