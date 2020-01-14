package odoo

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
)

// 创建工单
func (s *Service) postWorkorders(ctx iris.Context) {

	var workorders []interface{}
	err := ctx.ReadJSON(&workorders)
	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	for _, v1 := range workorders {
		orderData, _ := json.Marshal(v1)
		s.diag.Debug(fmt.Sprintf("收到下發的工单: %s", string(orderData)))
		s.handleWorkorder(orderData)
	}
	ctx.StatusCode(iris.StatusCreated)
	return
}

func (s *Service) deleteRoutingOpertions(ctx iris.Context) {
	var rds []storage.RoutingOperationDelete
	e := ctx.ReadJSON(&rds)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(e.Error())
		return
	}

	rdsStr, _ := json.Marshal(rds)
	s.diag.Debug(fmt.Sprintf("remove local operations:%s", rdsStr))

	if err := s.storageService.DeleteRoutingOperations(rds); err != nil {
		s.diag.Error("Delete Routing Operations Failed ", err)
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) deleteAllRoutingOpertions(ctx iris.Context) {
	err := s.storageService.DeleteAllRoutingOperations()

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

func (s *Service) putSyncRoutingOpertions(ctx iris.Context) {

	ro := RoutingOperation{}
	e := ctx.ReadJSON(&ro)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(e.Error())
		return
	}

	points, _ := json.Marshal(ro.Points)

	dbRo, err := s.storageService.GetRoutingOperations(ro.Name, ro.ProductType, ro.TignteningStepRef)
	if err != nil {
		dbRo = storage.RoutingOperations{
			OperationID: ro.OperationID,
		}
	}

	dbRo.Points = string(points)
	dbRo.WorkcenterCode = ro.WorkcenterCode
	dbRo.ProductType = ro.ProductType
	dbRo.ProductId = ro.ProductId
	dbRo.Img = ro.Img
	dbRo.Name = ro.Name
	dbRo.MaxOpTime = ro.MaxOpTime
	dbRo.Job = ro.Job
	dbRo.WorkcenterID = ro.WorkcenterID
	dbRo.TighteningStepRef = ro.TignteningStepRef
	dbRo.ProductTypeImage = ro.ProductTypeImage

	if err != nil {
		// 新增
		if err := s.storageService.Store(dbRo); err != nil {
			s.diag.Error("Create Routing Operation Failed ", err)
		}
	} else {
		// 更新
		if err := s.storageService.UpdateRoutingOperations(&dbRo); err != nil {
			s.diag.Error("Update Routing Operation Failed ", err)
		}
	}
}

func (s *Service) postMaintenance(ctx iris.Context) {
	maintanence := Maintenance{}

	err := ctx.ReadJSON(&maintanence)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}

	s.doDispatch(dispatcherbus.DispatcherMaintenanceInfo, maintanence)
}
