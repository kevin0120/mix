package controller

import (
	"encoding/json"
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
	curve_obj *CurveObject
}

type HandlerPkg struct {
	HandlerType string
	Pkg         interface{}
}

type Handlers struct {
	controllerService *Service
	HandlerBuf        chan HandlerPkg
	closing           chan struct{}
	workers           int
	wg                sync.WaitGroup
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
				h.PushAiis(pkg_aiis.needPush, pkg_aiis.r, pkg_aiis.workorder, pkg_aiis.result, pkg_aiis.curve_obj)

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

func (h *Handlers) PushAiis(needPush bool, r *storage.Results, workorder *storage.Workorders, result *ControllerResult, curve_obj *CurveObject) error {

	if needPush {

		// 结果推送AIIS
		aiisResult := aiis.AIISResult{}
		aiisResult.ID = r.Id
		if r.Result == storage.RESULT_OK {
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

		if result.ExceptionReason != "" {
			aiisResult.ExceptionReason = result.ExceptionReason + aiisResult.ExceptionReason
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

		curves, err := h.controllerService.DB.ListCurvesByResult(result.Result_id)
		if err != nil {
			return err
		}

		curve_exist := false
		aiisCurve := aiis.CURObject{}
		for _, v := range curves {
			if v.CurveFile == curve_obj.File {
				curve_exist = true
			}
			aiisCurve.OP = v.Count
			aiisCurve.File = v.CurveFile
			aiisResult.CURObjects = append(aiisResult.CURObjects, aiisCurve)
		}

		if !curve_exist {
			aiisCurve.OP = curve_obj.Count
			aiisCurve.File = curve_obj.File
			aiisResult.CURObjects = append(aiisResult.CURObjects, aiisCurve)
		}

		h.controllerService.diag.Debug("推送结果数据到AIIS ...")

		err = h.controllerService.Aiis.PutResult(r.ResultId, aiisResult)
		if err == nil {
			h.controllerService.diag.Debug("推送AIIS成功，更新本地结果标识")
		} else {
			h.controllerService.diag.Error("推送AIIS失败", err)
		}
	}

	return nil
}

//
func (h *Handlers) SaveResult(result *ControllerResult, dbresult *storage.Results, create bool) (*storage.Results, error) {
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

	h.controllerService.diag.Debug("缓存结果到数据库 ...")

	var err error = nil
	if create {
		err = h.controllerService.DB.Store(dbresult)
	} else {
		_, err = h.controllerService.DB.UpdateResult(dbresult)
	}

	if err != nil {
		h.controllerService.diag.Error("缓存结果失败", err)
	} else {
		h.controllerService.diag.Debug("缓存结果成功")
	}

	return dbresult, nil
}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, dbresult *storage.Results, dbworkorder *storage.Workorders, curve_obj *CurveObject) error {
	h.controllerService.diag.Debug("处理结果数据 ...")

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

	if dbresult.Count >= int(dbresult.MaxRedoTimes) || dbresult.Result == storage.RESULT_OK {
		needPushAiis = true
		dbresult.Stage = storage.RESULT_STAGE_FINAL

		if dbresult.ResultId == dbworkorder.LastResultID {
			h.controllerService.diag.Debug("工单已完成")
			dbworkorder.Status = "done"
			h.controllerService.DB.UpdateWorkorder(dbworkorder)
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
	//wsResult.Seq = dbresult.Seq
	ws_str, _ := json.Marshal(wsResult)

	h.controllerService.diag.Debug("Websocket推送结果到HMI")

	h.controllerService.WS.WSSendResult(dbworkorder.HMISN, string(ws_str))

	if needPushAiis {
		pkg_aiis := HandlerPkgAiis{
			needPush:  needPushAiis,
			r:         dbresult,
			workorder: dbworkorder,
			result:    result,
			curve_obj: curve_obj,
		}

		pkg := HandlerPkg{
			HandlerType: HANDLER_TYPE_AIIS,
			Pkg:         pkg_aiis,
		}

		h.HandlerBuf <- pkg
	}

	defer func() {
		h.controllerService.diag.Debug("缓存结果到数据库 ...")
		_, err := h.controllerService.DB.UpdateResult(dbresult)
		if err != nil {
			h.controllerService.diag.Error("缓存结果失败", err)
		} else {
			h.controllerService.diag.Debug("缓存结果成功")
		}
	}()

	return nil
}

// 处理波形数据
func (h *Handlers) handleCurve(curve *ControllerCurve, dat string) error {
	h.controllerService.diag.Debug("处理波形数据 ...")

	dbCurve := storage.Curves{}
	dbCurve.ResultID = curve.ResultID
	dbCurve.CurveData = curve.CurveData
	dbCurve.CurveFile = curve.CurveFile
	dbCurve.Count = curve.Count
	dbCurve.HasUpload = false
	loc, _ := time.LoadLocation("Local")
	dbCurve.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", dat, loc)

	// 保存波形到对象存储
	h.controllerService.diag.Debug("保存波形数据到对象存储 ...")
	err := h.controllerService.Minio.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		h.controllerService.diag.Error("对象存储保存失败", err)
	} else {
		dbCurve.HasUpload = true
		h.controllerService.diag.Debug("对象存储保存成功")
	}

	// 保存波形到数据库
	//exist, err := h.AudiVw.DB.CurveExist(&dbCurve)
	raw_curve, err := h.controllerService.DB.GetCurve(&dbCurve)
	if err != nil {
		return err
	} else {
		h.controllerService.diag.Debug("缓存波形数据到数据库 ...")
		if raw_curve != nil {
			dbCurve.Id = raw_curve.(storage.Curves).Id
			_, err := h.controllerService.DB.UpdateCurve(&dbCurve)
			if err != nil {
				h.controllerService.diag.Error("缓存波形失败", err)
				return err
			}
		} else {
			err := h.controllerService.DB.Store(dbCurve)
			if err != nil {
				h.controllerService.diag.Error("缓存波形失败", err)
				return err
			}
		}

		h.controllerService.diag.Debug("缓存波形成功")
	}

	return nil
}

// 处理收到的数据
func (h *Handlers) Handle(controllerResult interface{}, controllerCurveFile interface{}) {

	model_controllerResult := controllerResult.(ControllerResult)

	result := storage.Results{}
	result, err := h.controllerService.DB.GetResult(model_controllerResult.Result_id, 0)
	if err != nil {
		h.controllerService.diag.Error("Cannot find result", err)
		return
	}

	workorder := storage.Workorders{}
	workorder, err = h.controllerService.DB.GetWorkorder(model_controllerResult.Workorder_ID, true)
	if err != nil {
		h.controllerService.diag.Error("Cannot find workorder", err)
		return
	}

	// 车型-螺栓编号-结果sequence-结果id-拧接次数
	model_controllerResult.CurFile = fmt.Sprintf("%s_%s_%d_%d_%d.json",
		workorder.MO_Model, result.NutNo, result.Seq, result.ResultId, model_controllerResult.Count)

	if controllerCurveFile != nil {

		model_controllerCurveFile := controllerCurveFile.(ControllerCurveFile)
		sCurvedata, _ := json.Marshal(model_controllerCurveFile)
		controllerCurve := ControllerCurve{}
		controllerCurve.CurveData = string(sCurvedata)
		controllerCurve.Count = model_controllerResult.Count
		controllerCurve.CurveFile = model_controllerResult.CurFile
		controllerCurve.ResultID = model_controllerResult.Result_id

		pkg_curve := HandlerPkgCurve{
			curve: &controllerCurve,
			dat:   model_controllerResult.Dat,
		}

		pkg := HandlerPkg{
			HandlerType: HANDLER_TYPE_CURVE,
			Pkg:         pkg_curve,
		}
		h.HandlerBuf <- pkg
	}

	cur_obj := CurveObject{
		File:  model_controllerResult.CurFile,
		Count: model_controllerResult.Count,
	}

	h.handleResult(&model_controllerResult, &result, &workorder, &cur_obj)

}
