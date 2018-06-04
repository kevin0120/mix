package httpd

import (
	"github.com/kataras/iris"
	"github.com/masami10/rush/payload"
	"strings"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"fmt"
	"github.com/masami10/rush/services/storage"
)

type Methods struct {
	service	*Service
}

func (m *Methods) putPSets(ctx iris.Context) {

	var err error
	var pset payload.PSet
	err = ctx.ReadJSON(&pset)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	if pset.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
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
	_, err = m.service.DB.GetResult(pset.Result_id, 0)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	// 通过控制器设定程序
	//err = apiserver.CVI3.Service.PSet(pset.Controller_SN, pset.PSet, result.WorkorderID, pset.Result_id, pset.Count)
	//if err != nil {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString(err.Error())
	//	return
	//}
}

// 根据hmi序列号以及vin或knr取得工单
func (m *Methods) getWorkorder(ctx iris.Context) {
	hmi_sn := ctx.URLParam("hmi_sn")
	vin_or_knr := ctx.URLParam("vin_or_knr")

	if hmi_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn is required")
		return
	}

	if vin_or_knr == ""  {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("vin_or_knr is required")
		return
	}

	var vin = ""
	var knr = ""

	if strings.Contains(vin_or_knr, payload.KNR_KEY) {
		knr = vin_or_knr
	} else {
		vin = vin_or_knr
	}

	workorder, err := m.service.DB.FindWorkorder(hmi_sn, vin, knr)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	resp := payload.Workorder{}
	resp.HMI_sn = workorder.HMISN
	resp.PSet = workorder.PSet
	resp.Workorder_id = workorder.WorkorderID
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.Nut_total = workorder.NutTotal
	resp.Status = workorder.Status
	resp.WorkSheet = workorder.WorkSheet
	json.Unmarshal([]byte(workorder.ResultIDs), &resp.Result_ids)

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

// 创建工单
func (m *Methods) postWorkorders(ctx iris.Context) {
	var err error
	var workorders []payload.ODOOWorkorder
	err = ctx.ReadJSON(&workorders)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	var final_err error
	for _, v := range workorders {
		o := storage.Workorders{}
		o.Status = v.Status
		o.WorkorderID = v.ID
		o.PSet, _ = strconv.Atoi(v.PSet)
		o.HMISN = v.HMI.UUID
		o.Knr = v.KNR
		o.NutTotal = v.NutTotal
		o.Vin = v.VIN
		o.MaxOpTime = v.Max_op_time
		o.MaxRedoTimes = v.Max_redo_times

		ids, _ := json.Marshal(v.Result_IDs)
		o.ResultIDs = string(ids)

		e := m.service.DB.InsertWorkorder(o)
		if e != nil {
			final_err = e
		}
	}

	if final_err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(final_err.Error())

		return
	} else {

		ctx.StatusCode(iris.StatusCreated)
		return
	}
}

func (m *Methods) getDoc(ctx iris.Context) {
	f, _ := ioutil.ReadFile(m.service.ApiDoc)

	ctx.Header("content-type", "application/json")
	ctx.Write(f)
}

func (m *Methods) getResults(ctx iris.Context) {
	// 根据查询参数返回结果
	has_upload := ctx.URLParam("has_upload")
	if has_upload == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload is required")
		return
	}


	result := ctx.URLParams()["result"]

	if result == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result is required")
		return
	}
	re_list := strings.Split(result, ",")

	bool_has_upload, err := strconv.ParseBool(has_upload)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload value error")
		return
	}

	//list_result := []string{}
	//e := json.Unmarshal([]byte(result), &list_result)
	//if e != nil {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("result value error")
	//	return
	//}


	resp := []payload.ODOOResultSync{}
	results, _ := m.service.DB.FindUnuploadResults(bool_has_upload, re_list)
	target_results :=  map[int64]storage.Results{}
	for _, v := range results {
		tr, exist := target_results[v.ResultId]
		if exist {
			// 已存在
			if v.Count > tr.Count {
				target_results[v.ResultId] = v
			}
		} else {
			// 不存在
			target_results[v.ResultId] = v
		}
	}

	for _, v := range target_results {
		odoo_result := payload.ODOOResultSync{}
		stime := strings.Split(v.UpdateTime.Format("2006-01-02 15:04:05"), " ")
		odoo_result.Control_date = fmt.Sprintf("%sT%s+08:00", stime[0], stime[1])

		odoo_result.CURObjects = []payload.CURObject{}

		curves, err := m.service.DB.ListCurvesByResult(v.ResultId)
		if err != nil {
			for _, c := range curves {
				cur_object := payload.CURObject{}
				cur_object.File = c.CurveFile
				cur_object.OP = c.Count
				odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
			}
		}

		r := payload.ResultValue{}
		json.Unmarshal([]byte(v.ResultValue), &r)

		pset := payload.PSetDefine{}
		json.Unmarshal([]byte(v.PSetDefine), &pset)

		odoo_result.Measure_degree = r.Wi
		odoo_result.Measure_result = strings.ToLower(v.Result)
		odoo_result.Measure_t_don = r.Ti
		odoo_result.Measure_torque = r.Mi
		odoo_result.Op_time = v.Count
		odoo_result.Pset_m_max = pset.Mp
		odoo_result.Pset_m_min = pset.Mm
		odoo_result.Pset_m_target = pset.Ma
		odoo_result.Pset_m_threshold = pset.Ms
		odoo_result.Pset_strategy = pset.Strategy
		odoo_result.Pset_w_max = pset.Wp
		odoo_result.Pset_w_min = pset.Wm
		odoo_result.Pset_w_target = pset.Wa
		odoo_result.ID = v.ResultId

		resp = append(resp, odoo_result)
	}

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)

}

func (m *Methods) patchResult(ctx iris.Context) {
	id, err := strconv.Atoi(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	var e error
	var up payload.ResultPatch
	e = ctx.ReadJSON(&up)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	e = m.service.DB.UpdateResultByCount(int64(id), 0, up.HasUpload)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}
}

func (m *Methods) getHealthz(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNoContent)
	return
}

func (m *Methods) getStatus(ctx iris.Context) {
	// 返回控制器状态

	//sn := ctx.URLParam("controller_sn")
	//status, err := apiserver.CVI3.Service.GetControllersStatus(sn)
	//
	//if err != nil {
	//	ctx.StatusCode(iris.StatusNotFound)
	//	ctx.WriteString(err.Error())
	//	return
	//} else {
	//	body, _ := json.Marshal(status)
	//	ctx.Header("content-type", "application/json")
	//	ctx.Write(body)
	//}
}