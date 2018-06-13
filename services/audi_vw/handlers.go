package audi_vw

import (
	"fmt"
	"encoding/json"
	"github.com/masami10/rush/services/storage"
	"time"
	"strings"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/wsnotify"
	"encoding/xml"
)

const(
	ODOO_RESULT_PASS = "pass"
	ODOO_RESULT_FAIL = "fail"
)

type HandlerContext struct {
	cvi3_result CVI3Result
	controller_curve ControllerCurve
	controller_result ControllerResult
	controller_curve_file ControllerCurveFile
	db_curve	storage.Curves
	result_ids []int64
	ws_result wsnotify.WSResult
	aiis_result aiis.AIISResult
	aiis_curve aiis.CURObject
}

type Handlers struct {
	AudiVw	*Service
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, ctx *HandlerContext) (error) {
	h.AudiVw.diag.Debug("处理结果数据 ...")

	var need_push_aiis bool = false

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
		need_push_aiis = true
		r.Stage = RESULT_STAGE_FINAL

		json.Unmarshal([]byte(workorder.ResultIDs), &ctx.result_ids)
		if r.ResultId == ctx.result_ids[len(ctx.result_ids) - 1] {
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
	ctx.ws_result.Result_id = result.Result_id
	ctx.ws_result.Count = result.Count
	ctx.ws_result.Result = result.Result
	ctx.ws_result.MI = result.ResultValue.Mi
	ctx.ws_result.WI = result.ResultValue.Wi
	ctx.ws_result.TI = result.ResultValue.Ti
	ws_str, _ := json.Marshal(ctx.ws_result)

	h.AudiVw.diag.Debug("Websocket推送结果到HMI")

	h.AudiVw.WS.WSSendResult(workorder.HMISN, string(ws_str))

	if need_push_aiis {

		// 结果推送AIIS
		if r.Result == RESULT_OK {
			ctx.aiis_result.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				ctx.aiis_result.One_time_pass = ODOO_RESULT_PASS
			} else {
				ctx.aiis_result.One_time_pass = ODOO_RESULT_FAIL
			}
		} else {
			ctx.aiis_result.Final_pass = ODOO_RESULT_FAIL
			ctx.aiis_result.One_time_pass = ODOO_RESULT_FAIL
		}

		ctx.aiis_result.Control_date = r.UpdateTime.Format(time.RFC3339)

		ctx.aiis_result.Measure_degree = result.ResultValue.Wi
		ctx.aiis_result.Measure_result = strings.ToLower(result.Result)
		ctx.aiis_result.Measure_t_don = result.ResultValue.Ti
		ctx.aiis_result.Measure_torque = result.ResultValue.Mi
		ctx.aiis_result.Op_time = result.Count
		ctx.aiis_result.Pset_m_max = result.PSetDefine.Mp
		ctx.aiis_result.Pset_m_min = result.PSetDefine.Mm
		ctx.aiis_result.Pset_m_target = result.PSetDefine.Ma
		ctx.aiis_result.Pset_m_threshold = result.PSetDefine.Ms
		ctx.aiis_result.Pset_strategy = result.PSetDefine.Strategy
		ctx.aiis_result.Pset_w_max = result.PSetDefine.Wp
		ctx.aiis_result.Pset_w_min = result.PSetDefine.Wm
		ctx.aiis_result.Pset_w_target = result.PSetDefine.Wa
		ctx.aiis_result.Pset_w_threshold = 1
		ctx.aiis_result.UserID = result.UserID

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		ctx.aiis_result.CURObjects = []aiis.CURObject{}
		for _, v := range curves {
			ctx.aiis_curve.OP = v.Count
			ctx.aiis_curve.File = v.CurveFile
			ctx.aiis_result.CURObjects = append(ctx.aiis_result.CURObjects, ctx.aiis_curve)
		}


		h.PutResultToAIIS(ctx.aiis_result, r.ResultId)

	}

	return nil
}

func (h *Handlers) PutResultToAIIS(aiis_result aiis.AIISResult, r_id int64) error{
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
func (h *Handlers) handleCurve(curve *ControllerCurve, ctx *HandlerContext) (error) {
	h.AudiVw.diag.Debug("处理波形数据 ...")

	// 保存波形到数据库
	ctx.db_curve.ResultID = curve.ResultID
	ctx.db_curve.CurveData = curve.CurveData
	ctx.db_curve.CurveFile = curve.CurveFile
	ctx.db_curve.Count = curve.Count
	ctx.db_curve.HasUpload = false

	exist, err := h.AudiVw.DB.CurveExist(&ctx.db_curve)
	if err != nil {
		return err
	} else {
		h.AudiVw.diag.Debug("缓存波形数据到数据库 ...")
		if exist {
			_, err := h.AudiVw.DB.UpdateCurve(&ctx.db_curve)
			if err != nil {
				h.AudiVw.diag.Error("缓存波形失败", err)
				return err
			}
		} else {
			err := h.AudiVw.DB.Store(ctx.db_curve)
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
		ctx.db_curve.HasUpload = true
		h.AudiVw.diag.Debug("对象存储保存成功，更新本地结果标识")
		_, err = h.AudiVw.DB.UpdateCurve(&ctx.db_curve)
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string, ctx *HandlerContext) {

	h.AudiVw.diag.Debug(fmt.Sprintf("收到结果数据:%s\n", msg))

	err := xml.Unmarshal([]byte(msg), &ctx.cvi3_result)
	if err != nil {
		h.AudiVw.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
		return
	}

	//结果数据
	XML2Result(&ctx.cvi3_result, &ctx.controller_result)

	// 波形文件
	XML2Curve(&ctx.cvi3_result, &ctx.controller_curve_file)

	s_curvedata, _ := json.Marshal(ctx.controller_curve_file)
	ctx.controller_curve.CurveData = string(s_curvedata)
	ctx.controller_curve.Count = ctx.controller_result.Count
	ctx.controller_curve.CurveFile = ctx.controller_result.CurFile
	ctx.controller_curve.ResultID = ctx.controller_result.Result_id

	e := h.handleCurve(&ctx.controller_curve, ctx)
	if  e == nil {
		h.handleResult(&ctx.controller_result, ctx)
	} else {
		h.AudiVw.diag.Error("handleCurve err", err)
	}

}