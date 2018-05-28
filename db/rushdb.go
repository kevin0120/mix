package rushdb

import (
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/kataras/iris/core/errors"

	"github.com/masami10/rush/payload"
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
		fmt.Printf("new result:%d\n", result.Result_id)
	}

	return nil
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
		o.Workorder_id = v.ID
		o.PSet, _ = strconv.Atoi(v.PSet)
		o.HMI_sn = v.HMI.UUID
		o.Knr = v.KNR
		o.Nut_total = v.NutTotal
		o.Vin = v.VIN
		o.Max_op_time = v.Max_op_time
		o.Max_redo_times = v.Max_redo_times

		ids, _ := json.Marshal(v.Result_IDs)
		o.Result_ids = string(ids)

		_, err := engine.Insert(o)
		if err != nil {
			return workorders, err
		} else {
			new_orders = append(new_orders, v)
			fmt.Printf("new workorder:%d\n", o.Workorder_id)
		}

		//// 预保存结果
		for _, result_id := range v.Result_IDs {
			r := new(Results)
			r.Workorder_id = o.Workorder_id
			r.Result = payload.RESULT_NONE
			r.Controller_sn = ""
			r.Count = 1
			r.Cur_data = ""
			r.Cur_upload = false
			r.Result_id = result_id
			r.Result_data = ""
			r.Result_upload = false

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

	has, err := engine.Exist(&Workorders{ Workorder_id: id})
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

	sql := "update `results` set cur_upload = ?, result_upload = ?, update_time = ?, result_data = ?, cur_data = ?, controller_sn = ?, result = ?, need_upload = ? where result_id = ? and count = ?"
	_, err = engine.Exec(sql, result.Cur_upload, result. Result_upload, result.Update_time, result.Result_data, result.Cur_data, result.Controller_sn, result.Result, result.Need_upload, result.Result_id, result.Count)

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}