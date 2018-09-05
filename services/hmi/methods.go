package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
	"strings"
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

	//if te.GunSN == "" {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("gun_sn is required")
	//	return
	//}

	// 通过控制器设定程序
	c, exist := m.service.ControllerService.Controllers[te.Controller_SN]
	if !exist {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller not found")
		return
	}

	switch c.Protocol() {
	case controller.AUDIPROTOCOL:
		err = m.service.AudiVw.ToolControl(te.Controller_SN, te.Enable)
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.ToolControl(te.Controller_SN, te.Enable)

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

	if pset.UserID == 0 {
		pset.UserID = DEFAULT_USER_ID
	}

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

	if pset.Result_id == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result_id is required")
		return
	}

	// 检测count
	if pset.Count < 1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("tightning count should be greater than 0")
		return
	}

	// 检测结果id
	result, err := m.service.DB.GetResult(pset.Result_id, 0)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
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
	case controller.AUDIPROTOCOL:
		err = m.service.AudiVw.PSet(pset.Controller_SN, pset.PSet, result.WorkorderID, pset.Result_id, pset.Count, pset.UserID)
	case controller.OPENPROTOCOL:
		err = m.service.OpenProtocol.PSet(pset.Controller_SN, pset.PSet, pset.Result_id, pset.Count, pset.UserID)

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

		ex_info := fmt.Sprintf("%25s%25s%25s%25d", pset.Vin, pset.HmiSN, pset.CarType, pset.UserID)
		err = m.service.OpenProtocol.PSetManual(pset.Controller_SN, pset.PSet, pset.UserID, ex_info, pset.Count)

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

	if job.Job == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("job is required")
		return
	}

	if job.CarType == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("car type is required")
		return
	}

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

	switch c.Protocol() {
	case controller.OPENPROTOCOL:
		err = m.insertResultsForJob(&job)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
		//vin-cartype-hmisn-userid
		ex_info := fmt.Sprintf("%25s%25s%25s%25d", job.Vin, job.HmiSN, job.CarType, job.UserID)
		err = m.service.OpenProtocol.JobSetManual(job.Controller_SN, job.Job, job.UserID, ex_info)

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

func (m *Methods) insertResultsForJob(job *JobManual) error {
	if len(job.Points) == 0 {
		return errors.New("points is required")
	}

	key := fmt.Sprintf("%s:%s:%s:%d:%d", job.Vin, job.CarType, job.HmiSN, job.ProductID, job.WorkcenterID)

	db_results := []storage.Results{}
	for _, v := range job.Points {
		r := storage.Results{}
		r.ExInfo = key
		r.PSet = v.PSet
		r.OffsetX = v.X
		r.OffsetY = v.Y
		r.Seq = v.Seq
		r.MaxRedoTimes = v.MaxOpTime
		r.Stage = storage.RESULT_STAGE_INIT

		db_results = append(db_results, r)
	}

	err := m.service.DB.DeleteResultsForJob(key)
	err = m.service.DB.InsertWorkorder(nil, &db_results, false)

	return err
}

func (m *Methods) insertResultsForPSet(pset *PSetManual) error {

	key := fmt.Sprintf("%s:%s:%s:%d:%d", pset.Vin, pset.CarType, pset.HmiSN, pset.ProductID, pset.WorkcenterID)

	r := storage.Results{}
	r.ExInfo = key
	r.PSet = pset.PSet
	r.Stage = storage.RESULT_STAGE_INIT

	db_reuslt := []storage.Results{}
	db_reuslt = append(db_reuslt, r)

	err := m.service.DB.DeleteResultsForJob(key)
	err = m.service.DB.InsertWorkorder(nil, &db_reuslt, false)

	return err
}

// 根据hmi序列号以及vin或knr取得工单
func (m *Methods) getWorkorder(ctx iris.Context) {

	var err error
	hmi_sn := ctx.URLParam("hmi_sn")
	code := ctx.URLParam("code")

	if hmi_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn is required")
		return
	}

	if code == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("code is required")
		return
	}

	workorder, err := m.service.DB.FindWorkorder(hmi_sn, code)
	if err != nil {
		// 通过odoo定位并创建工单
		body, e := m.service.ODOO.GetWorkorder(m.service.SN, hmi_sn, code)
		if e != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("cannot find workorder")
			return
		} else {
			var odooWorkorders []odoo.ODOOWorkorder
			json.Unmarshal(body, &odooWorkorders)
			o, e := m.service.ODOO.CreateWorkorders(odooWorkorders)
			if e != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString("save workorder failed")
				return
			} else {
				workorder = o[0]
			}
		}
	}

	results, err := m.service.DB.FindResultsByWorkorder(workorder.WorkorderID)

	resp := Workorder{}
	resp.HMI_sn = workorder.HMISN
	resp.Workorder_id = workorder.WorkorderID
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.LongPin = workorder.LongPin
	resp.Status = workorder.Status
	resp.MaxOpTime = workorder.MaxOpTime
	resp.WorkSheet = workorder.WorkSheet
	resp.Job = workorder.JobID

	for _, v := range results {
		r := Result{}
		r.PSet = v.PSet
		r.GunSN = v.GunSN
		r.ID = v.ResultId
		r.Controller_SN = v.ControllerSN
		r.X = v.OffsetX
		r.Y = v.OffsetY
		r.MaxRedoTimes = v.MaxRedoTimes

		resp.Results = append(resp.Results, r)
	}

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
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
