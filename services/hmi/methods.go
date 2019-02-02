package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_USER_ID = 1
)

type Methods struct {
	service *Service
}

func (m *Methods) putToolControl(ctx iris.Context) {
	var err error
	var te ToolEnable
	err = ctx.ReadJSON(&te)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if te.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if te.GunSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("gun_sn is required")
		return
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[te.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.AUDIPROTOCOL:
		err = m.service.AudiVw.ToolControl(te.Controller_SN, te.GunSN, te.Enable)
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.ToolControl(te.Controller_SN, te.GunSN, te.Enable)

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

func (m *Methods) putPSets(ctx iris.Context) {

	var err error
	var pset PSet
	err = ctx.ReadJSON(&pset)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	str, _ := json.Marshal(pset)
	m.service.diag.Debug(fmt.Sprintf("new pset:%s", str))

	if pset.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if pset.GunSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("gun_sn is required")
		return
	}

	if pset.PSet == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset is required")
		return
	}

	if pset.Count == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count is required")
		return
	}

	if pset.WorkorderID == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workorder_id is required")
		return
	}

	// 检测count
	if pset.Count < 1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("tightening count should be greater than 0")
		return
	}

	// 检测工单
	workorder, err := m.service.DB.GetWorkorder(pset.WorkorderID, true)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	m.service.DB.UpdateGun(&storage.Guns{
		Serial:      pset.GunSN,
		WorkorderID: workorder.Id,
		Seq:         pset.GroupSeq,
		Count:       pset.Count,
	})

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[pset.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.AUDIPROTOCOL:
		err = m.service.AudiVw.PSet(pset.Controller_SN, pset.GunSN, pset.PSet, workorder.Id, int64(pset.GroupSeq), pset.Count, pset.UserID)
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.PSet(pset.Controller_SN, pset.GunSN, pset.PSet, pset.Result_id, pset.Count, pset.UserID)

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

func (m *Methods) putManualPSets(ctx iris.Context) {

	var err error
	var pset PSetManual
	err = ctx.ReadJSON(&pset)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if pset.UserID == 0 {
		pset.UserID = DEFAULT_USER_ID
	}

	if pset.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if pset.HmiSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi sn is required")
		return
	}

	//if pset.GunSN == "" {
	//	pset.GunSN = ""
	//}

	if pset.PSet == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset is required")
		return
	}

	if pset.Vin == "" {
		pset.Vin = "unknown"
	}

	if pset.CarType == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("CarType is required")
		return
	}

	if pset.Count < 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset count must be greater than or equal to 0")
		return
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[pset.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		err = m.insertResultsForPSet(&pset)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

		exInfo := fmt.Sprintf("%25s%25s%25s%25d", pset.Vin, pset.HmiSN, pset.CarType, pset.UserID)
		err = m.service.OpenProtocol.PSetManual(pset.Controller_SN, pset.GunSN, pset.PSet, pset.UserID, exInfo, pset.Count)

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

func (m *Methods) getPSetList(ctx iris.Context) {

	controller_sn := ctx.URLParam("controller_sn")

	if controller_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	c, exist := m.service.ControllerService.Controllers[controller_sn]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	var err error = nil
	var pset_list []int
	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		pset_list, err = m.service.OpenProtocol.GetPSetList(controller_sn)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	body, _ := json.Marshal(pset_list)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) getPSetDetail(ctx iris.Context) {

	controller_sn := ctx.URLParam("controller_sn")
	pset := ctx.URLParam("pset")

	if controller_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if pset == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset is required")
		return
	}

	v_pset, err := strconv.Atoi(pset)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset format error")
		return
	}

	c, exist := m.service.ControllerService.Controllers[controller_sn]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	var pset_detail openprotocol.PSetDetail
	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		pset_detail, err = m.service.OpenProtocol.GetPSetDetail(controller_sn, v_pset)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	body, _ := json.Marshal(pset_detail)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) getJobList(ctx iris.Context) {

	controller_sn := ctx.URLParam("controller_sn")

	if controller_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	c, exist := m.service.ControllerService.Controllers[controller_sn]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	var err error = nil
	var job_list []int
	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		job_list, err = m.service.OpenProtocol.GetJobList(controller_sn)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	body, _ := json.Marshal(job_list)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) getJobDetail(ctx iris.Context) {

	controller_sn := ctx.URLParam("controller_sn")
	job := ctx.URLParam("job")

	if controller_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if job == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("job is required")
		return
	}

	v_job, err := strconv.Atoi(job)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset format error")
		return
	}

	c, exist := m.service.ControllerService.Controllers[controller_sn]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	var job_detail openprotocol.JobDetail
	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		job_detail, err = m.service.OpenProtocol.GetJobDetail(controller_sn, v_job)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	body, _ := json.Marshal(job_detail)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) enableJobMode(ctx iris.Context) {
	var mode ControllerMode
	err := ctx.ReadJSON(&mode)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if mode.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if mode.Mode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("mode is required")
		return
	}

	c, exist := m.service.ControllerService.Controllers[mode.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		flag := true
		if mode.Mode == openprotocol.MODE_PSET {
			flag = false
		}
		err = m.service.OpenProtocol.JobOFF(mode.Controller_SN, flag)

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

func (m *Methods) putJobs(ctx iris.Context) {

	var err error
	var job Job
	err = ctx.ReadJSON(&job)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if job.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if job.Job == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("job is required")
		return
	}

	if job.WorkorderiD == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workorder_id is required")
		return
	}

	exist, _ := m.service.DB.WorkorderExists(job.WorkorderiD)
	if !exist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString("workorder not found")
		return
	}

	if job.UserID == 0 {
		job.UserID = DEFAULT_USER_ID
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[job.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.JobSet(job.Controller_SN, job.Job, job.WorkorderiD, job.UserID)

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	err = m.service.DB.InitWorkorderForJob(job.WorkorderiD)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

func (m *Methods) putManualJobs(ctx iris.Context) {

	var err error
	var job JobManual
	err = ctx.ReadJSON(&job)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if job.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if job.GunSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("gun_sn is required")
		return
	}

	if job.Job == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("job is required")
		return
	}

	//if job.CarType == "" {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("car type is required")
	//	return
	//}

	if job.HmiSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi sn is required")
		return
	}

	//if job.Vin == "" {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("vin is required")
	//	return
	//}

	if job.UserID == 0 {
		job.UserID = DEFAULT_USER_ID
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[job.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	var workorderID int64 = 0

	switch c.Protocol() {
	case controller.OPENPROTOCOL:

		var db_workorder *storage.Workorders
		ex_info := ""
		if !job.Skip {
			db_workorder, err = m.insertResultsForJob(&job)
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString(err.Error())
				return
			}

			workorderID = db_workorder.Id

			//vin-cartype-hmisn-userid
			ex_info = m.service.OpenProtocol.GenerateIDInfo(fmt.Sprintf("%d", db_workorder.Id))
		}

		if !job.HasSet {
			err = m.service.OpenProtocol.JobSetManual(job.Controller_SN, job.GunSN, job.Job, job.UserID, ex_info)
		} else {
			m.service.OpenProtocol.IDSet(job.Controller_SN, ex_info)
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	body, _ := json.Marshal(NewWorkorder{
		WorkorderID: workorderID,
	})

	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) insertResultsForJob(job *JobManual) (*storage.Workorders, error) {
	op, err := m.service.DB.GetOperation(job.OperationID, job.CarType)
	if err != nil {
		return nil, err
	}

	//points := []RoutingOperationPoint{}
	//err = json.Unmarshal([]byte(op.Points), &points)
	//if err != nil {
	//	return nil, errors.New("points is required")
	//}
	//
	//max_seq := 0
	//db_results := []storage.Results{}
	//for _, v := range points {
	//	if v.GroupSequence > max_seq {
	//		max_seq = v.GroupSequence
	//	}
	//	r := storage.Results{}
	//	r.PSet = v.PSet
	//	r.GroupSeq = v.GroupSequence
	//	r.OffsetX = v.X
	//	r.OffsetY = v.Y
	//	r.Seq = v.Seq
	//	r.MaxRedoTimes = v.MaxRedoTimes
	//	r.Stage = storage.RESULT_STAGE_INIT
	//	r.Result = storage.RESULT_NONE
	//	r.ControllerSN = job.Controller_SN
	//	r.UserID = job.UserID
	//	r.ToleranceMax = v.ToleranceMax
	//	r.ToleranceMin = v.ToleranceMin
	//	r.ToleranceMaxDegree = v.ToleranceMaxDegree
	//	r.ToleranceMinDegree = v.ToleranceMinDegree
	//	r.ConsuProductID = v.ConsuProductID
	//	r.Batch = fmt.Sprintf("%d/%d", v.GroupSequence, points[len(points)-1].GroupSequence)
	//	r.UpdateTime = time.Now()
	//	r.Count = 1
	//
	//	db_results = append(db_results, r)
	//}

	db_workorder := storage.Workorders{}
	db_workorder.Vin = job.Vin
	db_workorder.JobID = job.Job
	db_workorder.HMISN = job.HmiSN
	db_workorder.ProductID = op.ProductId
	db_workorder.WorkcenterID = op.WorkcenterID
	//db_workorder.MaxSeq = max_seq
	db_workorder.UserID = job.UserID
	db_workorder.MO_Model = job.CarType
	db_workorder.Mode = job.Mode
	db_workorder.UpdateTime = time.Now()
	db_workorder.WorkcenterCode = op.WorkcenterCode
	db_workorder.Consumes = op.Points

	err = m.service.DB.InsertWorkorder(&db_workorder, nil, false, false, true)

	return &db_workorder, err
}

func (m *Methods) insertResultsForPSet(pset *PSetManual) error {

	key := fmt.Sprintf("%s:%s:%s:%d:%d", pset.Vin, pset.CarType, pset.HmiSN, pset.ProductID, pset.WorkcenterID)

	r := storage.Results{}
	r.ExInfo = key
	r.PSet = pset.PSet
	r.Stage = storage.RESULT_STAGE_INIT

	ri := []storage.Results{r}
	//db_reuslt = append(db_reuslt, r)

	err := m.service.DB.DeleteResultsForJob(key)
	err = m.service.DB.InsertWorkorder(nil, &ri, false, false, true)

	return err
}

func (m *Methods) getNextWorkorder(ctx iris.Context) {
	var err error = nil
	hmi_sn := ctx.URLParam("hmi_sn")
	workcenterCode := ctx.URLParam("workcenter_code")

	if hmi_sn == "" && workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn or workcenter_code is required")
		return
	}

	workorder, err := m.service.DB.FindNextWorkorder(hmi_sn, workcenterCode)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.WriteString(err.Error())
			return
		}

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	nw := NextWorkorder{
		Vin:     workorder.Vin,
		LongPin: workorder.LongPin,
		Knr:     workorder.Knr,
		Model:   workorder.MO_Model,
		Lnr:     workorder.MO_Lnr,
	}

	body, _ := json.Marshal(nw)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

// 根据hmi序列号以及vin或knr取得工单
func (m *Methods) getWorkorder(ctx iris.Context) {
	//m.service.diag.Debug("getWorkorder start")

	var err error
	hmi_sn := ctx.URLParam("hmi_sn")
	workcenterCode := ctx.URLParam("workcenter_code")
	code := ctx.URLParam("code")

	if hmi_sn == "" && workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn or workcenter_code is required")
		return
	}

	if code == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("code is required")
		return
	}

	workorder, err := m.service.DB.FindWorkorder(hmi_sn, workcenterCode, code)
	if err != nil {
		// 通过odoo定位并创建工单
		body, e := m.service.ODOO.GetWorkorder("", hmi_sn, workcenterCode, code)
		if e != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("cannot find workorder")
			return
		} else {
			var odooWorkorders []odoo.ODOOWorkorder
			err := json.Unmarshal(body, &odooWorkorders)
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString(fmt.Sprintf("get workorder failed:%s", err.Error()))
				return
			}
			o, e := m.service.ODOO.CreateWorkorders(odooWorkorders)
			if e != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString(fmt.Sprintf("save workorder failed:%s", e.Error()))
				return
			} else {
				workorder = o[0] // 拥有获取的是第一张工单信息
			}
		}
	}

	//results, err := m.service.DB.FindResultsByWorkorder(workorder.Id)

	resp := Workorder{}
	resp.HMI_sn = workorder.HMISN
	resp.Workorder_id = workorder.Id
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.LongPin = workorder.LongPin
	resp.Status = workorder.Status
	resp.MaxOpTime = workorder.MaxOpTime
	//resp.WorkSheet = workorder.WorkSheet
	resp.Job = workorder.JobID
	resp.VehicleTypeImg = workorder.VehicleTypeImg
	resp.Lnr = workorder.MO_Lnr
	resp.Model = workorder.MO_Model

	op, err := m.service.DB.GetOperation(workorder.ImageOPID, workorder.MO_Model)
	if err == nil {
		resp.WorkSheet = op.Img
	}

	var consumes []odoo.ODOOConsume
	json.Unmarshal([]byte(workorder.Consumes), &consumes)

	for _, v := range consumes {
		r := Result{}
		r.PSet, _ = strconv.Atoi(v.PSet)
		r.GunSN = v.GunSN
		r.Controller_SN = v.ControllerSN
		r.X = v.X
		r.Y = v.Y
		r.MaxRedoTimes = v.Max_redo_times
		r.Seq = v.Seq
		r.GroupSeq = v.GroupSeq

		resp.Results = append(resp.Results, r)
	}

	var reasons []string
	if workorder.Status == "done" {
		// 工单已完成
		ctx.StatusCode(iris.StatusConflict)
		reasons = append(reasons, "done")
	}

	nextWorkorder, err := m.service.DB.FindNextWorkorder(hmi_sn, workcenterCode)
	if err == nil {
		if code != nextWorkorder.Knr && code != nextWorkorder.LongPin && code != nextWorkorder.Vin {
			// 车辆校验失败
			ctx.StatusCode(iris.StatusConflict)
			reasons = append(reasons, "conflict")
		}
	}

	resp.Reasons = reasons

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)

	//m.service.diag.Debug(fmt.Sprintf("getWorkorder finish with body:%s", string(body)))
}

func (m *Methods) getHealthz(ctx iris.Context) {

	ctx.StatusCode(iris.StatusNoContent)
	return
}

func (m *Methods) getHmiResults(ctx iris.Context) {

	result_id := ctx.URLParam("result_id")
	count := ctx.URLParam("count")

	if result_id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result_id is required")
		return
	}

	if count == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count is required")
		return
	}

	n_result_id, err := strconv.Atoi(result_id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result_id format error")
		return
	}

	n_count, err := strconv.Atoi(count)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count format error")
		return
	}

	db_result, err := m.service.DB.GetResult(int64(n_result_id), n_count)

	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString("result not found")
		return
	}

	ws_result := wsnotify.WSResult{}
	ws_result.Count = n_count
	ws_result.Result_id = int64(n_result_id)
	ws_result.Result = db_result.Result

	result_value := controller.ResultValue{}
	json.Unmarshal([]byte(db_result.ResultValue), &result_value)

	ws_result.TI = result_value.Ti
	ws_result.MI = result_value.Mi
	ws_result.WI = result_value.Wi

	body, _ := json.Marshal(ws_result)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
	ctx.StatusCode(iris.StatusOK)
}

func (m *Methods) getStatus(ctx iris.Context) {
	// 返回控制器状态

	sn := ctx.URLParam("controller_sn")

	sns := []string{}
	if sn != "" {
		vs := strings.Split(sn, ",")
		for _, v := range vs {
			if strings.TrimSpace(v) != "" {
				sns = append(sns, strings.TrimSpace(v))
			}
		}
	}

	status, err := m.service.AudiVw.GetControllersStatus(sns)

	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(err.Error())
		return
	} else {
		body, _ := json.Marshal(status)
		ctx.Header("content-type", "application/json")
		ctx.Write(body)
	}
}

func (m *Methods) putIOSet(ctx iris.Context) {
	io_set := IOSet{}
	err := ctx.ReadJSON(&io_set)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if io_set.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[io_set.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.IOSet(io_set.Controller_SN, &io_set.IOStatus)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}
}

func (m *Methods) putBarcodeTest(ctx iris.Context) {
	barcode := wsnotify.WSScanner{}
	err := ctx.ReadJSON(&barcode)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	str, _ := json.Marshal(barcode)
	m.service.OpenProtocol.WS.WSSendScanner(string(str))
}

func (m *Methods) putIOInputTest(ctx iris.Context) {
	inputs := openprotocol.IOMonitor{}
	err := ctx.ReadJSON(&inputs)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	str, _ := json.Marshal(inputs)
	m.service.OpenProtocol.WS.WSSendIOInput(string(str))
}

func (m *Methods) putResultTest(ctx iris.Context) {
	results := []wsnotify.WSResult{}
	err := ctx.ReadJSON(&results)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	str, _ := json.Marshal(results)
	m.service.OpenProtocol.WS.WSSend(wsnotify.WS_EVENT_RESULT, string(str))
}

func (m *Methods) postAK2(ctx iris.Context) {
	cr := controller.ControllerResult{}
	err := ctx.ReadJSON(&cr)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if cr.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if cr.GunSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("gun_sn is required")
		return
	}

	//if cr.PSet == 0 {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("pset is required")
	//	return
	//}

	if cr.Count == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count is required")
		return
	}

	if cr.Seq == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("seq is required")
		return
	}

	workorder, err := m.service.DB.GetWorkorder(cr.Workorder_ID, true)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workorder not found")
		return
	}

	workorder.Status = "done"
	m.service.DB.UpdateWorkorder(&workorder)

	_, exist := m.service.ControllerService.Controllers[cr.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	cr.NeedPushHmi = false
	cr.NeedPushAiis = true
	cr.Result = storage.RESULT_AK2
	cr.PSetDefine.Strategy = controller.STRATEGY_AK2

	m.service.ControllerService.Handle(&cr, nil)
}

func (m *Methods) putJobControll(ctx iris.Context) {
	jc := JobControl{}
	err := ctx.ReadJSON(&jc)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if jc.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[jc.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.JobControl(jc.Controller_SN, jc.Action)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("not supported")
		return
	}
}

func (m *Methods) getRoutingOpertions(ctx iris.Context) {
	code := ctx.Params().Get("code")
	if code == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workcenter code is required")
		return
	}

	carType := ctx.URLParam("carType")
	job, _ := strconv.Atoi(ctx.URLParam("job"))

	ro, err := m.service.DB.FindRoutingOperations(code, carType, job)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString("can not find RoutingOpertions")
		return
	}

	points := []RoutingOperationPoint{}
	json.Unmarshal([]byte(ro.Points), &points)

	rt_ro := RoutingOperation{}
	rt_ro.Points = points
	rt_ro.OperationID = ro.OperationID
	rt_ro.Job = ro.Job
	rt_ro.MaxOpTime = ro.MaxOpTime
	rt_ro.Name = ro.Name
	rt_ro.Img = ro.Img
	rt_ro.ProductId = ro.ProductId
	rt_ro.WorkcenterCode = ro.WorkcenterCode
	rt_ro.VehicleTypeImg = ro.VehicleTypeImg
	rt_ro.ProductType = ro.ProductType

	body, _ := json.Marshal(rt_ro)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) filterValue(filters string, key string, value interface{}) interface{} {
	if filters == "" || strings.Contains(filters, key) {
		return value
	}

	return nil
}

func (m *Methods) getLocalResults(ctx iris.Context) {
	hmi_sn := ctx.URLParam("hmi_sn")
	filters := ctx.URLParam("filters")
	limit, _ := strconv.Atoi(ctx.URLParam("limit"))

	results, err := m.service.DB.FindLocalResults(hmi_sn, limit)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	rt := []LocalResults{}
	sr := controller.ResultValue{}
	for _, v := range results {
		stime := v.Results.UpdateTime.Format("2006-01-02 15:04:05")
		dt, _ := time.Parse("2006-01-02 15:04:05", stime)

		json.Unmarshal([]byte(v.ResultValue), &sr)

		lr := LocalResults{
			HmiSN:        m.filterValue(filters, "hmi_sn", string(v.HMISN)),
			Vin:          m.filterValue(filters, "vin", string(v.Vin)),
			ControllerSN: m.filterValue(filters, "controller_sn", string(v.ControllerSN)),
			GunSN:        m.filterValue(filters, "gun_sn", string(v.GunSN)),
			Result:       m.filterValue(filters, "result", string(v.Result)),
			Torque:       m.filterValue(filters, "torque", float64(sr.Mi)),
			Angle:        m.filterValue(filters, "angle", float64(sr.Wi)),
			Spent:        m.filterValue(filters, "spent", float64(sr.Ti)),
			TimeStamp:    m.filterValue(filters, "timestamp", dt.Local()),
			Batch:        m.filterValue(filters, "batch", string(v.Batch)),
			VehicleType:  m.filterValue(filters, "vehicle_type", string(v.MO_Model)),
			JobID:        m.filterValue(filters, "job_id", int(v.JobID)),
			PSetID:       m.filterValue(filters, "pset_id", int(v.PSet)),
		}

		rt = append(rt, lr)
	}

	body, _ := json.Marshal(rt)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) listWorkorders(ctx iris.Context) {
	status := ctx.URLParam("status")
	hmi_sn := ctx.URLParam("hmi_sn")
	workcenterCode := ctx.URLParam("workcenter_code")

	if hmi_sn == "" && workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn or workcenter_code is required")
		return
	}

	workorders, err := m.service.DB.ListWorkorders(hmi_sn, workcenterCode, status)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("find workorders failed")
		return
	}

	rtWorkorders := []NextWorkorder{}
	for _, v := range workorders {
		rw := NextWorkorder{
			Vin:     v.Vin,
			Model:   v.MO_Model,
			Knr:     v.Knr,
			LongPin: v.LongPin,
			Lnr:     v.MO_Lnr,
		}

		rtWorkorders = append(rtWorkorders, rw)
	}

	body, _ := json.Marshal(rtWorkorders)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) testProtocol(ctx iris.Context) {

	tp := TestProtocol{
		ProtocolType: "op",
		Payload:      "0516006109981       010000020103工作组-test           04                         0500020600207000800000090001100000111122131141151161171181191200000000000210001002200000023000000240003232500720263600027000002802154290000030000003100000320003300034000350000003600000037000000380000003900000040000000410000000005420000043000004418C45596      452018-09-29:11:45:3346                   47                         481490250                      34051                         52                         53    540000005500000000005600570058",
	}
	err := ctx.ReadJSON(&tp)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("struct err")
		return
	}

	switch tp.ProtocolType {
	case "vw":
		m.service.AudiVw.Parse(tp.Payload)

	case "op":
		m.service.ControllerService.Controllers["0001"].(*openprotocol.Controller).Parse(tp.Payload)
	}
}
