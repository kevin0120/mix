package storage

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/masami10/aiis/services/rush"
	"github.com/pkg/errors"
	"strings"
	"sync/atomic"
	"time"
	"reflect"
)

type Diagnostic interface {
	Error(msg string, err error)
	OpenEngineSuccess(info string)
	UpdateResultSuccess(id int64)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	//eng         *xorm.Engine
	eng *gorm.DB

	results  chan *rush.ResultObject
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
		results: make(chan *rush.ResultObject, c.BatchSaveRowsLimit),
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	//c := s.Config()
	//info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
	//	c.User,
	//	c.Password,
	//	c.Url,
	//	c.DBName)
	//engine, err := xorm.NewEngine("postgres", info)
	//
	//if err != nil {
	//	return errors.Wrapf(err, "Create postgres engine fail")
	//}
	//
	//exist, err := engine.IsTableExist("operation_result")
	//if err != nil {
	//	return errors.Wrapf(err, "Check Table exist %s fail", "operation_result")
	//}
	//if !exist {
	//	return errors.New("Check Table exist operation_result fail, Please start Odoo first")
	//}
	//
	//engine.SetMaxOpenConns(c.MaxConnects) // always success
	//
	//s.eng = engine

	urls := strings.Split(s.Config().Url, ":")
	sConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		urls[0],
		urls[1],
		s.Config().User,
		s.Config().DBName,
		s.Config().Password)

	_db, err := gorm.Open("postgres", sConn)
	if err != nil {
		return err
	}

	_db.AutoMigrate(&OperationResultModel{})

	s.eng = _db

	go s.TaskResultsBatchSave()

	s.diag.OpenEngineSuccess(sConn)

	return nil
}

func (s *Service) Close() error {
	if s.eng == nil {
		return nil
	}
	s.eng.Close()

	return nil
}

func (s *Service) TaskResultsBatchSave() {
	idx := 0
	c := s.configValue.Load().(Config)
	results := make([]*rush.ResultObject, c.BatchSaveRowsLimit)

	for {
		select {
		case <- time.After(time.Duration(c.BatchSaveTimeLimit)):
			if idx > 0 {
				s.BatchSave(results)
				idx = 0
			}


		case data := <- s.results:
			results[idx] = data
			idx++

			if idx == s.Config().BatchSaveRowsLimit {
				s.BatchSave(results)
				idx = 0
			}
		}
	}

}

func (s *Service) BatchSave(results []*rush.ResultObject) error {

	arrSKeys := []string{}
	for _, v := range KEYS {
		arrSKeys = append(arrSKeys, fmt.Sprintf("s.%s", v))
	}
	strKeys := fmt.Sprintf("(%s)", strings.Join(KEYS, ","))
	strSKeys := fmt.Sprintf("(%s)", strings.Join(arrSKeys, ","))

	arrValues := []string{}
	for _, v := range results {

		arrValue := []string{}
		for _, k := range KEYS {
			if k == "id" {
				arrValue = append(arrValue, fmt.Sprintf("%d", v.ID))
				continue
			}

			if k == "sent" {
				if v.Send == 1 {
					arrValue = append(arrValue, "true")
				} else {
					arrValue = append(arrValue, "false")
				}
				continue
			}

			if k == "control_date" {
				arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
				continue
			}

			targetType := fmt.Sprintf("%T", v.OR[k])

			switch targetType {
			case "<nil>":
				arrValue = append(arrValue, "''")
			case "array":
				arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
			case "time":
				arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
			case "string":
				arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
			default:
				arrValue = append(arrValue, fmt.Sprintf("%v", v.OR[k]))
			}
		}

		arrValues = append(arrValues, fmt.Sprintf("(%s)", strings.Join(arrValue, ",")))
	}

	strValues := strings.Join(arrValues, ",")

	sql := fmt.Sprintf("UPDATE operation_result AS o SET %s = %s FROM(VALUES %s ) AS s %s WHERE o.id = s.id", strKeys, strSKeys, strValues, strKeys)

	result := s.eng.Exec(sql)
	if result.Error == nil {
		//result.
		fmt.Printf("success:%d\n", result.RowsAffected)
	}
	return result.Error
}

func (s *Service) AddResult(r *rush.ResultObject) {
	s.results <- r
}

func (s *Service) UpdateResults(result *rush.OperationResult, id int64, sent int) error {

	var r OperationResultModel
	r.PsetMThreshold = result.PsetMThreshold
	r.PsetMMax = result.PsetMMax
	r.ControlDate = result.ControlDate
	r.PsetWMax = result.PsetMMax
	r.UserId = result.UserId
	r.OneTimePass = result.OneTimePass
	r.PsetStrategy = result.PsetStrategy
	r.PsetWThreshold = result.PsetWThreshold

	cur, _ := json.Marshal(result.CurObjects)
	r.CurObjects = string(cur)

	r.PsetMTarget = result.PsetMTarget
	r.PsetMMin = result.PsetMMin
	r.FinalPass = result.FinalPass
	r.MeasureDegree = result.MeasureDegree
	r.MeasureTDon = result.MeasureTDone
	r.MeasureTorque = result.MeasureTorque
	r.MeasureResult = result.MeasureResult
	r.OpTime = result.OPTime
	r.PsetWMin = result.PsetWMin
	r.PsetWTarget = result.PsetWTarget
	r.UserId = result.UserId
	r.Lacking = "normal"
	r.QualityState = result.QualityState
	r.ExceptionReason = result.ExceptionReason
	r.WorkcenterId = result.WorkcenterID
	r.Sent = sent
	r.Batch = result.Batch

	if id == 0 {
		// 新增

		r.CreateTime = r.ControlDate
		r.ProductId = result.ProductID
		r.Vin = result.Vin
		r.GunID = result.GunID

		result := s.eng.Create(&r)

		if result.Error != nil {
			return errors.Wrapf(result.Error, "insert result record fail: %s", result.Error.Error())
		}

		s.diag.UpdateResultSuccess(result.RowsAffected)
	} else {
		// 更新

		r.Id = id
		result := s.eng.Save(&r)
		fmt.Printf("%d\n", id)

		if result.Error != nil {
			return errors.Wrapf(result.Error, "Update result record %d fail: %s", id, result.Error.Error())
		}

		s.diag.UpdateResultSuccess(result.RowsAffected)
	}

	return nil
}
