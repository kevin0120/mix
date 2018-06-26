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
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Workorders")
	}
	if !exist {
		if err := engine.Sync2(new(Workorders)); err != nil {
			return errors.Wrapf(err, "Create Table Workorders fail")
		}

	}

	exist, err = engine.IsTableExist("Results")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Results")
	}
	if !exist {
		if err := engine.Sync2(new(Results)); err != nil {
			return errors.Wrapf(err, "Create Table Results fail")
		}

	}

	exist, err = engine.IsTableExist("Curves")
	if err != nil {
		return errors.Wrapf(err, "Check Table exist %s fail", "Curves")
	}
	if !exist {
		if err := engine.Sync2(new(Curves)); err != nil {
			return errors.Wrapf(err, "Create Table Curves fail")
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

func (s *Service) CurveExist(curve *Curves) (bool, error) {

	has, err := s.eng.Exist(&Curves{ResultID: curve.ResultID, Count: curve.Count})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (s *Service) UpdateCurve(curve *Curves) (*Curves, error) {

	sql := "update `curves` set has_upload = ?, curve_file = ?, curve_data = ? where result_id = ? and count = ?"
	_, err := s.eng.Exec(sql,
		curve.HasUpload, curve.CurveFile, curve.CurveData, curve.ResultID, curve.Count)

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

func (s *Service) InsertWorkorder(workorder *Workorders, results *[]Results) error {

	has, err := s.WorkorderExists(workorder.WorkorderID)
	if has {
		// 忽略存在的工单
		return nil
	}

	session := s.eng.NewSession()
	defer session.Close()

	err = session.Begin()
	// 执行事务

	// 保存新工单
	_, err = session.Insert(workorder)
	if err != nil {
		session.Rollback()
		return errors.Wrapf(err, "store data fail")
	} else {
		s.diag.Debug(fmt.Sprintf("new workorder:%d", workorder.WorkorderID))
	}

	// 预保存结果
	for _, v := range *results {

		_, err = session.Insert(v)
		if err != nil {
			session.Rollback()
			return errors.Wrapf(err, "store data fail")
		} else {
			s.diag.Debug(fmt.Sprintf("new result:%d", v.ResultId))
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil
}

func (s *Service) WorkorderExists(id int64) (bool, error) {

	has, err := s.eng.Exist(&Workorders{WorkorderID: id})
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

func (s *Service) GetWorkorder(id int64) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.x_workorder_id = ?", id).Get(&workorder)

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

func (s *Service) FindWorkorder(hmi_sn string, code string) (Workorders, error) {

	var workorder Workorders

	rt, err := s.eng.Alias("w").Where("w.hmi_sn = ?", hmi_sn).And("w.long_pin = ?", code).Or("w.vin = ?", code).Or("w.knr = ?", code).Get(&workorder)

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

	sql := "update `results` set controller_sn = ?, result = ?, has_upload = ?, stage = ?, update_time = ?, pset_define = ?, result_value = ?, count = ? where x_result_id = ?"
	r, err := s.eng.Exec(sql,
		result.ControllerSN,
		result.Result,
		result.HasUpload,
		result.Stage,
		result.UpdateTime,
		result.PSetDefine,
		result.ResultValue,
		result.Count,
		result.ResultId)

	id, _ := r.RowsAffected()
	if err != nil {

		return id, errors.Wrapf(err, "Update result fail for id : %d", id)
	} else {
		return id, nil
	}
}

func (s *Service) UpdateWorkorder(workorder *Workorders) (*Workorders, error) {

	sql := "update `workorders` set status = ? where x_workorder_id = ?"
	_, err := s.eng.Exec(sql,
		workorder.Status,
		workorder.WorkorderID)

	if err != nil {
		return workorder, err
	} else {
		return workorder, nil
	}
}

func (s *Service) UpdateResultByCount(id int64, count int, flag bool) error {

	var err error
	if count > 0 {
		sql := "update `results` set has_upload = ? where x_result_id = ? and count = ?"
		_, err = s.eng.Exec(sql, flag, id, count)
	} else {
		sql := "update `results` set has_upload = ? where x_result_id = ?"
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

	sql := fmt.Sprintf("delete from `workorders` where status = 'finished' and update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
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
