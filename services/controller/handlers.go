package controller

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
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
	dbWorkorder      *storage.Workorders
	dbResult         *storage.Results
	consume          *Consume
	count            int
	batch            string
	curveFile        string
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

	// 取得工具信息
	//gun, err := h.controllerService.DB.GetGun(controllerResult.GunSN)
	//if err != nil {
	//	h.controllerService.diag.Error("get gun failed", err)
	//	return
	//}

	// 取得工单
	dbWorkorder, err := h.controllerService.DB.GetWorkorder(controllerResult.Workorder_ID, true)
	if err != nil {
		h.controllerService.diag.Error("get workorder failed", err)
		return
	}

	consumes := []Consume{}
	json.Unmarshal([]byte(dbWorkorder.Consumes), &consumes)

	targetConsume := consumes[0]
	for _, v := range consumes {
		if v.GroupSeq == controllerResult.Seq {
			targetConsume = v
		}
	}

	// 处理曲线
	curveFileName := ""
	if curve != nil {
		controllerCurve := curve.(*minio.ControllerCurve)
		curveFileName = fmt.Sprintf("%s_%s_%d_%d_%d.json",
			dbWorkorder.MO_Model, targetConsume.NutNo, dbWorkorder.WorkorderID, controllerResult.Seq, controllerResult.Count)

		controllerResult.CurFile = aiis.CURObject{
			File: curveFileName,
			OP:   controllerResult.Count,
		}

		controllerCurve.CurveFile = curveFileName
		//controllerCurve.ResultID = dbResult.Id
		controllerCurve.Count = controllerResult.Count
		controllerCurve.UpdateTime = controllerResult.Dat

		h.handleCurve(controllerCurve)
	}

	pkg := SavePackage{
		controllerResult: controllerResult,
		dbWorkorder:      &dbWorkorder,
		consume:          &targetConsume,
		count:            controllerResult.Count,
		batch:            fmt.Sprintf("%d/%d", targetConsume.GroupSeq, len(consumes)),
		curveFile:        curveFileName,
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
			GroupSeq: data.consume.GroupSeq,
			Batch:    data.batch,
		}

		wsResults := []wsnotify.WSResult{}
		wsResults = append(wsResults, wsResult)
		ws_str, _ := json.Marshal(wsResults)
		h.controllerService.WS.WSSendResult(data.dbWorkorder.HMISN, string(ws_str))
	}

}

func magicTrick(t time.Time) time.Time {
	return t.Add(-8 * time.Hour)
}

// 处理保存结果
func (h *Handlers) handleSaveResult(data *SavePackage) {

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

	dbResult.UpdateTime = magicTrick(dt.UTC())

	dbResult.Result = data.controllerResult.Result
	dbResult.ControllerSN = data.controllerResult.Controller_SN
	s_value, _ := json.Marshal(data.controllerResult.ResultValue)
	s_pset, _ := json.Marshal(data.controllerResult.PSetDefine)

	dbResult.ResultValue = string(s_value)
	dbResult.PSetDefine = string(s_pset)
	dbResult.TighteningID = data.controllerResult.TighteningID
	dbResult.Batch = data.batch
	dbResult.GunSN = data.controllerResult.GunSN
	dbResult.ExInfo = data.curveFile

	dbResult.PSet, _ = strconv.Atoi(data.consume.PSet)
	dbResult.ToleranceMax = data.consume.ToleranceMax
	dbResult.ToleranceMin = data.consume.ToleranceMin
	dbResult.ToleranceMaxDegree = data.consume.ToleranceMaxDegree
	dbResult.ToleranceMinDegree = data.consume.ToleranceMinDegree
	dbResult.NutNo = data.consume.NutNo

	dbResult.HasUpload = false
	dbResult.Count = data.count
	dbResult.UserID = 1

	dbResult.OffsetX = data.consume.X
	dbResult.OffsetY = data.consume.Y

	dbResult.GroupSeq = data.consume.GroupSeq
	dbResult.Seq = data.consume.Seq
	dbResult.MaxRedoTimes = data.consume.Max_redo_times
	dbResult.WorkorderID = data.dbWorkorder.Id

	if data.count >= int(data.consume.Max_redo_times) || dbResult.Result == storage.RESULT_OK {
		dbResult.Stage = storage.RESULT_STAGE_FINAL

		kvs := strings.Split(dbResult.Batch, "/")
		if kvs[0] == kvs[1] {
			h.controllerService.diag.Debug("工单已完成")
			data.dbWorkorder.Status = "done"
			h.controllerService.DB.UpdateWorkorder(data.dbWorkorder)
		}
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
