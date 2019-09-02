package openprotocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
	"strings"
)

func GetMidHandler(mid string) (MidHandler, error) {
	h, exist := MidHandlers[mid]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Handler Not Found, Mid: %s", mid))
	}

	return h, nil
}

var MidHandlers = map[string]MidHandler{
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

// 处理曲线
func handleMID_7410_LAST_CURVE(c *TighteningController, pkg *handlerPkg) error {
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

		var curve minio.ControllerCurve
		curve.CurveContent.CUR_M = Torque
		curve.CurveContent.CUR_W = Angle

		c.MID_7410_CURVE.Header = nil
		c.MID_7410_CURVE.Body = nil

		c.updateResult(nil, &curve)
		c.handleResultandClear()
	}

	return nil
}

// 处理结果
func handleMID_0061_LAST_RESULT(c *TighteningController, pkg *handlerPkg) error {
	result_data := ResultData{}
	err := result_data.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	return c.handleResult(&result_data, nil)
}

// 处理历史结果
func handleMID_0065_OLD_DATA(c *TighteningController, pkg *handlerPkg) error {
	result_data := ResultData{}
	err := result_data.DeserializeOld(pkg.Body)
	if err != nil {
		return err
	}

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

	return nil
}

// pset详情
func handleMID_0013_PSET_DETAIL_REPLY(c *TighteningController, pkg *handlerPkg) error {
	pset_detail := PSetDetail{}
	err := pset_detail.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	c.Response.update(MID_0012_PSET_DETAIL_REQUEST, pset_detail)

	return nil
}

// pset列表
func handleMID_0011_PSET_LIST_REPLY(c *TighteningController, pkg *handlerPkg) error {
	pset_list := PSetList{}
	err := pset_list.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	c.Response.update(MID_0010_PSET_LIST_REQUEST, pset_list)

	return nil
}

// job列表
func handleMID_0031_JOB_LIST_REPLY(c *TighteningController, pkg *handlerPkg) error {
	job_list := JobList{}
	err := job_list.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	c.Response.update(MID_0030_JOB_LIST_REQUEST, job_list)

	return nil
}

// job详情
func handleMID_0033_JOB_DETAIL_REPLY(c *TighteningController, pkg *handlerPkg) error {
	job_detail := JobDetail{}
	err := job_detail.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	c.Response.update(MID_0032_JOB_DETAIL_REQUEST, job_detail)

	return nil
}

// 请求错误
func handleMID_0004_CMD_ERR(c *TighteningController, pkg *handlerPkg) error {
	err_code := pkg.Body[4:6]
	c.Response.update(pkg.Body[0:4], request_errors[err_code])

	return nil
}

// 请求成功
func handleMID_0005_CMD_OK(c *TighteningController, pkg *handlerPkg) error {
	c.Response.update(pkg.Body, request_errors["00"])

	return nil
}

// job推送信息
func handleMID_0035_JOB_INFO(c *TighteningController, pkg *handlerPkg) error {
	job_info := JobInfo{}
	err := job_info.Deserialize(pkg.Body)
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

	str, _ := json.Marshal(inputs)
	c.Srv.WS.WSSendIOInput(string(str))

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
	err := ai.Deserialize(pkg.Body)
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
		c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

	case EVT_CONTROLLER_TOOL_DISCONNECT:
		c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
	}

	return nil
}

// 报警状态
func handleMID_0076_ALARM_STATUS(c *TighteningController, pkg *handlerPkg) error {
	var as AlarmStatus
	err := as.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	switch as.ErrorCode {
	case EVT_CONTROLLER_NO_ERR:
		c.UpdateToolStatus(controller.EVT_TOOL_CONNECTED)

	case EVT_CONTROLLER_TOOL_DISCONNECT:
		c.UpdateToolStatus(controller.EVT_TOOL_DISCONNECTED)
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

	if c.Status() == device.STATUS_OFFLINE {
		c.Response.update(MID_0040_TOOL_INFO_REQUEST, ti)
	} else {
		// 将数据通过api传给odoo
		if ti.ToolSN == "" {
			return errors.New("Tool Serial Number Is Empty String")
		}

		if ti.TotalTighteningCount == 0 || ti.CountSinLastService == 0 {
			//不需要尝试创建维修/标定单据
			return nil
		}

		go c.Srv.TryCreateMaintenance(ti) // 协程处理
	}

	return nil
}
