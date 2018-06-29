package audi_vw

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"strings"
	"time"
)

const (
	ODOO_RESULT_PASS = "pass"
	ODOO_RESULT_FAIL = "fail"

	QUALITY_STATE_PASS = "pass"
	QUALITY_STATE_FAIL = "fail"
	QUALITY_STATE_EX = "exception"
)

type HandlerContext struct {
	cvi3Result          CVI3Result
	controllerCurve     ControllerCurve
	controllerResult    ControllerResult
	controllerCurveFile ControllerCurveFile
	dbCurve             storage.Curves
	//resultIds           []int64
	wsResult            wsnotify.WSResult
	aiisResult          aiis.AIISResult
	aiisCurve           aiis.CURObject
}

type Handlers struct {
	AudiVw *Service
}

func (h *Handlers) PushAiis(needPush bool, r *storage.Results, workorder *storage.Workorders, result *ControllerResult, ctx *HandlerContext) error {

	defer func() {
		h.AudiVw.diag.Debug("缓存结果到数据库 ...")
		_, err := h.AudiVw.DB.UpdateResult(r)
		if err != nil {
			h.AudiVw.diag.Error("缓存结果失败", err)
		} else {
			h.AudiVw.diag.Debug("缓存结果成功")
		}
	}()

	if needPush {

		// 结果推送AIIS
		if r.Result == RESULT_OK {
			ctx.aiisResult.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				ctx.aiisResult.One_time_pass = ODOO_RESULT_PASS
			} else {
				ctx.aiisResult.One_time_pass = ODOO_RESULT_FAIL
			}

			if (result.ResultValue.Mi >= r.ToleranceMin && result.ResultValue.Mi <= r.ToleranceMax) &&
				(result.ResultValue.Wi >= r.ToleranceMinDegree && result.ResultValue.Wi <= r.ToleranceMaxDegree) {
				ctx.aiisResult.QualityState = QUALITY_STATE_PASS
				ctx.aiisResult.ExceptionReason = ""
			} else {
				ctx.aiisResult.QualityState = QUALITY_STATE_EX
				ctx.aiisResult.ExceptionReason = QUALITY_STATE_EX
			}

		} else {
			ctx.aiisResult.Final_pass = ODOO_RESULT_FAIL
			ctx.aiisResult.One_time_pass = ODOO_RESULT_FAIL

			if (result.ResultValue.Mi >= r.ToleranceMin && result.ResultValue.Mi <= r.ToleranceMax) &&
				(result.ResultValue.Wi >= r.ToleranceMinDegree && result.ResultValue.Wi <= r.ToleranceMaxDegree) {

				ctx.aiisResult.QualityState = QUALITY_STATE_EX
				ctx.aiisResult.ExceptionReason = QUALITY_STATE_EX
			} else {
				ctx.aiisResult.QualityState = QUALITY_STATE_FAIL
				ctx.aiisResult.ExceptionReason = ""
			}

		}

		ctx.aiisResult.Control_date = r.UpdateTime.Format(time.RFC3339)

		ctx.aiisResult.Measure_degree = result.ResultValue.Wi
		ctx.aiisResult.Measure_result = strings.ToLower(result.Result)
		ctx.aiisResult.Measure_t_don = result.ResultValue.Ti
		ctx.aiisResult.Measure_torque = result.ResultValue.Mi
		ctx.aiisResult.Op_time = result.Count
		ctx.aiisResult.Pset_m_max = result.PSetDefine.Mp
		ctx.aiisResult.Pset_m_min = result.PSetDefine.Mm
		ctx.aiisResult.Pset_m_target = result.PSetDefine.Ma
		ctx.aiisResult.Pset_m_threshold = result.PSetDefine.Ms
		ctx.aiisResult.Pset_strategy = result.PSetDefine.Strategy
		ctx.aiisResult.Pset_w_max = result.PSetDefine.Wp
		ctx.aiisResult.Pset_w_min = result.PSetDefine.Wm
		ctx.aiisResult.Pset_w_target = result.PSetDefine.Wa
		ctx.aiisResult.Pset_w_threshold = 1
		ctx.aiisResult.UserID = result.UserID
		ctx.aiisResult.Seq = r.Seq

		// mo相关
		ctx.aiisResult.MO_AssemblyLine = workorder.MO_AssemblyLine
		ctx.aiisResult.MO_EquipemntName = workorder.MO_EquipemntName
		ctx.aiisResult.MO_FactoryName = workorder.MO_FactoryName
		ctx.aiisResult.MO_Pin = workorder.MO_Pin
		ctx.aiisResult.MO_Pin_check_code = workorder.MO_Pin_check_code
		ctx.aiisResult.MO_Year = workorder.MO_Year
		ctx.aiisResult.MO_Lnr = workorder.MO_Lnr
		ctx.aiisResult.MO_NutNo = r.NutNo
		ctx.aiisResult.MO_Model = workorder.MO_Model

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		for _, v := range curves {
			ctx.aiisCurve.OP = v.Count
			ctx.aiisCurve.File = v.CurveFile
			ctx.aiisResult.CURObjects = append(ctx.aiisResult.CURObjects, ctx.aiisCurve)
		}

		h.AudiVw.diag.Debug("推送结果数据到AIIS ...")

		err = h.AudiVw.Aiis.PutResult(r.ResultId, ctx.aiisResult)
		if err == nil {
			r.HasUpload = true
			h.AudiVw.diag.Debug("推送AIIS成功，更新本地结果标识")
		} else {
			h.AudiVw.diag.Error("推送AIIS失败", err)
		}
	}

	return nil
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, ctx *HandlerContext) error {
	h.AudiVw.diag.Debug("处理结果数据 ...")

	needPushAiis := false

	r, err := h.AudiVw.DB.GetResult(result.Result_id, 0)
	if err != nil {
		return err
	}

	workorder, err := h.AudiVw.DB.GetWorkorder(result.Workorder_ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Local")
	r.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", result.Dat, loc)
	r.Result = result.Result
	r.Count = result.Count
	r.HasUpload = false
	r.ControllerSN = result.Controller_SN
	r.UserID = result.UserID
	s_value, _ := json.Marshal(result.ResultValue)
	s_pset, _ := json.Marshal(result.PSetDefine)

	r.ResultValue = string(s_value)
	r.PSetDefine = string(s_pset)

	if r.Count >= int(r.MaxRedoTimes) || r.Result == RESULT_OK {
		needPushAiis = true
		r.Stage = RESULT_STAGE_FINAL

		if r.ResultId == workorder.LastResultID {
			h.AudiVw.diag.Debug("工单已完成")
			workorder.Status = "done"
			h.AudiVw.DB.UpdateWorkorder(&workorder)
		}
	}

	// 结果推送hmi
	ctx.wsResult.Result_id = result.Result_id
	ctx.wsResult.Count = result.Count
	ctx.wsResult.Result = result.Result
	ctx.wsResult.MI = result.ResultValue.Mi
	ctx.wsResult.WI = result.ResultValue.Wi
	ctx.wsResult.TI = result.ResultValue.Ti
	ws_str, _ := json.Marshal(ctx.wsResult)

	h.AudiVw.diag.Debug("Websocket推送结果到HMI")

	h.AudiVw.WS.WSSendResult(workorder.HMISN, string(ws_str))

	go h.PushAiis(needPushAiis,&r, &workorder, result, ctx)

	return nil
}

// 处理波形数据
func (h *Handlers) handleCurve(curve *ControllerCurve, ctx *HandlerContext) error {
	h.AudiVw.diag.Debug("处理波形数据 ...")

	ctx.dbCurve.ResultID = curve.ResultID
	ctx.dbCurve.CurveData = curve.CurveData
	ctx.dbCurve.CurveFile = curve.CurveFile
	ctx.dbCurve.Count = curve.Count
	ctx.dbCurve.HasUpload = false
	loc, _ := time.LoadLocation("Local")
	ctx.dbCurve.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", ctx.controllerResult.Dat, loc)

	// 保存波形到对象存储
	h.AudiVw.diag.Debug("保存波形数据到对象存储 ...")
	err := h.AudiVw.Minio.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		h.AudiVw.diag.Error("对象存储保存失败", err)
		return err
	} else {
		ctx.dbCurve.HasUpload = true
		h.AudiVw.diag.Debug("波形存储成功")
	}

	// 保存波形到数据库
	exist, err := h.AudiVw.DB.CurveExist(&ctx.dbCurve)
	if err != nil {
		return err
	} else {
		h.AudiVw.diag.Debug("缓存波形数据到数据库 ...")
		if exist {
			_, err := h.AudiVw.DB.UpdateCurve(&ctx.dbCurve)
			if err != nil {
				h.AudiVw.diag.Error("缓存波形失败", err)
				return err
			}
		} else {
			err := h.AudiVw.DB.Store(ctx.dbCurve)
			if err != nil {
				h.AudiVw.diag.Error("缓存波形失败", err)
				return err
			}
		}

		h.AudiVw.diag.Debug("缓存波形成功")
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string, ctx *HandlerContext) {

	//h.AudiVw.diag.Debug(fmt.Sprintf("收到结果数据:%s\n", msg))

	err := xml.Unmarshal([]byte(msg), &ctx.cvi3Result)
	if err != nil {
		h.AudiVw.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
		return
	}

	//结果数据
	XML2Result(&ctx.cvi3Result, &ctx.controllerResult)

	// 波形文件
	XML2Curve(&ctx.cvi3Result, &ctx.controllerCurveFile)

	sCurvedata, _ := json.Marshal(ctx.controllerCurveFile)
	ctx.controllerCurve.CurveData = string(sCurvedata)
	ctx.controllerCurve.Count = ctx.controllerResult.Count
	ctx.controllerCurve.CurveFile = ctx.controllerResult.CurFile
	ctx.controllerCurve.ResultID = ctx.controllerResult.Result_id

	e := h.handleCurve(&ctx.controllerCurve, ctx)
	if e == nil {
		h.handleResult(&ctx.controllerResult, ctx)
	} else {
		h.AudiVw.diag.Error("handleCurve err", err)
	}

}
