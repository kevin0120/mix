package odoo

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"strconv"
	"strings"
	"time"
)

type Methods struct {
	service *Service
}

// TODO: 创建工单（工单模型/工步模型/结果模型）
// 创建工单
func (m *Methods) postWorkorders(ctx iris.Context) {

	//orderData, _ := ioutil.ReadAll(ctx.Request().Body)

	//return
	//
	var workorders []interface{}
	err := ctx.ReadJSON(&workorders)
	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	for _, v1 := range workorders {
		orderData, _ := json.Marshal(v1)
		m.service.diag.Debug(fmt.Sprintf("收到下發的工单: %s", string(orderData)))
		m.service.HandleWorkorder(orderData)
	}
	ctx.StatusCode(iris.StatusCreated)
	return
}

func (m *Methods) getResults(ctx iris.Context) {
	// 根据查询参数返回结果
	hasUploadParam := ctx.URLParam("has_upload")
	if hasUploadParam == "" {
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

	hasUpload, err := strconv.ParseBool(hasUploadParam)
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
			if v.Count > tr.Count || v.ResultId == 0 {
				targetResults[v.Id] = v
			}
		} else {
			// 不存在
			targetResults[v.Id] = v
		}
	}

	resp := make([]ODOOResultSync, len(targetResults))

	i := 0
	for _, v := range targetResults {
		odooResultSync := ODOOResultSync{}

		odooResultSync.LocalID = v.Id

		odooResultSync.Control_date = v.UpdateTime.Format(time.RFC3339)

		odooResultSync.CURObjects = []aiis.CURObject{}

		odooResultSync.Batch = v.Batch

		curves, err := m.service.DB.ListCurvesByResult(v.ResultId)
		if err == nil {
			for _, c := range curves {
				curObject := aiis.CURObject{}
				curObject.File = c.CurveFile
				curObject.OP = c.Count
				odooResultSync.CURObjects = append(odooResultSync.CURObjects, curObject)
			}
		}

		r := tightening_device.ResultValue{}
		json.Unmarshal([]byte(v.ResultValue), &r)

		pset := tightening_device.PSetDefine{}
		json.Unmarshal([]byte(v.PSetDefine), &pset)

		if v.Result == storage.RESULT_OK {
			odooResultSync.Final_pass = tightening_device.RESULT_PASS
			if v.Count == 1 {
				odooResultSync.One_time_pass = tightening_device.RESULT_PASS
			} else {
				odooResultSync.One_time_pass = tightening_device.RESULT_FAIL
			}

			if (r.Mi >= v.ToleranceMin && r.Mi <= v.ToleranceMax) &&
				(r.Wi >= v.ToleranceMinDegree && r.Wi <= v.ToleranceMaxDegree) {
				odooResultSync.QualityState = tightening_device.RESULT_PASS
				odooResultSync.ExceptionReason = ""
			} else {
				odooResultSync.QualityState = tightening_device.RESULT_EXCEPTION
				odooResultSync.ExceptionReason = tightening_device.RESULT_EXCEPTION
			}

		} else {
			odooResultSync.Final_pass = tightening_device.RESULT_FAIL
			odooResultSync.One_time_pass = tightening_device.RESULT_FAIL

			if (r.Mi >= v.ToleranceMin && r.Mi <= v.ToleranceMax) &&
				(r.Wi >= v.ToleranceMinDegree && r.Wi <= v.ToleranceMaxDegree) {

				odooResultSync.QualityState = tightening_device.RESULT_EXCEPTION
				odooResultSync.ExceptionReason = tightening_device.RESULT_EXCEPTION
			} else {
				odooResultSync.QualityState = tightening_device.RESULT_FAIL
				odooResultSync.ExceptionReason = ""
			}

		}

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

	up := aiis.ResultPatch{}
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

func (m *Methods) deleteRoutingOpertions(ctx iris.Context) {
	rds := []storage.RoutingOperationDelete{}
	e := ctx.ReadJSON(&rds)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	rds_str, _ := json.Marshal(rds)
	m.service.diag.Debug(fmt.Sprintf("remove local operations:%s", rds_str))

	m.service.DB.DeleteRoutingOperations(rds)

	ctx.StatusCode(iris.StatusNoContent)
}

func (m *Methods) deleteAllRoutingOpertions(ctx iris.Context) {
	err := m.service.DB.DeleteAllRoutingOperations()

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

func (m *Methods) putSyncRoutingOpertions(ctx iris.Context) {

	//orderData, _ := ioutil.ReadAll(ctx.Request().Body)
	//m.service.diag.Debug(fmt.Sprintf("收到下發的工单: %s", string(orderData)))

	ro := RoutingOperation{}
	e := ctx.ReadJSON(&ro)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	points, _ := json.Marshal(ro.Points)

	db_ro, err := m.service.DB.GetRoutingOperations(ro.Name, ro.ProductType)

	db_ro.Points = string(points)
	db_ro.VehicleTypeImg = ro.VehicleTypeImg
	db_ro.WorkcenterCode = ro.WorkcenterCode
	db_ro.ProductType = ro.ProductType
	db_ro.ProductId = ro.ProductId
	db_ro.Img = ro.Img
	db_ro.Name = ro.Name
	db_ro.MaxOpTime = ro.MaxOpTime
	db_ro.Job = ro.Job
	db_ro.WorkcenterID = ro.WorkcenterID
	db_ro.Tigntening_step_ref = ro.Tigntening_step_ref
	db_ro.ProductTypeImage = ro.ProductTypeImage

	if err != nil {
		// 新增
		db_ro.OperationID = ro.OperationID
		m.service.DB.Store(db_ro)
	} else {
		// 更新
		m.service.DB.UpdateRoutingOperations(&db_ro)
	}
}

func (m *Methods) postMaintenance(ctx iris.Context) {
	maintanence := Maintenance{}

	err := ctx.ReadJSON(&maintanence)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	msg, _ := json.Marshal(maintanence)
	m.service.WS.NotifyAll(wsnotify.WS_EVENT_MAINTENANCE, string(msg))
}
