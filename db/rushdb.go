package rushdb

import (
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/kataras/iris/core/errors"

	"github.com/masami10/rush/payload"
	"time"
)

type DB struct {
	DBName string	`yaml:"dbname"`
	URL string		`yaml:"url"`
	Port int		`yaml:"port"`
	User string		`yaml:"user"`
	Pwd string		`yaml:"pwd"`
}

func (db *DB)Init() (error) {

	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return err
	}

	var exist bool
	var e error

	exist, e = engine.IsTableExist(&Workorders{})
	if e != nil {
		return e
	} else {
		if !exist {
			e = engine.CreateTables(&Workorders{})
			if e != nil {
				return e
			}
		}
	}

	exist, e = engine.IsTableExist(&Results{})
	if e != nil {
		return e
	} else {
		if !exist {
			e = engine.CreateTables(&Results{})
			if e != nil {
				return e
			}
		}
	}

	exist, e = engine.IsTableExist(&Curves{})
	if e != nil {
		return e
	} else {
		if !exist {
			e = engine.CreateTables(&Curves{})
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func (db *DB) UpdateResults(id int, count int, flag bool)(error) {
	//results := []Results{}
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return err
	}

	if count > 0 {
		sql := "update `results` set result_upload = ? where result_id = ? and count = ?"
		_, err = engine.Exec(sql, flag, id, count)
	} else {
		sql := "update `results` set result_upload = ? where result_id = ?"
		_, err = engine.Exec(sql, flag, id)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (db *DB) FindResults(result_upload bool, result []string)([]Results, error) {
	results := []Results{}
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return results, err
	}

	e := engine.Alias("r").Where("r.result_upload = ?", result_upload).And("r.need_upload = ?", true).And("r.result <> ?", "NONE").Find(&results)
	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (db *DB) ListResults(id int)([]Results, error) {
	results := []Results{}

	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return results, err
	}

	e := engine.Alias("r").Where("r.result_id = ?", id).Find(&results)
	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (db *DB) InsertResults(result Results) (error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return err
	}

	_, e := engine.Insert(result)
	if e != nil {
		return e
	} else {
		fmt.Printf("new result:%d\n", result.ResultId)
	}

	return nil
}

func (db *DB) CurveExist(curve Curves) (bool, error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return false, err
	}

	has, err := engine.Exist(&Curves{ ResultID: curve.ResultID, Count: curve.Count})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (db *DB) InsertCurve(curve Curves) (error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return err
	}

	_, e := engine.Insert(curve)
	if e != nil {
		return e
	} else {
		return nil
	}
}

func (db *DB) UpdateCurve(curve Curves) (Curves, error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return curve, err
	}

	sql := "update `curves` set has_upload = ?, curve_file = ?, curve_data = ? where result_id = ? and count = ?"
	_, err = engine.Exec(sql,
		curve.HasUpload, curve.CurveFile, curve.CurveData, curve.ResultID, curve.Count)

	if err != nil {
		return curve, err
	} else {
		return curve, nil
	}
}

func (db *DB) ListCurves(result_id int) ([]Curves, error) {
	var curves []Curves = []Curves{}
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return curves, err
	}

	e := engine.Alias("c").Where("c.result_id = ?", result_id).Find(&curves)
	if e != nil {
		return curves, e
	} else {
		return curves, nil
	}
}


func (db *DB) InsertWorkorders(workorders []payload.ODOOWorkorder) ([]payload.ODOOWorkorder, error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return workorders, err
	}

	new_orders := []payload.ODOOWorkorder{}

	for _, v := range workorders {
		has, _ := db.WorkorderExists(v.ID)
		if has {
			// 忽略存在的记录
			continue
		}

		if len(v.Result_IDs) == 0 {
			//fmt.Printf("has no results\n")
			continue
		}

		o := new(Workorders)
		o.Status = v.Status
		o.WorkorderID = v.ID
		o.PSet, _ = strconv.Atoi(v.PSet)
		o.HMISN = v.HMI.UUID
		o.Knr = v.KNR
		o.NutTotal = v.NutTotal
		o.Vin = v.VIN
		o.MaxOpTime = v.Max_op_time
		o.MaxRedoTimes = v.Max_redo_times

		ids, _ := json.Marshal(v.Result_IDs)
		o.ResultIDs = string(ids)

		_, err := engine.Insert(o)
		if err != nil {
			return workorders, err
		} else {
			new_orders = append(new_orders, v)
			fmt.Printf("new workorder:%d\n", o.WorkorderID)
		}

		// 预保存结果
		for _, result_id := range v.Result_IDs {
			r := new(Results)
			r.ResultId = result_id
			r.ControllerSN = ""
			r.WorkorderID = o.WorkorderID
			r.Result = payload.RESULT_NONE
			r.HasUpload = false
			r.Stage = "init"
			r.UpdateTime = time.Now()
			r.PSetDefine = ""
			r.ResultValue = ""
			r.Count = 1

			_, err := engine.Insert(r)
			if err != nil {
				return workorders, err
			} else {
				fmt.Printf("new result:%d\n", result_id)
			}
		}

	}

	return new_orders, nil

}

func (db *DB) WorkorderExists(id int) (bool, error) {
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return false, err
	}

	has, err := engine.Exist(&Workorders{ WorkorderID: id})
	if err != nil {
		return false, err
	} else {
		return has, nil
	}
}

func (db *DB) GetResult(id int, count int) (Results, error) {
	var err error
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	result := Results{}

	if err != nil {
		return result, err
	}

	var rt bool
	if count == 0 {
		rt, err = engine.Alias("r").Where("r.result_id = ?", id).Limit(1).Get(&result)
	} else {
		rt, err = engine.Alias("r").Where("r.result_id = ?", id).And("r.count = ?", count).Limit(1).Get(&result)
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

func (db *DB) GetWorkorder(id int) (Workorders, error) {
	var err error
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	workorder := Workorders{}

	if err != nil {
		return workorder, err
	}

	var rt bool
	rt, err = engine.Alias("w").Where("w.workorder_id = ?", id).Get(&workorder)

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

func (db *DB) FindWorkorder(hmi_sn string, vin string, knr string) (Workorders, error) {
	var err error
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	workorder := Workorders{}

	if err != nil {
		return workorder, err
	}

	var rt bool
	if vin != "" {
		rt, err = engine.Alias("w").Where("w.hmi_sn = ?", hmi_sn).And("w.vin = ?", vin).Get(&workorder)
	} else {
		rt, err = engine.Alias("w").Where("w.hmi_sn = ?", hmi_sn).And("w.knr = ?", knr).Get(&workorder)
	}

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

func (db *DB) ListNeedPushResults() ([]Results, error) {
	results := []Results{}

	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return results, err
	}

	e := engine.Alias("r").Where("r.need_upload = ?", true).And("r.result_upload = ?", false).Find(&results)
	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (db *DB) UpdateResult(result Results) (Results, error) {
	var err error
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pwd,
		db.URL,
		db.Port,
		db.DBName))

	if err != nil {
		return result, err
	}

	sql := "update `results` set ControllerSN = ?, Result = ?, HasUpload = ?, Stage = ?, UpdateTime = ?, PSetDefine = ?, ResultValue = ?, Count = ? where result_id = ?"
	_, err = engine.Exec(sql,
		result.ControllerSN,
		result.Result,
		result.HasUpload,
		result.Stage,
		result.UpdateTime,
		result.PSetDefine,
		result.ResultValue,
		result.Count,
		result.ResultId)

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}