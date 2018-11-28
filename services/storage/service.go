package storage

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"sync/atomic"
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
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Url,
		c.DBName)
	engine, err := xorm.NewEngine("postgres", info)

	if err != nil {
		return errors.Wrapf(err, "Create postgres engine fail")
	}

	exist, err := engine.IsTableExist("Workorders")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Workorders)); err != nil {
				return errors.Wrapf(err, "Create Table Workorders fail")
			}
		}
	}

	exist, err = engine.IsTableExist("Results")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Results)); err != nil {
				return errors.Wrapf(err, "Create Table Results fail")
			}
		}
	}

	exist, err = engine.IsTableExist("Curves")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Curves)); err != nil {
				return errors.Wrapf(err, "Create Table Curves fail")
			}
		}
	}

	exist, err = engine.IsTableExist("Controllers")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Controllers)); err != nil {
				return errors.Wrapf(err, "Create Table Controllers fail")
			}
		}
	}

	exist, err = engine.IsTableExist("Guns")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Guns)); err != nil {
				return errors.Wrapf(err, "Create Table Guns fail")
			}
		}
	}

	exist, err = engine.IsTableExist("RoutingOperations")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(RoutingOperations)); err != nil {
				return errors.Wrapf(err, "Create Table RoutingOperations fail")
			}
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

	s.eng.Close()

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

func (s *Service) FindResultsByWorkorder(workorder_id int64) ([]Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", workorder_id).OrderBy("r.seq")

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
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

func (s *Service) ListUnuploadResults() ([]Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.has_upload = ?", false).And("r.stage = ?", "final")

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (s *Service) CurveExist(curve *Curves) (bool, error) {

	has, err := s.eng.Exist(&Curves{ResultID: curve.ResultID, Count: curve.Count})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (s *Service) GetCurve(curve *Curves) (interface{}, error) {

	var rt_curve Curves

	rt, err := s.eng.Alias("c").Where("c.result_id = ?", curve.ResultID).And("c.count = ?", curve.Count).Get(&rt_curve)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, nil
		} else {
			return rt_curve, nil
		}
	}
}

func (s *Service) GetGun(serial string) (Guns, error) {

	var rt_gun Guns

	rt, err := s.eng.Alias("g").Where("g.serial = ?", serial).Get(&rt_gun)

	if err != nil {
		return rt_gun, err
	} else {
		if !rt {
			return rt_gun, errors.New("found gun failed")
		} else {
			return rt_gun, nil
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

	e := s.eng.Alias("c").Where("c.has_upload = ?", false).Find(&curves)
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

func (s *Service) GetResult(resultId int64, count int) (Results, error) {
	var err error

	result := Results{}

	var rt bool
	if count == 0 {
		rt, err = s.eng.Alias("r").Where("r.x_result_id = ?", resultId).Limit(1).Get(&result)
	} else {
		rt, err = s.eng.Alias("r").Where("r.x_result_id = ?", resultId).And("r.count = ?", count).Limit(1).Get(&result)
	}

	if err != nil {
		return result, err
	} else {
		if !rt {
			return result, errors.New("result does not exist")
		} else {
			return result, nil
		}
	}
}

func (s *Service) ListWorkorders(hmi_sn string, status string) ([]Workorders, error) {
	var workorders []Workorders

	sql := fmt.Sprintf("select * from workorders where hmi_sn = '%s'", hmi_sn)
	if status != "" {
		sql = sql + fmt.Sprintf(" and workorders.status = '%s'", status)
	}

	sql = sql + " order by lnr asc"

	err := s.eng.SQL(sql).Find(&workorders)

	return workorders, err
}

func (s *Service) GetWorkorder(id int64, raw bool) (Workorders, error) {

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

func (s *Service) FindWorkorder(hmi_sn string, workcenter_code string, code string) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.hmi_sn = ? or w.workcenter_code = ?", hmi_sn, workcenter_code).And("w.long_pin = ? or w.vin = ? or w.knr = ?", code, code, code).Get(&workorder)

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

func (s *Service) FindNextWorkorder(hmi_sn string) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.hmi_sn = ?", hmi_sn).And("w.status = ?", "ready").Asc("w.update_time").Get(&workorder)

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

func (s *Service) UpdateResultUpload(upload bool, r_id int64) (int64, error) {
	sql := "update `results` set has_upload = ? where x_result_id = ?"

	r, err := s.eng.Exec(sql, upload, r_id)

	id, _ := r.RowsAffected()
	if err != nil {

		return id, errors.Wrapf(err, "Update result upload status fail for id : %d", id)
	} else {
		return id, nil
	}
}

func (s *Service) UpdateResult(result *Results) (int64, error) {

	sql := "update `results` set controller_sn = ?, result = ?, has_upload = ?, stage = ?, update_time = ?, pset_define = ?, result_value = ?, count = ?, batch = ?, gun_sn = ?, spent = ? where id = ?"
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
		result.GunSN,
		result.Spent,
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

	return nil
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

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", raw_workorder_id).And("r.stage = ? or r.seq = ?", RESULT_STAGE_INIT, 0).OrderBy("r.seq")

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
	sql := "update `routing_operations` set job = ?, max_op_time = ?, name = ?, img = ?, product_id = ?, product_type = ?, workcenter_code = ?, vehicle_type_img = ?, points = ?, workcenter_id = ? where operation_id = ?"
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
		ro.OperationID)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) GetRoutingOperations(op_id int64) (RoutingOperations, error) {

	var ro RoutingOperations

	rt, err := s.eng.Alias("r").Where("r.operation_id = ?", op_id).Get(&ro)

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

	sql := "select * from results, workorders where results.x_workorder_id = workorders.id and results.stage = 'final' "
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

func (s *Service) UpdateResultTriggerTime(trigger_type string, trigger_time time.Time, controller_sn string) error {

	sql := fmt.Sprintf("update `controllers` set %s = ? where controller_sn = ?", trigger_type, trigger_time)
	_, err := s.eng.Exec(sql, controller_sn)

	if err != nil {
		return err
	} else {
		return nil
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
