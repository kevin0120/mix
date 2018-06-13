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
	var workorders []ODOOWorkorder
	err := ctx.ReadJSON(&workorders)

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

	hasUpload, err := strconv.ParseBool(has_upload)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload value error")
		return
	}

	queryMeasResultParam := strings.ToUpper(result)

	results, _ := m.service.DB.FindUnuploadResults(hasUpload, strings.Split(queryMeasResultParam, ","))
	targetResults := map[int64]storage.Results{}
	for _, v := range results {
		tr, exist := targetResults[v.ResultId]
		if exist {
			// 已存在
			if v.Count > tr.Count {
				targetResults[v.ResultId] = v
			}
		} else {
			// 不存在
			targetResults[v.ResultId] = v
		}
	}

	resp := make([]ODOOResultSync, len(targetResults))

	i := 0
	for _, v := range targetResults {
		odooResultSync := ODOOResultSync{}

		odooResultSync.Control_date = v.UpdateTime.Format(time.RFC3339)

		odooResultSync.CURObjects = []aiis.CURObject{}

		curves, err := m.service.DB.ListCurvesByResult(v.ResultId)
		if err != nil {
			for _, c := range curves {
				curObject := aiis.CURObject{}
				curObject.File = c.CurveFile
				curObject.OP = c.Count
				odooResultSync.CURObjects = append(odooResultSync.CURObjects, curObject)
			}
		}

		r := audi_vw.ResultValue{}
		json.Unmarshal([]byte(v.ResultValue), &r)

		pset := audi_vw.PSetDefine{}
		json.Unmarshal([]byte(v.PSetDefine), &pset)

		odooResultSync.Measure_degree = r.Wi
		odooResultSync.Measure_result = strings.ToLower(v.Result)
		odooResultSync.Measure_t_don = r.Ti
		odooResultSync.Measure_torque = r.Mi
		odooResultSync.Op_time = v.Count
		odooResultSync.Pset_m_max = pset.Mp
		odooResultSync.Pset_m_min = pset.Mm
		odooResultSync.Pset_m_target = pset.Ma
		odooResultSync.Pset_m_threshold = pset.Ms
		odooResultSync.Pset_strategy = pset.Strategy
		odooResultSync.Pset_w_max = pset.Wp
		odooResultSync.Pset_w_min = pset.Wm
		odooResultSync.Pset_w_target = pset.Wa
		odooResultSync.ID = v.ResultId

		resp[i] = odooResultSync
		i += 1
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

	up := ResultPatch{}
	e := ctx.ReadJSON(&up)
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
