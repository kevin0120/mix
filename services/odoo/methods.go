package odoo

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/storage"
	"strconv"
	"strings"
	"time"
)

type Methods struct {
	service *Service
}

// 创建工单
func (m *Methods) postWorkorders(ctx iris.Context) {
	var err error
	var workorders []ODOOWorkorder
	err = ctx.ReadJSON(&workorders)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	_, err = m.service.CreateWorkorders(workorders)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	} else {

		ctx.StatusCode(iris.StatusCreated)
		return
	}
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
	for i, v := range re_list {
		re_list[i] = strings.ToUpper(v)
	}

	bool_has_upload, err := strconv.ParseBool(has_upload)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload value error")
		return
	}

	resp := []ODOOResultSync{}
	results, _ := m.service.DB.FindUnuploadResults(bool_has_upload, re_list)
	target_results := map[int64]storage.Results{}
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
		odoo_result := ODOOResultSync{}

		odoo_result.Control_date = v.UpdateTime.Format(time.RFC3339)

		odoo_result.CURObjects = []aiis.CURObject{}

		curves, err := m.service.DB.ListCurvesByResult(v.ResultId)
		if err != nil {
			for _, c := range curves {
				cur_object := aiis.CURObject{}
				cur_object.File = c.CurveFile
				cur_object.OP = c.Count
				odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
			}
		}

		r := audi_vw.ResultValue{}
		json.Unmarshal([]byte(v.ResultValue), &r)

		pset := audi_vw.PSetDefine{}
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
	var up ResultPatch
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
