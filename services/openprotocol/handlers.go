package openprotocol

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/scanner"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils/ascii"
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
	MID_9999_ALIVE:                handleMid9999Alive,
	MID_0002_START_ACK:            handleMid0002StartAck,
	MID_0005_CMD_OK:               handleMid0005CmdOk,
	MID_0004_CMD_ERR:              handleMid0004CmdErr,
	MID_7410_LAST_CURVE:           handleMid7410LastCurve,
	MID_0061_LAST_RESULT:          handleMid0061LastResult,
	MID_0065_OLD_DATA:             handleMid0065OldData,
	MID_0013_PSET_DETAIL_REPLY:    handleMid0013PsetDetailReply,
	MID_0011_PSET_LIST_REPLY:      handleMid0011PsetListReply,
	MID_0031_JOB_LIST_REPLY:       handleMid0031JobListReply,
	MID_0033_JOB_DETAIL_REPLY:     handleMid0033JobDetailReply,
	MID_0035_JOB_INFO:             handleMid0035JobInfo,
	MID_0211_INPUT_MONITOR:        handleMid0211InputMonitor,
	MID_0101_MULTI_SPINDLE_RESULT: handleMid0101MultiSpindleResult,
	MID_0052_VIN:                  handleMid0052Vin,
	MID_0071_ALARM:                handleMid0071Alarm,
	MID_0076_ALARM_STATUS:         handleMid0076AlarmStatus,
	MID_0041_TOOL_INFO_REPLY:      handleMid0041ToolInfoReply,
}

type MidHandler func(controller *TighteningController, pkg *handlerPkg) error

func handleMid9999Alive(c *TighteningController, pkg *handlerPkg) error {
	return nil
}

func handleMid0002StartAck(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	seq := <-client.requestChannel
	client.response.update(seq, request_errors["00"])

	go c.ProcessSubscribeControllerInfo(pkg.SN)

	return nil
}

// 处理曲线
func handleMid7410LastCurve(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	//讲收到的曲线先做反序列化处理
	var curve = CurveBody{}
	err := ascii.Unmarshal(pkg.Body, &curve)
	if err != nil {
		c.diag.Error("ascii.Unmarshal", err)
	}
	if curve.ToolChannelNumber == 0 {
		e := errors.New("收到的结果曲线数据不合法，未指定工具号")
		c.diag.Error("handleMID_7410_LAST_CURVE", e)
		return e
	}

	torqueCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TorqueString), 64)
	angleCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.AngleString), 64)
	Torque, Angle := c.CurveDataDecoding([]byte(curve.Data), torqueCoefficient, angleCoefficient, c.diag)

	client.tempResultCurve.CUR_M = append(client.tempResultCurve.CUR_M, Torque...)
	client.tempResultCurve.CUR_W = append(client.tempResultCurve.CUR_W, Angle...)
	//当本次数据为本次拧紧曲线的最后一次数据时
	if curve.Num == curve.Id {
		//若取到的点的数量大于协议解析出来该曲线的点数，多出的部分删掉，否则有多少发多少.
		if curve.MeasurePoints < len(client.tempResultCurve.CUR_M) {
			client.tempResultCurve.CUR_M = client.tempResultCurve.CUR_M[0:curve.MeasurePoints]
			client.tempResultCurve.CUR_W = client.tempResultCurve.CUR_W[0:curve.MeasurePoints]
		}

		//本次曲线全部解析完毕后,降临时存储的数据清空
		tool, err := c.getInstance().GetToolViaChannel(curve.ToolChannelNumber)
		if err != nil {
			return err
		}

		sn := tool.SerialNumber()

		//defer delete(client.tempResultCurve, curve.ToolChannelNumber)
		client.tempResultCurve.ToolSN = sn
		client.tempResultCurve.UpdateTime = time.Now()
		c.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherCurve), client.tempResultCurve)
		// dispatch完后创建新的缓存拧紧曲线
		client.tempResultCurve = tightening_device.NewTighteningCurve()
	}
	return nil
}

// 处理结果
func handleMid0061LastResult(c *TighteningController, pkg *handlerPkg) error {
	resultData := ResultData{}
	err := ascii.Unmarshal(pkg.Body, &resultData)
	if err != nil {
		return err
	}

	tighteningResult := resultData.ToTighteningResult()
	return c.handleResult(tighteningResult)
}

// 处理历史结果
func handleMid0065OldData(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	resultData := ResultDataOld{}
	err := ascii.Unmarshal(pkg.Body, &resultData)
	if err != nil {
		return err
	}

	seq := <-client.requestChannel
	client.response.update(seq, resultData.ToTighteningResult())

	return nil
}

// pset详情
func handleMid0013PsetDetailReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	psetDetail, err := DeserializePSetDetail(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-client.requestChannel
	client.response.update(seq, psetDetail)

	return nil
}

// pset列表
func handleMid0011PsetListReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	psetList := PSetList{}
	err := psetList.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-client.requestChannel
	client.response.update(seq, psetList)

	return nil
}

// job列表
func handleMid0031JobListReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	jobList := JobList{}
	err := jobList.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-client.requestChannel
	client.response.update(seq, jobList)

	return nil
}

// job详情
func handleMid0033JobDetailReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	jobDetaill, err := DeserializeJobDetail(pkg.Body)
	if err != nil {
		return err
	}

	seq := <-client.requestChannel
	client.response.update(seq, jobDetaill)

	return nil
}

// 请求错误
func handleMid0004CmdErr(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	errCode := pkg.Body[4:6]

	seq := <-client.requestChannel
	client.response.update(seq, request_errors[errCode])

	return nil
}

// 请求成功
func handleMid0005CmdOk(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	seq := <-client.requestChannel
	client.response.update(seq, request_errors["00"])

	return nil
}

// job推送信息
func handleMid0035JobInfo(c *TighteningController, pkg *handlerPkg) error {
	jobInfo := JobInfo{}
	err := ascii.Unmarshal(pkg.Body, &jobInfo)
	if err != nil {
		return err
	}

	// 加入判断，防止重复推送
	if jobInfo.JobStatus == JobInfoNotCompleted &&
		jobInfo.JobBatchCounter == 0 &&
		jobInfo.JobCurrentStep > 0 &&
		jobInfo.JobTotalStep > 0 &&
		jobInfo.JobTighteningStatus == 0 {
		// 推送job选择信息

		jobSelect := tightening_device.JobInfo{
			Job: jobInfo.JobID,
		}

		c.doDispatch(dispatcherbus.DispatcherJob, jobSelect)
	}

	return nil

}

// 控制器输入变化
func handleMid0211InputMonitor(c *TighteningController, pkg *handlerPkg) error {
	inputs := IOMonitor{}
	err := inputs.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	inputs.ControllerSN = c.deviceConf.SN

	c.inputs = inputs.Inputs

	c.NotifyIOContact(io.IoTypeInput, c.inputs)

	return nil
}

// 多轴结果
func handleMid0101MultiSpindleResult(c *TighteningController, pkg *handlerPkg) error {
	ms := MultiSpindleResult{}
	ms.Deserialize(pkg.Body)

	return nil
}

// 收到条码推送
func handleMid0052Vin(c *TighteningController, pkg *handlerPkg) error {
	ids := DeserializeIDS(pkg.Body)

	bc := ""
	for _, v := range c.ProtocolService.config().VinIndex {
		if v < 0 || v > (MaxIdsNum-1) {
			continue
		}

		bc += ids[v]
	}
	ss := scanner.ScannerRead{
		Src:     tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER,
		SN:      c.deviceConf.SN,
		Barcode: bc,
	}
	c.doDispatch(dispatcherbus.DispatcherScannerData, ss)

	return nil
}

// 报警信息
func handleMid0071Alarm(c *TighteningController, pkg *handlerPkg) error {
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
	case EvtControllerToolConnect:
		c.getInstance().UpdateToolStatus(pkg.SN, device.BaseDeviceStatusOnline)

	case EvtControllerToolDisconnect:
		c.getInstance().UpdateToolStatus(pkg.SN, device.BaseDeviceStatusOffline)
	}

	return nil
}

// 报警状态
func handleMid0076AlarmStatus(c *TighteningController, pkg *handlerPkg) error {
	var as = AlarmStatus{}
	err := ascii.Unmarshal(pkg.Body, &as)
	if err != nil {
		return err
	}

	switch as.ErrorCode {
	case EvtControllerNoErr:
		c.getInstance().UpdateToolStatus(pkg.SN, device.BaseDeviceStatusOnline)

	case EvtControllerToolDisconnect:
		c.getInstance().UpdateToolStatus(pkg.SN, device.BaseDeviceStatusOffline)
	}

	return nil
}

// 工具状态(维护)
func handleMid0041ToolInfoReply(c *TighteningController, pkg *handlerPkg) error {
	var ti ToolInfo
	err := ti.Deserialize(pkg.Body)

	if err != nil {
		return err
	}

	if ti.ToolSN == "" {
		return errors.New("Tool Serial Number Is Empty String ")
	}

	if ti.TotalTighteningCount == 0 || ti.CountSinLastService == 0 {
		//不需要尝试创建维修/标定单据
		return nil
	}

	c.doDispatch(dispatcherbus.DispatcherToolMaintenance, ti.ToMaintenanceInfo())
	return nil
}
