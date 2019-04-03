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
	//eng         *xorm.Engine
	eng *gorm.DB
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
	sConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=10",
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

func (s *Service) BatchSave(results []*ResultObject) error {
	for _, v := range results {
		cur, _ := json.Marshal(v.OR["cur_objects"])

		oneTimePass := false
		if reflect.ValueOf(v.OR["one_time_pass"]).String() == "pass" {
			oneTimePass = true
		}

		sql := fmt.Sprintf("SELECT create_operation_result(%v, %v, timestamp '%s', %v, %v, %v, '%s', '%s', %v, '%s', %v, %v, '%s', %v, %v, %v, %v, %v, %v, '%s', '%s', '%s', %v, '%s', %v, '%s', %v, '%s', '%s', '%s')",
			v.OR["pset_m_threshold"],
			v.OR["pset_m_max"],
			reflect.ValueOf(v.OR["control_date"]).String(),
			v.OR["pset_w_max"],
			v.OR["user_id"],
			oneTimePass,
			reflect.ValueOf(v.OR["pset_strategy"]).String(),
			reflect.ValueOf(v.OR["measure_result"]).String(),
			v.OR["pset_w_threshold"],
			string(cur),
			v.OR["pset_m_target"],
			v.OR["pset_m_min"],
			reflect.ValueOf(v.OR["final_pass"]).String(),
			v.OR["measure_degree"],
			v.OR["measure_t_don"],
			v.OR["measure_torque"],
			v.OR["op_time"],
			v.OR["pset_w_min"],
			v.OR["pset_w_target"],
			reflect.ValueOf(v.OR["lacking"]).String(),
			reflect.ValueOf(v.OR["quality_state"]).String(),
			reflect.ValueOf(v.OR["exception_reason"]).String(),
			v.Send,
			reflect.ValueOf(v.OR["batch"]).String(),
			v.OR["workorder_id"],
			v.OR["nut_no"],
			v.OR["tightening_id"],
			reflect.ValueOf(v.OR["vin"]).String(),
			reflect.ValueOf(v.OR["model"]).String(),
			reflect.ValueOf(v.OR["tool_sn"]).String())

		err := s.eng.Exec(sql)
		if err.Error != nil {
			if strings.HasSuffix(err.Error.Error(), "tid_vin_gun_uniq\"") || strings.HasSuffix(err.Error.Error(), "tid_wo_gun_uniq\"") {
				return nil
			}
			return err.Error
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

func (s *Service) UpdateResults(result *OperationResult, id int64, sent int) error {

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
