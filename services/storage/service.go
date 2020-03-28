package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
	OpenEngineSuccess(info string)
	Close()
	Closed()
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	eng         *xorm.Engine
	validator   *validator.Validate
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
	}

	s.configValue.Store(c)
	s.validator = validator.New()

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}

	info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Url,
		c.DBName)
	engine, err := xorm.NewEngine("postgres", info)

	engine.DatabaseTZ = time.UTC
	engine.TZLocation = time.UTC

	if err != nil {
		return errors.Wrapf(err, "Create postgres engine fail")
	}

	_, err = engine.IsTableExist("workorders")
	if err == nil {
		if err := engine.Sync2(new(Workorders)); err != nil {
			return errors.Wrapf(err, "Create Table Workorders fail")
		}
	}

	_, err = engine.IsTableExist("results")
	if err == nil {
		if err := engine.Sync2(new(Results)); err != nil {
			return errors.Wrapf(err, "Create Table Results fail")
		}
	}

	_, err = engine.IsTableExist("curves")
	if err == nil {
		if err := engine.Sync2(new(Curves)); err != nil {
			return errors.Wrapf(err, "Create Table Curves fail")
		}
	}

	_, err = engine.IsTableExist("controllers")
	if err == nil {
		if err := engine.Sync2(new(Controllers)); err != nil {
			return errors.Wrapf(err, "Create Table Controllers fail")
		}
	}

	_, err = engine.IsTableExist("guns")
	if err == nil {
		if err := engine.Sync2(new(Tools)); err != nil {
			return errors.Wrapf(err, "Create Table Tools fail")
		}
	}

	_, err = engine.IsTableExist("routing_operations")
	if err := engine.Sync2(new(RoutingOperations)); err != nil {
		return errors.Wrapf(err, "Create Table RoutingOperations fail")
	}

	_, err = engine.IsTableExist("steps")
	if err == nil {
		if err := engine.Sync2(new(Steps)); err != nil {
			return errors.Wrapf(err, "Create Table Steps fail")
		}
	}

	engine.SetMaxOpenConns(c.MaxConnects) // always success

	s.eng = engine

	s.diag.OpenEngineSuccess(info)

	go s.DropTableManage() //启动drop数据协程
	return nil
}

func (s *Service) Close() error {
	s.diag.Close()
	if s.eng != nil {
		s.eng.Close()
	}

	s.diag.Closed()

	return nil
}

func (s *Service) Store(data interface{}) error {

	session := s.eng.NewSession()
	defer session.Close()

	// add Begin() before any action
	err := session.Begin()
	_, err = session.Insert(data)
	if err != nil {
		session.Rollback()
		return errors.Wrapf(err, "store data fail")
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil
}

func (s *Service) FindUnuploadResults(result_upload bool, result []string) ([]Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.has_upload = ?", result_upload).And("r.stage = ?", "final").In("r.result", result)

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (s *Service) ListUnUploadResults() ([]Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.has_upload = ?", false).And("r.curve_file != ?", "")

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (s *Service) GetTool(serial string) (Tools, error) {

	var tools Tools

	rt, err := s.eng.Alias("g").Where("g.serial = ?", serial).Get(&tools)

	if err != nil {
		return tools, err
	} else {
		if !rt {
			return tools, errors.New("found gun failed")
		} else {
			return tools, nil
		}
	}
}

func (s *Service) UpdateTool(gun *Tools) error {
	g, err := s.GetTool(gun.Serial)
	if err == nil {
		// update
		_, err := s.eng.Id(g.Id).Update(gun)
		if err != nil {
			return err
		}
	} else {
		// insert
		_, err := s.eng.Insert(gun)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Service) GetOperation(id int64, model string) (RoutingOperations, error) {
	var op RoutingOperations

	rt, err := s.eng.Alias("g").Where("g.operation_id = ?", id).And("g.product_type = ?", model).Get(&op)

	if err != nil {
		return op, err
	} else {
		if !rt {
			return op, errors.New("found op failed")
		} else {
			return op, nil
		}
	}
}

func (s *Service) UpdateCurve(curve *Curves) (*Curves, error) {

	sql := "update `curves` set has_upload = ?, curve_file = ?, curve_data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		curve.HasUpload, curve.CurveFile, curve.CurveData, curve.Id)

	if err != nil {
		return curve, err
	} else {
		return curve, nil
	}
}

func (s *Service) ListCurvesByResult(result_id int64) ([]Curves, error) {
	var curves []Curves

	e := s.eng.Alias("c").Where("c.result_id = ?", result_id).Find(&curves)
	if e != nil {
		return curves, e
	} else {
		return curves, nil
	}
}

func (s *Service) ListUnuploadCurves() ([]Curves, error) {
	var curves []Curves

	e := s.eng.Alias("c").Where("c.has_upload = ?", false).And("c.tightening_id != ?", "").Find(&curves)
	if e != nil {
		return curves, e
	} else {
		return curves, nil
	}
}

func (s *Service) InsertWorkorder(workorder *Workorders, results *[]Results, checkWorkorder bool, checkResult bool, rawid bool) error {

	session := s.eng.NewSession()
	defer session.Close()

	err := session.Begin()
	// 执行事务

	// 保存新工单
	var raw_workorderid int64
	if workorder != nil {
		if checkWorkorder {
			has, _ := s.WorkorderExists(workorder.WorkorderID)
			if has {
				// 忽略存在的工单
				return nil
			}
		}

		_, err = session.Insert(workorder)
		if err != nil {
			session.Rollback()
			return errors.Wrapf(err, "store data fail")
		} else {
			s.diag.Debug(fmt.Sprintf("new workorder:%d", workorder.WorkorderID))
			raw_workorderid = workorder.Id
		}
	}

	// 预保存结果
	if results != nil {
		for _, v := range *results {

			if checkResult {
				has, _ := s.ResultExists(v.ResultId)
				if has {
					continue
				}
			}

			if rawid {
				v.WorkorderID = raw_workorderid
			}

			_, err = session.Insert(v)
			if err != nil {
				session.Rollback()
				return errors.Wrapf(err, "store data fail")
			} else {
				s.diag.Debug(fmt.Sprintf("new result:%d", v.ResultId))
			}
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil
}

func (s *Service) DeleteResultsForJob(key string) error {
	sql := fmt.Sprintf("delete from `results` where exinfo = '%s'", key)
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) WorkorderExists(id int64) (bool, error) {

	has, err := s.eng.Exist(&Workorders{WorkorderID: id})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (s *Service) ResultExists(id int64) (bool, error) {

	has, err := s.eng.Exist(&Results{ResultId: id})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (s *Service) GetResultByID(id int64) (*Results, error) {
	result := Results{}
	rt, err := s.eng.Alias("r").Where("r.id = ?", id).Limit(1).Get(&result)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, errors.New("result does not exist")
		} else {
			return &result, nil
		}
	}
}

func (s *Service) ListWorkorders(hmi_sn string, workcenterCode string, status string) ([]Workorders, error) {
	var workorders []Workorders

	sql := fmt.Sprintf("select * from workorders where hmi_sn = '%s' or workcenter_code = '%s'", hmi_sn, workcenterCode)
	sql = sql + fmt.Sprintf(" and workorders.x_workorder_id > 0 ")
	if status != "" {
		sql = sql + fmt.Sprintf(" and workorders.status = '%s'", status)
	}

	sql = sql + " order by update_time, lnr asc"

	err := s.eng.SQL(sql).Find(&workorders)

	return workorders, err
}

func (s *Service) GetWorkOrder(id int64, raw bool) (Workorders, error) {

	var workorder Workorders

	key := "x_workorder_id"
	if raw {
		key = "id"
	}

	rt, err := s.eng.Alias("w").Where(fmt.Sprintf("w.%s = ?", key), id).Get(&workorder)

	if err != nil {
		return workorder, err
	} else {
		if !rt {
			return workorder, errors.New("workorder does not exist")
		} else {
			return workorder, nil
		}
	}
}

func (s *Service) GetStep(id int64) (Steps, error) {

	var step Steps

	rt, err := s.eng.Alias("w").Where("w.id = ?", id).Get(&step)

	if err != nil {
		return step, err
	} else {
		if !rt {
			return step, errors.New("Step does not exist")
		} else {
			return step, nil
		}
	}
}

func (s *Service) FindWorkorder(hmi_sn string, workcenter_code string, code string) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.hmi_sn = ? or w.workcenter_code = ?", hmi_sn, workcenter_code).And("w.long_pin = ? or w.vin = ? or w.knr = ?", code, code, code).And("w.x_workorder_id > ?", 0).Get(&workorder)

	if workorder.LongPin == "" {
		return workorder, errors.New("workorder does not exist")
	} else if !rt || err != nil {
		return workorder, errors.New("workorder does not exist")
	}
	return workorder, nil
}

func (s *Service) GetOperationByID(opid int64) (RoutingOperations, error) {

	var routingOp RoutingOperations

	rt, err := s.eng.Alias("r").Where("r.id = ?", opid).Get(&routingOp)

	if err != nil {
		return routingOp, err
	} else {
		if !rt {
			return routingOp, errors.New("workorder does not exist")
		} else {
			return routingOp, nil
		}
	}
}

func (s *Service) FindNextWorkorder(hmi_sn string, workcenter_code string) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.hmi_sn = ? or w.workcenter_code = ?", hmi_sn, workcenter_code).And("w.status = ?", "ready").And("w.x_workorder_id > ?", 0).Asc("w.update_time").Asc("w.lnr").Get(&workorder)

	if err != nil {
		return workorder, err
	} else {
		if !rt {
			return workorder, errors.New("workorder not found")
		} else {
			return workorder, nil
		}
	}
}

func (s *Service) UpdateWorkorderUserID(id int64, userID int64) error {
	sql := "update `workorders` set user_id = ? where id = ?"
	_, err := s.eng.Exec(sql,
		userID,
		id)

	return err
}

func (s *Service) UpdateResult(result *Results) (int64, error) {

	sql := "update `results` set controller_sn = ?, result = ?, has_upload = ?, stage = ?, update_time = ?, pset_define = ?, result_value = ?, count = ?, batch = ?, gun_sn = ?, spent = ?, tightening_id = ? where id = ?"
	r, err := s.eng.Exec(sql,
		result.ControllerSN,
		result.Result,
		result.HasUpload,
		result.Stage,
		result.UpdateTime,
		result.PSetDefine,
		result.ResultValue,
		result.Count,
		result.Batch,
		result.ToolSN,
		result.Spent,
		result.TighteningID,
		result.Id)

	if err != nil {

		return 0, errors.Wrapf(err, "Update result fail for id : %d", result.Id)
	} else {
		id, _ := r.RowsAffected()
		return id, nil
	}
}

func (s *Service) UpdateWorkorder(workorder *Workorders) (*Workorders, error) {

	sql := "update `workorders` set status = ? where id = ?"
	_, err := s.eng.Exec(sql,
		workorder.Status,
		workorder.Id)

	if err != nil {
		return workorder, err
	} else {
		return workorder, nil
	}
}

func (s *Service) UpdateResultByCount(id int64, count int, flag bool) error {

	var err error
	if count > 0 {
		sql := "update `results` set has_upload = ? where id = ? and count = ?"
		_, err = s.eng.Exec(sql, flag, id, count)
	} else {
		sql := "update `results` set has_upload = ? where id = ?"
		_, err = s.eng.Exec(sql, flag, id)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) DeleteRoutingOperations(rds []RoutingOperationDelete) error {
	for _, v := range rds {
		sql := fmt.Sprintf("delete from `routing_operations` where operation_id = %d and product_type = '%s'", v.OperationID, v.ProductType)
		s.eng.Exec(sql)
	}

	return nil
}

func (s *Service) DeleteAllRoutingOperations() error {
	_, err := s.eng.Exec("delete from `routing_operations`")

	return err
}

func (s *Service) DeleteInvalidResults(keep time.Time) error {
	sql := fmt.Sprintf("delete from `results` where has_upload = true and update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) DeleteInvalidCurves(keep time.Time) error {

	sql := fmt.Sprintf("delete from `curves` where has_upload = true and update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) DeleteInvalidWorkorders(keep time.Time) error {

	sql := fmt.Sprintf("delete from `workorders` where status = 'done' and update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) InitWorkorderForJob(workorder_id int64) error {
	sql := "update `results` set result = ?, stage = ?, count = ? where x_workorder_id = ?"
	_, err := s.eng.Exec(sql, RESULT_NONE, RESULT_STAGE_INIT, 0, workorder_id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) FindTargetResultForJob(workorder_id int64) (Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", workorder_id).And("r.stage = ?", RESULT_STAGE_INIT).OrderBy("r.seq")

	e := ss.Find(&results)

	if e != nil {
		return Results{}, e
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return Results{}, errors.New("result not found")
		}
	}
}

func (s *Service) FindTargetResultForJobManual(raw_workorder_id int64) (Results, error) {
	var results []Results

	//ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", raw_workorder_id).OrderBy("r.update_time").OrderBy("r.seq").OrderBy("r.count").Desc("r.update_time").Desc("r.seq").Desc("r.count")

	sql := fmt.Sprintf("select * from results where results.x_workorder_id = %d order by results.update_time desc, results.tightening_id desc, results.seq desc, results.count desc", raw_workorder_id)

	e := s.eng.SQL(sql).Find(&results)

	if e != nil {
		return Results{}, e
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return Results{}, errors.New("result not found")
		}
	}
}

func (s *Service) CreateController(controller_sn string) (Controllers, error) {
	var controller Controllers

	rt, err := s.eng.Alias("c").Where("c.controller_sn = ?", controller_sn).Get(&controller)

	if err != nil {
		return controller, err
	} else {
		if !rt {
			// 创建
			controller.SN = controller_sn
			controller.LastID = "0"
			err = s.Store(controller)
			if err != nil {
				return controller, err
			} else {
				rt, err = s.eng.Alias("c").Where("c.controller_sn = ?", controller_sn).Get(&controller)
				return controller, nil
			}
		} else {
			return controller, nil
		}
	}
}

func (s *Service) UpdateTightning(id int64, last_id string) error {
	sql := "update `controllers` set last_id = ? where id = ?"
	_, err := s.eng.Exec(sql,
		last_id,
		id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) ResetTightning(controller_sn string) error {
	sql := "update `controllers` set last_id = ? where controller_sn = ?"
	_, err := s.eng.Exec(sql,
		"0",
		controller_sn)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) UpdateRoutingOperations(ro *RoutingOperations) error {
	sql := "update `routing_operations` set job = ?, max_op_time = ?, name = ?, img = ?, product_id = ?, product_type = ?, workcenter_code = ?, vehicle_type_img = ?, points = ?, workcenter_id = ? where id = ?"
	_, err := s.eng.Exec(sql,
		ro.Job,
		ro.MaxOpTime,
		ro.Name,
		ro.Img,
		ro.ProductId,
		ro.ProductType,
		ro.WorkcenterCode,
		ro.VehicleTypeImg,
		ro.Points,
		ro.WorkcenterID,
		ro.Id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) GetRoutingOperations(name string, model string, step string) (RoutingOperations, error) {

	var ro RoutingOperations

	rt, err := s.eng.Alias("r").Where("r.name = ?", name).And("r.product_type = ?", model).And("r.tightening_step_ref = ?", step).Get(&ro)

	if err != nil {
		return ro, err
	} else {
		if !rt {
			return ro, errors.New("found RoutingOperations failed")
		} else {
			return ro, nil
		}
	}
}

func (s *Service) FindRoutingOperations(workcenter_code string, cartype string, job int) (RoutingOperations, error) {

	var ros []RoutingOperations

	ss := s.eng.Alias("r").Where("r.workcenter_code = ?", workcenter_code).And("r.product_type = ? or r.job = ?", cartype, job)

	e := ss.Find(&ros)

	if e != nil {
		return RoutingOperations{}, e
	} else {
		if len(ros) > 0 {
			return ros[0], nil
		} else {
			return RoutingOperations{}, errors.New("result not found")
		}
	}
}

func (s *Service) FindLocalResults(hmi_sn string, limit int) ([]ResultsWorkorders, error) {
	var results []ResultsWorkorders

	sql := "select * from results, workorders where results.x_workorder_id = workorders.id "
	if hmi_sn != "" {
		sql = sql + fmt.Sprintf(" and workorders.hmi_sn = '%s'", hmi_sn)
	}

	sql = sql + " order by results.update_time desc"

	err := s.eng.SQL(sql).Find(&results)

	if limit > 0 {
		if limit < len(results) {
			return results[0:limit], err
		}
	}

	return results, err
}

func (s *Service) IsMultiResult(workorderID int64, batch string) bool {
	var results []Results

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", workorderID).And("r.batch = ?", batch)

	e := ss.Find(&results)

	if e != nil {
		return false
	} else {
		if len(results) > 1 {
			return true
		} else {
			return false
		}
	}
}

func (s *Service) GetController(sn string) (interface{}, error) {

	var rt_controller Controllers

	rt, err := s.eng.Alias("c").Where("c.controller_sn = ?", sn).Get(&rt_controller)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, nil
		} else {
			return rt_controller, nil
		}
	}
}

func (s *Service) ResetResult(id int64) error {
	sql := "update `results` set result = ?, has_upload = ?, stage = ?, count = ? where id = ?"

	_, err := s.eng.Exec(sql,
		RESULT_NONE,
		false,
		RESULT_STAGE_INIT,
		1,
		id)

	return err
}

func (s *Service) DeleteCurvesByResult(id int64) error {

	sql := fmt.Sprintf("delete from `curves` where result_id = %d", id)
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) DropTableManage() error {
	c := s.Config()
	for {
		start := time.Now()

		keep := time.Now().Add(time.Duration(c.DataKeep) * -1)

		// 清理过期结果
		s.DeleteInvalidResults(keep)

		// 清理过期波形
		s.DeleteInvalidCurves(keep)

		// 清理过期工单
		s.DeleteInvalidWorkorders(keep)

		diff := time.Since(start) // 执行的间隔时间

		time.Sleep(time.Duration(c.VacuumPeriod) - diff)
	}
}

func (s *Service) Workorders(par []byte) ([]Workorders, error) {
	orderPar := WorkorderListPar{}
	err := json.Unmarshal(par, &orderPar)
	if err != nil {
		return nil, err
	}

	q := s.eng.Alias("w")
	if orderPar.Status != "" {
		q = q.Where("w.status = ?", orderPar.Status)
	} else {
		q = q.Where("w.status != ?", orderPar.Status)
	}

	if orderPar.Time_from != "" {
		q = q.And("w.date_planned_start >= ?", orderPar.Time_from)
	}
	if orderPar.Time_to != "" {
		q = q.And("w.date_planned_complete <= ?", orderPar.Time_to)
	}
	q.Desc("id")
	if orderPar.Page_size == 0 {
		orderPar.Page_size = 20
	}
	q.Limit(orderPar.Page_size, orderPar.Page_no*orderPar.Page_size)

	var rt []Workorders

	err = q.Find(&rt)

	for i := 0; i < len(rt); i++ {
		rt[i].ProductTypeImage, err = s.findOrderPicture(rt[i].ProductCode)
	}

	if err != nil {
		return nil, err
	} else {
		return rt, nil
	}
}

func (s *Service) Steps(workorderID int64) ([]Steps, error) {
	var rt []Steps

	q := s.eng.Alias("s").Where("s.workorder_id = ?", workorderID).Asc("s.id")

	err := q.Find(&rt)

	if err != nil {
		return nil, err
	} else {
		return rt, nil
	}
}

func (s *Service) WorkorderStep(workorderID int64) (*WorkorderStep, error) {
	workorder, err := s.GetWorkOrder(workorderID, true)
	if err != nil {
		return nil, err
	}

	var p interface{}
	err = json.Unmarshal([]byte(workorder.Payload), &p)
	if err == nil {
		workorder.MarshalPayload = p
	}

	steps, err := s.Steps(workorderID)
	if err != nil {
		return nil, err
	}

	for k, v := range steps {
		err := json.Unmarshal([]byte(v.Payload), &steps[k].MarshalPayload)
		if err != nil {
			s.diag.Error(err.Error(), err)
		}
	}

	return &WorkorderStep{
		Workorders: workorder,
		Steps:      steps,
	}, nil
}

func (s *Service) DeleteWorkAndStep(ss *xorm.Session, code string, uniqueNum int64) (bool, error) {
	var workorder Workorders

	bool, e := ss.Alias("r").Where("r.code = ?", code).Get(&workorder)
	if e != nil {
		return false, e
	}
	if !bool {
		return false, nil
	}
	if workorder.UniqueNum > uniqueNum {
		return true, nil
	}
	sql1 := "delete from `workorders` where id = ?"
	sql2 := "delete from `steps` where x_workorder_id = ?"

	_, err := ss.Exec(sql1, workorder.Id)

	_, err = ss.Exec(sql2, workorder.Id)
	return false, err
}

func (s *Service) UpdateStep(step *Steps) (*Steps, error) {

	sql := "update `steps` set status = ? where id = ?"
	_, err := s.eng.Exec(sql,
		step.Status,
		step.Id)

	if err != nil {
		return step, err
	} else {
		return step, nil
	}
}

func (s *Service) UpdateStepData(step *Steps) (*Steps, error) {

	sql := "update `steps` set data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		step.Data,
		step.Id)

	if err != nil {
		return step, err
	} else {
		return step, nil
	}
}

func (s *Service) UpdateOrderData(order *Workorders) (*Workorders, error) {

	sql := "update `workorders` set data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		order.Data,
		order.Id)

	if err != nil {
		return order, err
	} else {
		return order, nil
	}
}

func (s *Service) GetLastIncompleteCurve(toolSN string) (*Curves, error) {
	curve := Curves{}
	e := s.eng.Alias("c").Where("c.tightening_id = ?", "").And("c.tool_sn = ?", toolSN).OrderBy("c.update_time").Desc("c.update_time").Find(&curve)

	if e != nil {
		return &curve, e
	} else {
		return &curve, nil
	}
}

func (s *Service) getLasIncompleteCurveForTool(session *xorm.Session, toolSN string) (*Curves, error) {

	curve := Curves{}
	exist, err := session.Alias("c").Where("c.tightening_id = ?", "").And("c.tool_sn = ?", toolSN).OrderBy("c.update_time").Desc("c.update_time").Get(&curve)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("Curve Not Found")
	}

	return &curve, nil
}

func (s *Service) updateIncompleteCurve(session *xorm.Session, curve *Curves) error {

	sql := "update `curves` set tightening_id = ?, curve_file = ? where id = ?"
	_, err := session.Exec(sql,
		curve.TighteningID,
		curve.CurveFile,
		curve.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateIncompleteCurveAndSaveResult(result *Results) error {

	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 获取最近不完整的曲线
	curve, err := s.getLasIncompleteCurveForTool(session, result.ToolSN)

	if err == nil {
		// 更新曲线
		curve.TighteningID = result.TighteningID
		curve.CurveFile = fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)
		err = s.updateIncompleteCurve(session, curve)
		if err != nil {
			return err
		}

		result.CurveFile = curve.CurveFile
	}

	// 保存结果
	_, err = session.Insert(result)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Failed")
	}

	return nil
}

func (s *Service) getLasIncompleteResultForTool(session *xorm.Session, toolSN string) (*Results, error) {

	result := Results{}
	exist, err := session.Alias("r").Where("r.curve_file = ?", "").And("r.gun_sn = ?", toolSN).OrderBy("r.update_time").Desc("r.update_time").Get(&result)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("Result Not Found")
	}

	return &result, nil
}

func (s *Service) updateIncompleteResult(session *xorm.Session, result *Results) error {

	sql := "update `results` set curve_file = ? where id = ?"
	_, err := session.Exec(sql,
		result.CurveFile,
		result.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateIncompleteResultAndSaveCurve(curve *Curves) error {
	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 获取最近不完整的结果
	result, err := s.getLasIncompleteResultForTool(session, curve.ToolSN)

	if err == nil {
		// 更新结果
		result.CurveFile = fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)
		err = s.updateIncompleteResult(session, result)
		if err != nil {
			return err
		}

		curve.CurveFile = result.CurveFile
		curve.TighteningID = result.TighteningID
	}

	// 保存曲线
	_, err = session.Insert(curve)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Fail")
	}

	return nil
}

func (s *Service) getIncompleteResults(session *xorm.Session, toolSN string) ([]Results, error) {
	results := []Results{}

	err := session.Alias("r").Where("r.curve_file = ?", "").And("r.gun_sn = ?", toolSN).OrderBy("r.update_time").Desc("r.update_time").Find(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Service) getIncompleteCurves(session *xorm.Session, toolSN string) ([]Curves, error) {
	curves := []Curves{}

	err := session.Alias("c").Where("c.tightening_id = ?", "").And("c.tool_sn = ?", toolSN).OrderBy("c.update_time").Desc("c.update_time").Find(&curves)
	if err != nil {
		return nil, err
	}

	return curves, nil
}

func (s *Service) ClearToolResultAndCurve(toolSN string) error {
	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 获取所有不完整的结果
	results, err := s.getIncompleteResults(session, toolSN)
	if err != nil {
		return err
	}

	// 处理不完整的结果
	for _, r := range results {
		r.CurveFile = fmt.Sprintf("%s_%s.json", r.ToolSN, r.TighteningID)
		err := s.updateIncompleteResult(session, &r)
		if err != nil {
			return err
		}
	}

	// 获取所有不完整的曲线
	curves, err := s.getIncompleteCurves(session, toolSN)
	if err != nil {
		return err
	}

	// 处理不完整的曲线
	for _, c := range curves {
		c.TighteningID = utils.GenerateID()
		c.CurveFile = fmt.Sprintf("invalid_%s_%s.json", c.ToolSN, c.TighteningID)
		err := s.updateIncompleteCurve(session, &c)
		if err != nil {
			return err
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Fail")
	}

	return nil
}

func (s *Service) GetWorkorderByCode(code string) (*Workorders, error) {
	order := Workorders{}
	rt, err := s.eng.Alias("w").Where("w.code = ?", code).Limit(1).Get(&order)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, errors.New("Workorder Does Not Exist")
		} else {
			return &order, nil
		}
	}
}

func (s *Service) GetStepByCode(code string) (*Steps, error) {
	step := Steps{}
	rt, err := s.eng.Alias("s").Where("s.code = ?", code).Limit(1).Get(&step)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, errors.New("Step Does Not Exist")
		} else {
			return &step, nil
		}
	}
}

// seq, count
func (s *Service) CalBatch(workorderID int64) (int, int) {
	result, err := s.FindTargetResultForJobManual(workorderID)
	if err != nil {
		return 1, 1
	}

	if result.Result == RESULT_OK {
		return result.GroupSeq + 1, 1
	} else {
		return result.GroupSeq, result.Count + 1
	}
}

func (s *Service) PatchResultFromDB(result *Results, mode string) error {
	dbTool, err := s.GetTool(result.ToolSN)
	if err == nil && dbTool.CurrentWorkorderID != 0 {

		if mode == "job" {
			result.Seq, result.Count = s.CalBatch(dbTool.CurrentWorkorderID)
		} else {
			result.Seq = dbTool.Seq
			result.Count = dbTool.Count
		}

		result.WorkorderID = dbTool.CurrentWorkorderID
		result.UserID = dbTool.UserID
		result.Batch = fmt.Sprintf("%d/%d", result.Seq, dbTool.Total)

		dbStep, err := s.GetStep(dbTool.StepID)
		if err != nil {
			s.diag.Error("Get Step Failed", err)
			return err
		}

		consume, err := s.GetConsumeBySeqInStep(&dbStep, result.Seq)
		if err != nil {
			s.diag.Error("Get Consume Failed", err)
			return err
		}

		result.NutNo = consume.NutNo
		result.ScannerCode = dbTool.ScannerCode
	}

	return nil
}

func (s *Service) GetConsumeBySeqInStep(step *Steps, seq int) (*StepComsume, error) {
	if step == nil {
		return nil, errors.New("Step Is Nil")
	}

	ts := TighteningStep{}
	if err := json.Unmarshal([]byte(step.Step), &ts); err != nil {
		return nil, err
	}

	if len(ts.TighteningPoints) == 0 {
		return nil, errors.New("Consumes Is Empty")
	}

	for k, v := range ts.TighteningPoints {
		if v.Seq == seq {
			return &ts.TighteningPoints[k], nil
		}
	}

	return nil, errors.New("Consume Not Found")
}

func (s *Service) UpdateToolLocation(toolSN string, location string) error {
	tool, err := s.GetTool(toolSN)
	if err != nil {
		return err
	}

	tool.Location = location
	_, err = s.eng.Id(tool.Id).Update(tool)
	return err
}

func (s *Service) GetToolLocation(toolSN string) (string, error) {
	tool, err := s.GetTool(toolSN)
	if err != nil {
		return "", err
	}

	return tool.Location, nil
}
