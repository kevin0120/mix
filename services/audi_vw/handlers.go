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
	HandlerContext HandlerContext
}

func (h *Handlers) Init() {
	h.HandlerContext = HandlerContext {
		cvi3_result: CVI3Result{},
		controller_curve: ControllerCurve{},
		controller_result: ControllerResult{},
		controller_curve_file: ControllerCurveFile{},
		db_curve: storage.Curves{},
		ws_result: wsnotify.WSResult{},
		aiis_result: aiis.AIISResult{},
		aiis_curve: aiis.CURObject{},
	}
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult) (error) {
	fmt.Printf("处理结果数据 ...\n")

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
	s_value, _ := json.Marshal(result.ResultValue)
	s_pset, _ := json.Marshal(result.PSetDefine)

	r.ResultValue = string(s_value)
	r.PSetDefine = string(s_pset)

	fmt.Printf("缓存结果到数据库 ...\n")

	if r.Count >= int(workorder.MaxRedoTimes) || r.Result == RESULT_OK {
		need_push_aiis = true
		r.Stage = RESULT_STAGE_FINAL

		json.Unmarshal([]byte(workorder.ResultIDs), &h.HandlerContext.result_ids)
		if r.ResultId == h.HandlerContext.result_ids[len(h.HandlerContext.result_ids) - 1] {
			// 标记工单已完成
			workorder.Status = "finished"
			h.AudiVw.DB.UpdateWorkorder(&workorder)

		}
	}

	_, err = h.AudiVw.DB.UpdateResult(&r)
	if err != nil {
		fmt.Printf("缓存结果失败:%s\n", err.Error())
		return nil
	} else {
		fmt.Printf("缓存结果成功\n")
	}

	// 结果推送hmi
	h.HandlerContext.ws_result.Result_id = result.Result_id
	h.HandlerContext.ws_result.Count = result.Count
	h.HandlerContext.ws_result.Result = result.Result
	h.HandlerContext.ws_result.MI = result.ResultValue.Mi
	h.HandlerContext.ws_result.WI = result.ResultValue.Wi
	h.HandlerContext.ws_result.TI = result.ResultValue.Ti
	ws_str, _ := json.Marshal(h.HandlerContext.ws_result)

	fmt.Printf("Websocket推送结果到HMI\n")
	h.AudiVw.WS.WSSendResult(workorder.HMISN, string(ws_str))

	if need_push_aiis {

		// 结果推送AIIS
		if r.Result == RESULT_OK {
			h.HandlerContext.aiis_result.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				h.HandlerContext.aiis_result.One_time_pass = ODOO_RESULT_PASS
			} else {
				h.HandlerContext.aiis_result.One_time_pass = ODOO_RESULT_FAIL
			}
		} else {
			h.HandlerContext.aiis_result.Final_pass = ODOO_RESULT_FAIL
			h.HandlerContext.aiis_result.One_time_pass = ODOO_RESULT_FAIL
		}

		h.HandlerContext.aiis_result.Control_date = r.UpdateTime.Format(time.RFC3339)

		h.HandlerContext.aiis_result.Measure_degree = result.ResultValue.Wi
		h.HandlerContext.aiis_result.Measure_result = strings.ToLower(result.Result)
		h.HandlerContext.aiis_result.Measure_t_don = result.ResultValue.Ti
		h.HandlerContext.aiis_result.Measure_torque = result.ResultValue.Mi
		h.HandlerContext.aiis_result.Op_time = result.Count
		h.HandlerContext.aiis_result.Pset_m_max = result.PSetDefine.Mp
		h.HandlerContext.aiis_result.Pset_m_min = result.PSetDefine.Mm
		h.HandlerContext.aiis_result.Pset_m_target = result.PSetDefine.Ma
		h.HandlerContext.aiis_result.Pset_m_threshold = result.PSetDefine.Ms
		h.HandlerContext.aiis_result.Pset_strategy = result.PSetDefine.Strategy
		h.HandlerContext.aiis_result.Pset_w_max = result.PSetDefine.Wp
		h.HandlerContext.aiis_result.Pset_w_min = result.PSetDefine.Wm
		h.HandlerContext.aiis_result.Pset_w_target = result.PSetDefine.Wa

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		h.HandlerContext.aiis_result.CURObjects = []aiis.CURObject{}
		for _, v := range curves {
			h.HandlerContext.aiis_curve.OP = v.Count
			h.HandlerContext.aiis_curve.File = v.CurveFile
			h.HandlerContext.aiis_result.CURObjects = append(h.HandlerContext.aiis_result.CURObjects, h.HandlerContext.aiis_curve)
		}

		h.PutResultToAIIS(&h.HandlerContext.aiis_result, &r)

	}

	return nil
}

func (h *Handlers) PutResultToAIIS(aiis_result *aiis.AIISResult, db_result *storage.Results) {
	fmt.Printf("推送结果数据到AIIS ...\n")

	_, err := h.AudiVw.Aiis.PutResult(db_result.ResultId, aiis_result)
	if err == nil {
		// 发送成功
		db_result.HasUpload = true
		fmt.Printf("推送AIIS成功，更新本地结果标识\n")
		_, err := h.AudiVw.DB.UpdateResult(db_result)
		if err != nil {

		}
	} else {
		fmt.Printf("推送AIIS失败\n")

	}
}

// 处理波形数据
func (h *Handlers) handleCurve(curve *ControllerCurve) (error) {
	fmt.Printf("处理波形数据 ...\n")

	// 保存波形到数据库
	h.HandlerContext.db_curve.ResultID = curve.ResultID
	h.HandlerContext.db_curve.CurveData = curve.CurveData
	h.HandlerContext.db_curve.CurveFile = curve.CurveFile
	h.HandlerContext.db_curve.Count = curve.Count
	h.HandlerContext.db_curve.HasUpload = false

	exist, err := h.AudiVw.DB.CurveExist(&h.HandlerContext.db_curve)
	if err != nil {
		return err
	} else {
		fmt.Printf("缓存波形数据到数据库 ...\n")
		if exist {
			_, err := h.AudiVw.DB.UpdateCurve(&h.HandlerContext.db_curve)
			if err != nil {
				fmt.Printf("缓存波形失败\n")
				return err
			}
		} else {
			err := h.AudiVw.DB.Store(h.HandlerContext.db_curve)
			if err != nil {
				fmt.Printf("缓存波形失败\n")
				return err
			}
		}

		fmt.Printf("缓存波形成功\n")
	}

	// 保存波形到对象存储
	fmt.Printf("保存波形数据到对象存储 ...\n")
	err = h.AudiVw.Minio.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		fmt.Printf("对象存储保存失败\n")
		return err
	} else {
		h.HandlerContext.db_curve.HasUpload = true
		fmt.Printf("对象存储保存成功，更新本地结果标识\n")
		_, err = h.AudiVw.DB.UpdateCurve(&h.HandlerContext.db_curve)
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string) {

	fmt.Printf("收到结果数据:%s\n", msg)

	//h.HandlerContext.controller_curve_file.CUR_M = []float64{}
	//h.HandlerContext.controller_curve_file.CUR_W = []float64{}
	//h.HandlerContext.result_ids = []int64{}
	//h.HandlerContext.aiis_result.CURObjects = []aiis.CURObject{}

	err := xml.Unmarshal([]byte(msg), &h.HandlerContext.cvi3_result)
	if err != nil {
		fmt.Printf("HandleMsg err(struct):%s\n", err.Error())
	}

	//结果数据
	XML2Result(&h.HandlerContext.cvi3_result, &h.HandlerContext.controller_result)

	// 波形文件
	XML2Curve(&h.HandlerContext.cvi3_result, &h.HandlerContext.controller_curve_file)

	s_curvedata, _ := json.Marshal(h.HandlerContext.controller_curve_file)
	h.HandlerContext.controller_curve.CurveData = string(s_curvedata)
	h.HandlerContext.controller_curve.Count = h.HandlerContext.controller_result.Count
	h.HandlerContext.controller_curve.CurveFile = h.HandlerContext.controller_result.CurFile
	h.HandlerContext.controller_curve.ResultID = h.HandlerContext.controller_result.Result_id

	e := h.handleCurve(&h.HandlerContext.controller_curve)
	if  e == nil {
		h.handleResult(&h.HandlerContext.controller_result)
	} else {
		fmt.Printf("HandleMsg err:%s\n", e.Error())
	}

}