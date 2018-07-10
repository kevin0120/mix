package audi_vw

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"strings"
	"sync"
	"time"
)

const (
	ODOO_RESULT_PASS = "pass"
	ODOO_RESULT_FAIL = "fail"

	QUALITY_STATE_PASS = "pass"
	QUALITY_STATE_FAIL = "fail"
	QUALITY_STATE_EX   = "exception"

	HANDLER_TYPE_CURVE = "curve"
	HANDLER_TYPE_AIIS  = "aiis"
)

type HandlerPkgCurve struct {
	curve *ControllerCurve
	dat   string
}

type HandlerPkgAiis struct {
	needPush  bool
	r         *storage.Results
	workorder *storage.Workorders
	result    *ControllerResult
}

type HandlerPkg struct {
	HandlerType string
	Pkg         interface{}
}

type Handlers struct {
	AudiVw     *Service
	HandlerBuf chan HandlerPkg
	closing    chan struct{}
	workers    int
	wg         sync.WaitGroup
}

func (h *Handlers) Init(workers int) {

	h.workers = workers
	h.closing = make(chan struct{})
	h.HandlerBuf = make(chan HandlerPkg, workers)

	for i := 0; i < h.workers; i++ {
		h.wg.Add(1)
		go h.asyncHandlerProcess()
	}
}

func (h *Handlers) Release() {
	for i := 0; i < h.workers; i++ {
		h.closing <- struct{}{}
	}

	h.wg.Wait()
}

func (h *Handlers) asyncHandlerProcess() {
	for {
		select {
		case pkg := <-h.HandlerBuf:
			switch pkg.HandlerType {
			case HANDLER_TYPE_AIIS:
				pkg_aiis := pkg.Pkg.(HandlerPkgAiis)
				h.PushAiis(pkg_aiis.needPush, pkg_aiis.r, pkg_aiis.workorder, pkg_aiis.result)

			case HANDLER_TYPE_CURVE:
				pkg_curve := pkg.Pkg.(HandlerPkgCurve)
				h.handleCurve(pkg_curve.curve, pkg_curve.dat)
			}

		case <-h.closing:
			h.wg.Done()
			return
		}
	}
}

func (h *Handlers) PushAiis(needPush bool, r *storage.Results, workorder *storage.Workorders, result *ControllerResult) error {

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
		aiisResult := aiis.AIISResult{}
		if r.Result == RESULT_OK {
			aiisResult.Final_pass = ODOO_RESULT_PASS
			if r.Count == 1 {
				aiisResult.One_time_pass = ODOO_RESULT_PASS
			} else {
				aiisResult.One_time_pass = ODOO_RESULT_FAIL
			}

			if (result.ResultValue.Mi >= r.ToleranceMin && result.ResultValue.Mi <= r.ToleranceMax) &&
				(result.ResultValue.Wi >= r.ToleranceMinDegree && result.ResultValue.Wi <= r.ToleranceMaxDegree) {
				aiisResult.QualityState = QUALITY_STATE_PASS
				aiisResult.ExceptionReason = ""
			} else {
				aiisResult.QualityState = QUALITY_STATE_EX
				aiisResult.ExceptionReason = QUALITY_STATE_EX
			}

		} else {
			aiisResult.Final_pass = ODOO_RESULT_FAIL
			aiisResult.One_time_pass = ODOO_RESULT_FAIL

			if (result.ResultValue.Mi >= r.ToleranceMin && result.ResultValue.Mi <= r.ToleranceMax) &&
				(result.ResultValue.Wi >= r.ToleranceMinDegree && result.ResultValue.Wi <= r.ToleranceMaxDegree) {

				aiisResult.QualityState = QUALITY_STATE_EX
				aiisResult.ExceptionReason = QUALITY_STATE_EX
			} else {
				aiisResult.QualityState = QUALITY_STATE_FAIL
				aiisResult.ExceptionReason = ""
			}

		}

		aiisResult.Control_date = r.UpdateTime.Format(time.RFC3339)

		aiisResult.Measure_degree = result.ResultValue.Wi
		aiisResult.Measure_result = strings.ToLower(result.Result)
		aiisResult.Measure_t_don = result.ResultValue.Ti
		aiisResult.Measure_torque = result.ResultValue.Mi
		aiisResult.Op_time = result.Count
		aiisResult.Pset_m_max = result.PSetDefine.Mp
		aiisResult.Pset_m_min = result.PSetDefine.Mm
		aiisResult.Pset_m_target = result.PSetDefine.Ma
		aiisResult.Pset_m_threshold = result.PSetDefine.Ms
		aiisResult.Pset_strategy = result.PSetDefine.Strategy
		aiisResult.Pset_w_max = result.PSetDefine.Wp
		aiisResult.Pset_w_min = result.PSetDefine.Wm
		aiisResult.Pset_w_target = result.PSetDefine.Wa
		aiisResult.Pset_w_threshold = 1
		aiisResult.UserID = result.UserID
		aiisResult.Seq = r.Seq

		// mo相关
		aiisResult.MO_AssemblyLine = workorder.MO_AssemblyLine
		aiisResult.MO_EquipemntName = workorder.MO_EquipemntName
		aiisResult.MO_FactoryName = workorder.MO_FactoryName
		aiisResult.MO_Pin = workorder.MO_Pin
		aiisResult.MO_Pin_check_code = workorder.MO_Pin_check_code
		aiisResult.MO_Year = workorder.MO_Year
		aiisResult.MO_Lnr = workorder.MO_Lnr
		aiisResult.MO_NutNo = r.NutNo
		aiisResult.MO_Model = workorder.MO_Model

		curves, err := h.AudiVw.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		aiisCurve := aiis.CURObject{}
		for _, v := range curves {
			aiisCurve.OP = v.Count
			aiisCurve.File = v.CurveFile
			aiisResult.CURObjects = append(aiisResult.CURObjects, aiisCurve)
		}

		h.AudiVw.diag.Debug("推送结果数据到AIIS ...")

		err = h.AudiVw.Aiis.PutResult(r.ResultId, aiisResult)
		if err == nil {
			h.AudiVw.diag.Debug("推送AIIS成功，更新本地结果标识")
		} else {
			h.AudiVw.diag.Error("推送AIIS失败", err)
		}
	}

	h.AudiVw.diag.Debug(fmt.Sprintf("PushAiis end:%s", time.Now().Format("2006-01-02 15:04:05.999999999")))
	return nil
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, dbresult *storage.Results, dbworkorder *storage.Workorders) error {
	h.AudiVw.diag.Debug("处理结果数据 ...")

	needPushAiis := false

	loc, _ := time.LoadLocation("Local")
	dbresult.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", result.Dat, loc)
	dbresult.Result = result.Result
	dbresult.Count = result.Count
	dbresult.HasUpload = false
	dbresult.ControllerSN = result.Controller_SN
	dbresult.UserID = result.UserID
	s_value, _ := json.Marshal(result.ResultValue)
	s_pset, _ := json.Marshal(result.PSetDefine)

	dbresult.ResultValue = string(s_value)
	dbresult.PSetDefine = string(s_pset)

	if dbresult.Count >= int(dbresult.MaxRedoTimes) || dbresult.Result == RESULT_OK {
		needPushAiis = true
		dbresult.Stage = RESULT_STAGE_FINAL

		if dbresult.ResultId == dbworkorder.LastResultID {
			h.AudiVw.diag.Debug("工单已完成")
			dbworkorder.Status = "done"
			h.AudiVw.DB.UpdateWorkorder(dbworkorder)
		}
	}

	// 结果推送hmi
	wsResult := wsnotify.WSResult{}
	wsResult.Result_id = result.Result_id
	wsResult.Count = result.Count
	wsResult.Result = result.Result
	wsResult.MI = result.ResultValue.Mi
	wsResult.WI = result.ResultValue.Wi
	wsResult.TI = result.ResultValue.Ti
	ws_str, _ := json.Marshal(wsResult)

	h.AudiVw.diag.Debug("Websocket推送结果到HMI")

	h.AudiVw.WS.WSSendResult(dbworkorder.HMISN, string(ws_str))

	pkg_aiis := HandlerPkgAiis{
		needPush:  needPushAiis,
		r:         dbresult,
		workorder: dbworkorder,
		result:    result,
	}

	pkg := HandlerPkg{
		HandlerType: HANDLER_TYPE_AIIS,
		Pkg:         pkg_aiis,
	}
	h.HandlerBuf <- pkg

	return nil
}

// 处理波形数据
func (h *Handlers) handleCurve(curve *ControllerCurve, dat string) error {
	h.AudiVw.diag.Debug("处理波形数据 ...")

	dbCurve := storage.Curves{}
	dbCurve.ResultID = curve.ResultID
	dbCurve.CurveData = curve.CurveData
	dbCurve.CurveFile = curve.CurveFile
	dbCurve.Count = curve.Count
	dbCurve.HasUpload = false
	loc, _ := time.LoadLocation("Local")
	dbCurve.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", dat, loc)

	// 保存波形到对象存储
	h.AudiVw.diag.Debug("保存波形数据到对象存储 ...")
	err := h.AudiVw.Minio.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		h.AudiVw.diag.Error("对象存储保存失败", err)
	} else {
		dbCurve.HasUpload = true
		h.AudiVw.diag.Debug("对象存储保存成功")
	}

	// 保存波形到数据库
	//exist, err := h.AudiVw.DB.CurveExist(&dbCurve)
	raw_curve, err := h.AudiVw.DB.GetCurve(&dbCurve)
	if err != nil {
		return err
	} else {
		h.AudiVw.diag.Debug("缓存波形数据到数据库 ...")
		if raw_curve != nil {
			dbCurve.Id = raw_curve.(storage.Curves).Id
			_, err := h.AudiVw.DB.UpdateCurve(&dbCurve)
			if err != nil {
				h.AudiVw.diag.Error("缓存波形失败", err)
				return err
			}
		} else {
			err := h.AudiVw.DB.Store(dbCurve)
			if err != nil {
				h.AudiVw.diag.Error("缓存波形失败", err)
				return err
			}
		}

		h.AudiVw.diag.Debug("缓存波形成功")
	}

	h.AudiVw.diag.Debug(fmt.Sprintf("handleCurve end:%s", time.Now().Format("2006-01-02 15:04:05.999999999")))
	return nil
}

// 处理收到的数据
func (h *Handlers) HandleMsg(msg string) {

	h.AudiVw.diag.Debug(fmt.Sprintf("HandleMsg begin:%s", time.Now().Format("2006-01-02 15:04:05.999999999")))
	//h.AudiVw.diag.Debug(fmt.Sprintf("收到结果数据:%s\n", msg))

	cvi3Result := CVI3Result{}
	err := xml.Unmarshal([]byte(msg), &cvi3Result)
	if err != nil {
		h.AudiVw.diag.Error(fmt.Sprint("HandlerMsg err:", msg), err)
		return
	}

	//结果数据
	controllerResult := ControllerResult{}
	XML2Result(&cvi3Result, &controllerResult)

	// 波形文件
	controllerCurveFile := ControllerCurveFile{}
	XML2Curve(&cvi3Result, &controllerCurveFile)

	result := storage.Results{}
	result, err = h.AudiVw.DB.GetResult(controllerResult.Result_id, 0)
	if err != nil {
		h.AudiVw.diag.Error("Cannot find result", err)
		return
	}

	workorder := storage.Workorders{}
	workorder, err = h.AudiVw.DB.GetWorkorder(controllerResult.Workorder_ID)
	if err != nil {
		h.AudiVw.diag.Error("Cannot find workorder", err)
		return
	}

	controllerResult.CurFile = fmt.Sprintf("%s_%s_%d_%d_%d.json",
		workorder.MO_Model, result.NutNo, result.Seq, result.ResultId, controllerResult.Count)

	sCurvedata, _ := json.Marshal(controllerCurveFile)
	controllerCurve := ControllerCurve{}
	controllerCurve.CurveData = string(sCurvedata)
	controllerCurve.Count = controllerResult.Count
	controllerCurve.CurveFile = controllerResult.CurFile
	controllerCurve.ResultID = controllerResult.Result_id

	pkg_curve := HandlerPkgCurve{
		curve: &controllerCurve,
		dat:   controllerResult.Dat,
	}

	pkg := HandlerPkg{
		HandlerType: HANDLER_TYPE_CURVE,
		Pkg:         pkg_curve,
	}
	h.HandlerBuf <- pkg

	h.handleResult(&controllerResult, &result, &workorder)

	h.AudiVw.diag.Debug(fmt.Sprintf("HandleMsg end:%s", time.Now().Format("2006-01-02 15:04:05.999999999")))
}
