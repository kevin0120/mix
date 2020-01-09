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
		ctx.WriteString(err.Error())
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
	rds := []storage.RoutingOperationDelete{}
	e := ctx.ReadJSON(&rds)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	rds_str, _ := json.Marshal(rds)
	s.diag.Debug(fmt.Sprintf("remove local operations:%s", rds_str))

	s.storageService.DeleteRoutingOperations(rds)

	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) deleteAllRoutingOpertions(ctx iris.Context) {
	err := s.storageService.DeleteAllRoutingOperations()

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

func (s *Service) putSyncRoutingOpertions(ctx iris.Context) {

	ro := RoutingOperation{}
	e := ctx.ReadJSON(&ro)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	points, _ := json.Marshal(ro.Points)

	db_ro, err := s.storageService.GetRoutingOperations(ro.Name, ro.ProductType)

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
	db_ro.TighteningStepRef = ro.Tigntening_step_ref
	db_ro.ProductTypeImage = ro.ProductTypeImage

	if err != nil {
		// 新增
		db_ro.OperationID = ro.OperationID
		s.storageService.Store(db_ro)
	} else {
		// 更新
		s.storageService.UpdateRoutingOperations(&db_ro)
	}
}

func (s *Service) postMaintenance(ctx iris.Context) {
	maintanence := Maintenance{}

	err := ctx.ReadJSON(&maintanence)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	s.dispatcherBus.Dispatch(dispatcherbus.DispatcherMaintenanceInfo, maintanence)
}
