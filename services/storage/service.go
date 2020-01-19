package storage

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	OpenEngineSuccess(info string)
	UpdateResultSuccess(id int64)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	open        bool
	eng         *gorm.DB
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
		open: false,
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) openStorageEngine() error {
	c := s.Config()
	urls := strings.Split(c.Url, ":")
	sConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=10",
		urls[0],
		urls[1],
		c.User,
		c.DBName,
		c.Password)

	_db, err := gorm.Open("postgres", sConn)
	if err != nil {
		return err
	}

	_db.AutoMigrate(&OperationResultModel{})

	s.eng = _db
	s.open = true

	s.diag.OpenEngineSuccess(sConn)

	return nil
}

func (s *Service) Open() error {
	return s.openStorageEngine()
}

func (s *Service) Close() error {
	eng := s.eng
	if eng == nil {
		return nil
	}
	if err := eng.Close(); err != nil {
		s.diag.Error("Storage Service Close Engine Error", err)
	}

	return nil
}

func (s *Service) BatchSave(results []*ResultObject) error {
	eng, err := s.ensureStorageEng()

	if eng == nil || err != nil {
		s.diag.Error("Cannot Update Results", err)
		return err
	}

	for _, v := range results {
		cur, _ := json.Marshal(v.OR["cur_objects"])

		//oneTimePass := false
		//if reflect.ValueOf(v.OR["one_time_pass"]).String() == "pass" {
		//	oneTimePass = true
		//}

		sql := fmt.Sprintf("SELECT create_operation_result(%v, %v, timestamp '%s', %v, %v, '%s', '%s', %v, '%s', %v, %v, %v, %v, %v, %v, %v, %v, '%s', '%s', %v, '%s', '%s', '%s', %v, '%s', '%s', '%s')",
			v.OR["pset_m_threshold"],
			v.OR["pset_m_max"],
			reflect.ValueOf(v.OR["control_date"]).String(),
			v.OR["pset_w_max"],
			v.OR["user_id"],
			//oneTimePass,
			reflect.ValueOf(v.OR["pset_strategy"]).String(),
			reflect.ValueOf(v.OR["measure_result"]).String(),
			v.OR["pset_w_threshold"],
			string(cur),
			v.OR["pset_m_target"],
			v.OR["pset_m_min"],
			//reflect.ValueOf(v.OR["final_pass"]).String(),
			v.OR["measure_degree"],
			v.OR["measure_t_don"],
			v.OR["measure_torque"],
			v.OR["op_time"],
			v.OR["pset_w_min"],
			v.OR["pset_w_target"],
			//reflect.ValueOf(v.OR["lacking"]).String(),
			reflect.ValueOf(v.OR["quality_state"]).String(),
			reflect.ValueOf(v.OR["exception_reason"]).String(),
			v.Send,
			reflect.ValueOf(v.OR["batch"]).String(),
			reflect.ValueOf(v.OR["workorder_name"]).String(),
			v.OR["nut_no"],
			v.OR["tightening_id"],
			reflect.ValueOf(v.OR["vin"]).String(),
			reflect.ValueOf(v.OR["model"]).String(),
			reflect.ValueOf(v.OR["tool_sn"]).String())

		e := eng.Exec(sql)
		if e == nil {
			return errors.New("BatchSave Result Error")
		}
		for _, err := range e.GetErrors() {
			if err != nil {
				fmt.Println(err.Error())
				if strings.Contains(err.Error(), "tid_tool") {
					continue
				}
				return err
			}
		}

	}

	return nil
}

//func (s *Service) BatchSave(results []*ResultObject) error {
//
//	var updateResults []*ResultObject
//	var insertResults []*ResultObject
//
//	for _, v := range results {
//		if v.ID == 0 {
//			insertResults = append(insertResults, v)
//		} else {
//			updateResults = append(updateResults, v)
//		}
//	}
//
//	// 更新
//	if len(updateResults) > 0 {
//
//		var arrKeys []string
//		for _, v := range KEYS {
//			if v == "workcenter_id" || v == "product_id" || v == "gun_id" || v == "time" || v == "consu_product_id" {
//				continue
//			}
//
//			arrKeys = append(arrKeys, v)
//		}
//
//		var arrSKeys []string
//		for _, v := range arrKeys {
//			arrSKeys = append(arrSKeys, fmt.Sprintf("s.%s", v))
//		}
//		strKeys := fmt.Sprintf("(%s)", strings.Join(arrKeys, ","))
//		strSKeys := fmt.Sprintf("(%s)", strings.Join(arrSKeys, ","))
//
//		arrValues := []string{}
//		for _, v := range updateResults {
//			if v == nil {
//				continue
//			}
//
//			arrValue := []string{}
//			for _, k := range arrKeys {
//				if k == "id" {
//					arrValue = append(arrValue, fmt.Sprintf("%d", v.ID))
//					continue
//				}
//
//				//if k == "lacking" {
//				//	arrValue = append(arrValue, "'normal'")
//				//	continue
//				//}
//
//				if k == "time" {
//					//arrValue = append(arrValue, "'normal'")
//					continue
//				}
//
//				if k == "sent" {
//					if v.Send == 1 {
//						arrValue = append(arrValue, "true")
//					} else {
//						arrValue = append(arrValue, "false")
//					}
//					continue
//				}
//
//				if k == "cur_objects" {
//					curs, _ := json.Marshal(v.OR[k])
//					//fmt.Printf("%s\n", string(curs))
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", string(curs)))
//					continue
//				}
//
//				if k == "control_date" {
//					arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
//					continue
//				}
//
//				targetType := fmt.Sprintf("%T", v.OR[k])
//
//				switch targetType {
//				case "<nil>":
//					arrValue = append(arrValue, "''")
//				case "array":
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
//				case "time":
//					arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
//				case "string":
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
//				default:
//					arrValue = append(arrValue, fmt.Sprintf("%v", v.OR[k]))
//				}
//			}
//
//			arrValues = append(arrValues, fmt.Sprintf("(%s)", strings.Join(arrValue, ",")))
//		}
//		strValues := strings.Join(arrValues, ",")
//
//		sql := fmt.Sprintf("UPDATE operation_result AS o SET %s = %s FROM(VALUES %s ) AS s %s WHERE o.id = s.id", strKeys, strSKeys, strValues, strKeys)
//
//		result := s.eng.Exec(sql)
//		if result.Error != nil {
//			return result.Error
//		}
//	}
//
//	// 新增
//	if len(insertResults) > 0 {
//
//		arrSKeys := []string{}
//		for _, v := range KEYS {
//			arrSKeys = append(arrSKeys, v)
//		}
//
//		arrInsertValues := []string{}
//		for _, v := range insertResults {
//			if v == nil {
//				continue
//			}
//
//			arrValue := []string{}
//			for _, k := range KEYS {
//				if k == "id" {
//					arrValue = append(arrValue, "nextval('operation_result_id_seq')")
//					continue
//				}
//
//				//if k == "lacking" {
//				//	arrValue = append(arrValue, "'normal'")
//				//	continue
//				//}
//
//				if k == "gun_id" {
//					//fmt.Printf("%s:%d\n", k, v.OR[k])
//					if reflect.ValueOf(v.OR[k]).Float() == 0 {
//						for n, m := range arrSKeys {
//							if m == k {
//								kk := n + 1
//								arrSKeys = append(arrSKeys[:n], arrSKeys[kk:]...)
//							}
//						}
//						continue
//					}
//				}
//
//				if k == "sent" {
//					if v.Send == 1 {
//						arrValue = append(arrValue, "true")
//					} else {
//						arrValue = append(arrValue, "false")
//					}
//					continue
//				}
//
//				if k == "cur_objects" {
//					curs, _ := json.Marshal(v.OR[k])
//					//fmt.Printf("%s\n", string(curs))
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", string(curs)))
//					continue
//				}
//
//				if k == "control_date" {
//					arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
//					continue
//				}
//
//				if k == "time" {
//					arrValue = append(arrValue, fmt.Sprintf("date '%s'", reflect.ValueOf(v.OR["control_date"]).String()))
//					continue
//				}
//
//				targetType := fmt.Sprintf("%T", v.OR[k])
//
//				switch targetType {
//				case "<nil>":
//					arrValue = append(arrValue, "''")
//				case "array":
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
//				case "time":
//					arrValue = append(arrValue, fmt.Sprintf("timestamp '%s'", reflect.ValueOf(v.OR[k]).String()))
//				case "string":
//					arrValue = append(arrValue, fmt.Sprintf("'%s'", reflect.ValueOf(v.OR[k]).String()))
//				default:
//					arrValue = append(arrValue, fmt.Sprintf("%v", v.OR[k]))
//				}
//			}
//
//			arrInsertValues = append(arrInsertValues, fmt.Sprintf("(%s)", strings.Join(arrValue, ",")))
//		}
//
//		strInsertValues := strings.Join(arrInsertValues, ",")
//
//		strKeys := fmt.Sprintf("(%s)", strings.Join(arrSKeys, ","))
//		insertSql := fmt.Sprintf("INSERT INTO \"operation_result\" %s VALUES %s RETURNING id", strKeys, strInsertValues)
//
//		insertResult := s.eng.Exec(insertSql)
//		if insertResult.Error != nil {
//			return insertResult.Error
//		}
//	}
//
//	return nil
//}

func (s *Service) ensureStorageEng() (*gorm.DB, error) {
	if s.eng != nil {
		return s.eng, nil
	}

	if err := s.openStorageEngine(); err != nil {
		s.diag.Error("openStorageEngine Error:", err)
		return nil, err
	}

	return s.eng, nil
}

func (s *Service) UpdateResults(result *OperationResult, id int64, sent int) error {

	eng, err := s.ensureStorageEng()

	if eng == nil && err != nil {
		s.diag.Error("Cannot Update Results", err)
		return err
	}

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

		result := eng.Create(&r)

		if result.Error != nil {
			return errors.Wrapf(result.Error, "insert result record fail: %s", result.Error.Error())
		}

		s.diag.UpdateResultSuccess(result.RowsAffected)
	} else {
		// 更新

		r.Id = id
		result := eng.Save(&r)
		fmt.Printf("%d\n", id)

		if result.Error != nil {
			return errors.Wrapf(result.Error, "Update result record %d fail: %s", id, result.Error.Error())
		}

		s.diag.UpdateResultSuccess(result.RowsAffected)
	}

	return nil
}
