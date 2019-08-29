package controller

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
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
	dbWorkorder      *storage.Workorders
	dbResult         *storage.Results
	//consume          *Consume
	count     int
	batch     string
	curveFile string
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

func (h *Handlers) Handle(result ControllerResult, curve minio.ControllerCurve) {

	// 取得工具信息
	//gun, err := h.controllerService.StorageService.GetGun(controllerResult.GunSN)
	//if err != nil {
	//	h.controllerService.diag.Error("get gun failed", err)
	//	return
	//}

	// 取得工单
	//dbWorkorder, err := h.controllerService.StorageService.GetWorkorder(controllerResult.Workorder_ID, true)
	//if err != nil {
	//	h.controllerService.diag.Error("get workorder failed", err)
	//	return
	//}

	// 取得工步
	//dbStep, err := h.controllerService.StorageService.GetStep(controllerResult.StepID)
	//if err != nil {
	//	h.controllerService.diag.Error("get step failed", err)
	//	return
	//}

	//consumes := []Consume{}
	//json.Unmarshal([]byte(dbWorkorder.Consumes), &consumes)
	//
	//targetConsume := consumes[0]
	//maxGroupSeq := 0
	//for _, v := range consumes {
	//	if v.GroupSeq == controllerResult.Seq {
	//		targetConsume = v
	//	}
	//
	//	if v.GroupSeq >= maxGroupSeq {
	//		maxGroupSeq = v.GroupSeq
	//	}
	//}

	// 处理曲线
	curveFile := ""
	result.CurFile = aiis.CURObject{
		File: curve.CurveFile,
		OP:   curve.Count,
	}

	curveFile = curve.CurveFile

	//controllerCurve.ResultID = dbResult.Id
	curve.Count = result.Count
	curve.UpdateTime = result.Dat

	h.handleCurve(&curve)

	pkg := SavePackage{
		controllerResult: &result,
		//dbWorkorder:      &dbWorkorder,
		//consume:          &targetConsume,
		count: result.Count,
		//batch:            fmt.Sprintf("%d/%d", targetConsume.GroupSeq, maxGroupSeq),
		curveFile: curveFile,
	}

	dbResult := h.doSaveResult(&pkg)

	pkg.dbResult = &dbResult

	// 处理结果
	h.saveResult(&pkg)

}

// 处理曲线数据
func (h *Handlers) handleCurve(curve *minio.ControllerCurve) {
	// 保存对象存储
	h.controllerService.Minio.Save(curve)
}

//func (h *Handlers) HandleCurve(curve *minio.ControllerCurve) {
//	// 保存对象存储
//	h.handleCurve(curve)
//}

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
			GroupSeq: data.controllerResult.GroupSeq,
			Seq:      data.controllerResult.Seq,
			Batch:    data.controllerResult.Batch,
			ToolSN:   data.controllerResult.GunSN,
		}

		wsResults := []wsnotify.WSResult{}
		wsResults = append(wsResults, wsResult)

		msg := wsnotify.WSMsg{
			Type: WS_TIGHTENING_RESULT,
			Data: wsResults,
		}

		str, _ := json.Marshal(msg)
		h.controllerService.WS.WSSend(wsnotify.WS_EVENT_RESULT, string(str))

	}

}

// 处理保存结果
func (h *Handlers) handleSaveResult(data *SavePackage) {

	h.controllerService.diag.Debug(fmt.Sprintf("receive result: %s\n", data.controllerResult.Raw))

	// 推送aiis
	aiisResult, err := h.controllerService.Aiis.ResultToAiisResult(data.dbResult)

	if err == nil {
		h.controllerService.Aiis.PutResult(data.dbResult.ResultId, aiisResult)
	}
}

func (h *Handlers) doSaveResult(data *SavePackage) storage.Results {
	dbResult := storage.Results{}

	dt := time.Now()
	if data.controllerResult.Dat != "" {
		loc, _ := time.LoadLocation("Local")
		dt, _ = time.ParseInLocation("2006-01-02 15:04:05", data.controllerResult.Dat, loc)
	}

	dbResult.UpdateTime = dt.UTC()

	dbResult.Result = data.controllerResult.Result
	dbResult.ControllerSN = data.controllerResult.Controller_SN
	s_value, _ := json.Marshal(data.controllerResult.ResultValue)
	s_pset, _ := json.Marshal(data.controllerResult.PSetDefine)

	dbResult.ResultValue = string(s_value)
	dbResult.PSetDefine = string(s_pset)
	dbResult.TighteningID = data.controllerResult.TighteningID
	dbResult.Batch = data.controllerResult.Batch
	dbResult.GunSN = data.controllerResult.GunSN
	dbResult.ExInfo = data.curveFile

	dbResult.PSet = data.controllerResult.PSet
	dbResult.ToleranceMax = data.controllerResult.PSetDefine.Mp
	dbResult.ToleranceMin = data.controllerResult.PSetDefine.Mm
	dbResult.ToleranceMaxDegree = data.controllerResult.PSetDefine.Wp
	dbResult.ToleranceMinDegree = data.controllerResult.PSetDefine.Wm
	//dbResult.NutNo = data.consume.NutNo

	dbResult.HasUpload = false
	dbResult.Count = data.count
	dbResult.UserID = 1

	//dbResult.OffsetX = data.consume.X
	//dbResult.OffsetY = data.consume.Y

	//dbResult.GroupSeq = data.consume.GroupSeq
	//dbResult.Seq = data.consume.Seq
	//dbResult.MaxRedoTimes = data.consume.Max_redo_times
	dbResult.WorkorderID = data.controllerResult.Workorder_ID
	//dbResult.StepID = data.

	if data.count >= int(data.controllerResult.MaxRedoTime) || dbResult.Result == storage.RESULT_OK {
		dbResult.Stage = storage.RESULT_STAGE_FINAL

		//kvs := strings.Split(dbResult.Batch, "/")
		//if kvs[0] == kvs[1] {
		//	h.controllerService.diag.Debug("工单已完成")
		//	data.dbWorkorder.Status = "done"
		//	h.controllerService.StorageService.UpdateWorkorder(data.dbWorkorder)
		//}
	} else {

		dbResult.Stage = storage.RESULT_STAGE_INIT
	}

	// 保存结果
	err := h.controllerService.DB.CreateResult(&dbResult)
	if err != nil {
		h.controllerService.diag.Error("缓存结果失败", err)
	} else {
		h.controllerService.diag.Debug("缓存结果成功")
	}

	return dbResult
}
