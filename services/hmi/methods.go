package hmi

import (
	"github.com/kataras/iris"
	"encoding/json"
	"github.com/masami10/rush/services/storage"
)

type Methods struct {
	service	*Service
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
	result, err := m.service.DB.GetResult(pset.Result_id, 0)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	// 通过控制器设定程序
	err = m.service.AudiVw.PSet(pset.Controller_SN, pset.PSet, result.WorkorderID, pset.Result_id, pset.Count)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

// 根据hmi序列号以及vin或knr取得工单
func (m *Methods) getWorkorder(ctx iris.Context) {
	var err error
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

	var workorder storage.Workorders
	workorder, err = m.service.DB.FindWorkorder(hmi_sn, vin_or_knr, "")
	if err != nil {
		workorder, err = m.service.DB.FindWorkorder(hmi_sn, "", vin_or_knr)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
	}

	resp := Workorder{}
	resp.HMI_sn = workorder.HMISN
	resp.PSet = workorder.PSet
	resp.Workorder_id = workorder.WorkorderID
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.Nut_total = workorder.NutTotal
	resp.Status = workorder.Status
	resp.MaxRedoTimes = workorder.MaxRedoTimes
	resp.MaxOpTime = workorder.MaxOpTime

	json.Unmarshal([]byte(workorder.WorkSheet), &resp.WorkSheet)
	json.Unmarshal([]byte(workorder.ResultIDs), &resp.Result_ids)

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (m *Methods) getHealthz(ctx iris.Context) {

	ctx.StatusCode(iris.StatusNoContent)
	return
}

func (m *Methods) getStatus(ctx iris.Context) {
	// 返回控制器状态

	sn := ctx.URLParam("controller_sn")
	status, err := m.service.AudiVw.GetControllersStatus(sn)

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