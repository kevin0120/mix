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
)

type HandlerContext struct {
	cvi3Result          CVI3Result
	controllerCurve     ControllerCurve
	controllerResult    ControllerResult
	controllerCurveFile ControllerCurveFile
	dbCurve             storage.Curves
	resultIds           []int64
	wsResult            wsnotify.WSResult
	aiisResult          aiis.AIISResult
	aiisCurve           aiis.CURObject
}

type Handlers struct {
	AudiVw *Service
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, ctx *HandlerContext) error {
	h.AudiVw.diag.Debug("处理结果数据 ...")

	var needPushAiis bool = false

	var err error
	r, err := h.AudiVw.DB.GetResult(result.Result_id, 0)
	if err != nil {
		return err
	}

	workorder, err := h.AudiVw.DB.GetWorkorder(result.Workorder_ID)
	if err != nil {
		return err
	}

	// 保存结果
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

	h.AudiVw.diag.Debug("缓存结果到数据库 ...")

	if r.Count >= int(workorder.MaxRedoTimes) || r.Result == RESULT_OK {
		needPushAiis = true
		r.Stage = RESULT_STAGE_FINAL

		json.Unmarshal([]byte(workorder.ResultIDs), &ctx.resultIds)
		if r.ResultId == ctx.resultIds[len(ctx.resultIds)-1] {
			// 标记工单已完成
			workorder.Status = "finished"
			h.AudiVw.DB.UpdateWorkorder(&workorder)

		}
	}

	_, err = h.AudiVw.DB.UpdateResult(r)
	if err != nil {
		h.AudiVw.diag.Error("缓存结果失败", err)
		return nil
	} else {
		h.AudiVw.diag.Debug("缓存结果成功")
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

	if needPushAiis {

		// 结果推送AIIS
		if r.Result == RESULT_OK {
			ctx.aiisResult.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				ctx.aiisResult.One_time_pass = ODOO_RESULT_PASS
			} else {
				ctx.aiisResult.One_time_pass = ODOO_RESULT_FAIL
			}
		} else {
			ctx.aiisResult.Final_pass = ODOO_RESULT_FAIL
			ctx.aiisResult.One_time_pass = ODOO_RESULT_FAIL
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

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		ctx.aiisResult.CURObjects = []aiis.CURObject{}
		for _, v := range curves {
			ctx.aiisCurve.OP = v.Count
			ctx.aiisCurve.File = v.CurveFile
			ctx.aiisResult.CURObjects = append(ctx.aiisResult.CURObjects, ctx.aiisCurve)
		}

		h.PutResultToAIIS(ctx.aiisResult, r.ResultId)

	}

	return nil
}

func (h *Handlers) PutResultToAIIS(aiis_result aiis.AIISResult, r_id int64) error {
	h.AudiVw.diag.Debug("推送结果数据到AIIS")

	err := h.AudiVw.Aiis.PutResult(r_id, aiis_result)
	if err == nil {
		// 发送成功
		//db_result.HasUpload = true
		h.AudiVw.diag.Debug("推送AIIS成功，更新本地结果标识")
		_, err := h.AudiVw.DB.UpdateResultUpload(true, r_id)

		if err != nil {
			return err
		}
		return nil
	} else {
		h.AudiVw.diag.Error("推送AIIS失败", err)
		return err

	}
	return nil
}

// 处理波形数据
func (h *Handlers) handleCurve(curve *ControllerCurve, ctx *HandlerContext) error {
	h.AudiVw.diag.Debug("处理波形数据 ...")

	// 保存波形到数据库
	ctx.dbCurve.ResultID = curve.ResultID
	ctx.dbCurve.CurveData = curve.CurveData
	ctx.dbCurve.CurveFile = curve.CurveFile
	ctx.dbCurve.Count = curve.Count
	ctx.dbCurve.HasUpload = false

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

	// 保存波形到对象存储
	h.AudiVw.diag.Debug("保存波形数据到对象存储 ...")
	err = h.AudiVw.Minio.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		h.AudiVw.diag.Error("对象存储保存失败", err)
		return err
	} else {
		ctx.dbCurve.HasUpload = true
		h.AudiVw.diag.Debug("对象存储保存成功，更新本地结果标识")
		_, err = h.AudiVw.DB.UpdateCurve(&ctx.dbCurve)
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string, ctx *HandlerContext) {

	h.AudiVw.diag.Debug(fmt.Sprintf("收到结果数据:%s\n", msg))

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
