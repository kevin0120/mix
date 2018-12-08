package controller

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/minio"
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

type SavePackage struct {
	controllerResult *ControllerResult
	dbResult         *storage.Results
	dbWorkorder      *storage.Workorders
}

type Handlers struct {
	controllerService *Service
	closing           chan struct{}
	workers           int
	wg                sync.WaitGroup

	saveBuffer chan *SavePackage
}

func (h *Handlers) Init(workers int) {

	h.workers = workers
	h.closing = make(chan struct{})
	h.saveBuffer = make(chan *SavePackage, 1024)

	go h.saveProcess()
}

func (h *Handlers) Release() {

	h.closing <- struct{}{}
	h.wg.Wait()
}

func (h *Handlers) Handle(result interface{}, curve interface{}) {
	controllerResult := result.(*ControllerResult)

	// 取得工单
	dbWorkorder, err := h.controllerService.DB.GetWorkorder(controllerResult.Workorder_ID, true)
	if err != nil {
		h.controllerService.diag.Error("get workorder failed", err)
		return
	}

	// 取得结果
	dbResult, err := h.controllerService.DB.FindTargetResultForJobManual(controllerResult.Workorder_ID)
	if err != nil {
		h.controllerService.diag.Error("get result failed", err)
		return
	}

	// 处理曲线
	if curve != nil {
		controllerCurve := curve.(*minio.ControllerCurve)
		curveFileName := fmt.Sprintf("%s_%s_%d_%d_%d.json",
			dbWorkorder.MO_Model, dbResult.NutNo, dbResult.Seq, dbResult.ResultId, dbResult.Count)

		controllerResult.CurFile = curveFileName

		controllerCurve.CurveFile = curveFileName
		controllerCurve.ResultID = dbResult.Id
		controllerCurve.Count = dbResult.Count
		controllerCurve.UpdateTime = controllerResult.Dat

		h.handleCurve(controllerCurve)
	}

	// 处理结果
	h.handleResult(controllerResult, &dbResult, &dbWorkorder)

}

// 处理结果数据
func (h *Handlers) handleResult(result *ControllerResult, dbResult *storage.Results, dbWorkorder *storage.Workorders) {

	// 保存结果
	pkg := SavePackage{
		controllerResult: result,
		dbResult:         dbResult,
		dbWorkorder:      dbWorkorder,
	}

	h.saveResult(&pkg)
}

// 处理曲线数据
func (h *Handlers) handleCurve(curve *minio.ControllerCurve) {
	// 保存对象存储
	h.controllerService.Minio.Save(curve)
}

// 异步保存
func (h *Handlers) saveProcess() {
	for {
		select {
		case data := <-h.saveBuffer:
			h.handleSaveResult(data)

		case <-h.closing:
			h.wg.Done()
			return
		}
	}
}

// 保存结果
func (h *Handlers) saveResult(data *SavePackage) {
	h.saveBuffer <- data

	// 推送hmi
	if data.controllerResult.NeedPushHmi {
		wsResult := wsnotify.WSResult{
			Result:   data.controllerResult.Result,
			MI:       data.controllerResult.ResultValue.Mi,
			WI:       data.controllerResult.ResultValue.Wi,
			TI:       data.controllerResult.ResultValue.Ti,
			GroupSeq: data.dbResult.GroupSeq,
			Batch:    data.dbResult.Batch,
		}

		wsResults := []wsnotify.WSResult{}
		wsResults = append(wsResults, wsResult)
		ws_str, _ := json.Marshal(wsResults)
		h.controllerService.WS.WSSendResult(data.dbWorkorder.HMISN, string(ws_str))
	}

}

// 处理保存结果
func (h *Handlers) handleSaveResult(data *SavePackage) {

	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation("2006-01-02 15:04:05", data.controllerResult.Dat, loc)
	data.dbResult.UpdateTime = dt.UTC()

	data.dbResult.Result = data.controllerResult.Result
	data.dbResult.ControllerSN = data.controllerResult.Controller_SN
	s_value, _ := json.Marshal(data.controllerResult.ResultValue)
	s_pset, _ := json.Marshal(data.controllerResult.PSetDefine)

	data.dbResult.ResultValue = string(s_value)
	data.dbResult.PSetDefine = string(s_pset)
	data.dbResult.TighteningID = data.controllerResult.TighteningID

	if data.dbResult.Batch == "" {
		data.dbResult.Batch = fmt.Sprintf("%d/%d", data.dbResult.Seq, data.dbWorkorder.MaxSeq)
	}

	if data.dbResult.GunSN == "" {
		data.dbResult.GunSN = data.controllerResult.GunSN
	}

	if data.dbResult.Count >= int(data.dbResult.MaxRedoTimes) || data.dbResult.Result == storage.RESULT_OK {
		data.dbResult.Stage = storage.RESULT_STAGE_FINAL

		if data.dbResult.Batch != "" {
			kvs := strings.Split(data.dbResult.Batch, "/")
			if kvs[0] == kvs[1] {
				h.controllerService.diag.Debug("工单已完成")
				data.dbWorkorder.Status = "done"
				h.controllerService.DB.UpdateWorkorder(data.dbWorkorder)
			}
		}
	} else {

		data.dbResult.Count += 1
	}

	// 保存结果
	_, err := h.controllerService.DB.UpdateResult(data.dbResult)
	if err != nil {
		h.controllerService.diag.Error("缓存结果失败", err)
	} else {
		h.controllerService.diag.Debug("缓存结果成功")
	}

	// 推送aiis
	if data.controllerResult.NeedPushAiis || data.dbResult.Stage == storage.RESULT_STAGE_FINAL {

		aiisResult, err := h.controllerService.Aiis.ResultToAiisResult(data.dbResult)

		if err == nil {
			h.controllerService.Aiis.PutResult(data.dbResult.ResultId, aiisResult)
		}
	}
}
