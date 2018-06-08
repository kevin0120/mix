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

type Handlers struct {
	AudiVw	*Service
}

// 处理结果数据
func (h *Handlers) handleResult(result ControllerResult) (error) {
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

		ids := []int64{}
		json.Unmarshal([]byte(workorder.ResultIDs), &ids)
		if r.ResultId == ids[len(ids) - 1] {
			// 标记工单已完成
			workorder.Status = "finished"
			h.AudiVw.DB.UpdateWorkorder(workorder)

		}
	}

	_, err = h.AudiVw.DB.UpdateResult(r)
	if err != nil {
		fmt.Printf("缓存结果失败:%s\n", err.Error())
		return nil
	} else {
		fmt.Printf("缓存结果成功\n")
	}

	// 结果推送hmi
	ws_result := wsnotify.WSResult{}
	ws_result.Result_id = result.Result_id
	ws_result.Count = result.Count
	ws_result.Result = result.Result
	ws_result.MI = result.ResultValue.Mi
	ws_result.WI = result.ResultValue.Wi
	ws_result.TI = result.ResultValue.Ti
	ws_str, _ := json.Marshal(ws_result)

	fmt.Printf("Websocket推送结果到HMI\n")
	go h.AudiVw.WS.WSSendResult(workorder.HMISN, string(ws_str))

	if need_push_aiis {

		// 结果推送AIIS
		odoo_result := aiis.AIISResult{}
		if r.Result == RESULT_OK {
			odoo_result.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				odoo_result.One_time_pass = ODOO_RESULT_PASS
			} else {
				odoo_result.One_time_pass = ODOO_RESULT_FAIL
			}
		} else {
			odoo_result.Final_pass = ODOO_RESULT_FAIL
			odoo_result.One_time_pass = ODOO_RESULT_FAIL
		}

		odoo_result.Control_date = r.UpdateTime.Format(time.RFC3339)

		odoo_result.Measure_degree = result.ResultValue.Wi
		odoo_result.Measure_result = strings.ToLower(result.Result)
		odoo_result.Measure_t_don = result.ResultValue.Ti
		odoo_result.Measure_torque = result.ResultValue.Mi
		odoo_result.Op_time = result.Count
		odoo_result.Pset_m_max = result.PSetDefine.Mp
		odoo_result.Pset_m_min = result.PSetDefine.Mm
		odoo_result.Pset_m_target = result.PSetDefine.Ma
		odoo_result.Pset_m_threshold = result.PSetDefine.Ms
		odoo_result.Pset_strategy = result.PSetDefine.Strategy
		odoo_result.Pset_w_max = result.PSetDefine.Wp
		odoo_result.Pset_w_min = result.PSetDefine.Wm
		odoo_result.Pset_w_target = result.PSetDefine.Wa

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}
		for _, v := range curves {
			curobject := aiis.CURObject{}
			curobject.OP = v.Count
			curobject.File = v.CurveFile
			odoo_result.CURObjects = append(odoo_result.CURObjects, curobject)
		}

		go h.PutResultToAIIS(&odoo_result, &r)

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
		_, err := h.AudiVw.DB.UpdateResult(*db_result)
		if err != nil {

		}
	} else {
		fmt.Printf("推送AIIS失败\n")

	}
}

// 处理波形数据
func (h *Handlers) handleCurve(curve ControllerCurve) (error) {
	fmt.Printf("处理波形数据 ...\n")

	// 保存波形到数据库
	c := storage.Curves{}
	c.ResultID = curve.ResultID
	c.CurveData = curve.CurveData
	c.CurveFile = curve.CurveFile
	c.Count = curve.Count
	c.HasUpload = false

	exist, err := h.AudiVw.DB.CurveExist(c)
	if err != nil {
		return err
	} else {
		fmt.Printf("缓存波形数据到数据库 ...\n")
		if exist {
			_, err := h.AudiVw.DB.UpdateCurve(c)
			if err != nil {
				fmt.Printf("缓存波形失败\n")
				return err
			}
		} else {
			err := h.AudiVw.DB.Store(c)
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
		c.HasUpload = true
		fmt.Printf("对象存储保存成功，更新本地结果标识\n")
		_, err = h.AudiVw.DB.UpdateCurve(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string) {

	fmt.Printf("收到结果数据:%s\n", msg)

	result := CVI3Result{}
	err := xml.Unmarshal([]byte(msg), &result)
	if err != nil {
		fmt.Printf("HandleMsg err(struct):%s\n", err.Error())
	}

	// 结果数据
	result_data := XML2Result(result)

	// 波形文件
	curve_file := XML2Curve(result)

	curve := ControllerCurve{}
	s_curvedata, _ := json.Marshal(curve_file)
	curve.CurveData = string(s_curvedata)
	curve.Count = result_data.Count
	curve.CurveFile = result_data.CurFile
	curve.ResultID = result_data.Result_id

	e := h.handleCurve(curve)
	if  e == nil {
		h.handleResult(result_data)
	} else {
		fmt.Printf("HandleMsg err:%s\n", e.Error())
	}

}