package storage

import (
	"encoding/json"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/

func (s *Service) WorkorderIn(in []byte) (string, error) {

	session := s.eng.NewSession()
	defer session.Close()

	var work Workorders
	err := json.Unmarshal(in, &work)

	workorder1 := Workorders{
		Workorder:    string(in),
		Code:         work.Code,
		Track_code:   work.Track_code,
		Product_code: work.Product_code,
		Status:       "ready",
	}
	//
	s.DeleteWorkAndStep(work.Code)

	//var hh map[string]interface{}
	//
	//err = json.Unmarshal(wor, &hh)
	_, err = s.eng.Insert(&workorder1)
	// INSERT INTO struct () values ()
	//engine.Ping()
	//有的数据库超时断开ping可以重连。可以通过起一个定期Ping的Go程来保持连接鲜活。
	if err != nil {
		session.Rollback()
		return "", errors.Wrapf(err, "store data fail")
	}

	var workorderMap map[string]interface{}
	var step []map[string]interface{}

	err = json.Unmarshal(in, &workorderMap)

	cc, _ := json.Marshal(workorderMap["steps"])

	err = json.Unmarshal(cc, &step)

	for i := 0; i < len(step); i++ {
		a, _ := json.Marshal(step[i])
		var msg Steps
		err = json.Unmarshal(a, &msg)

		step := Steps{
			WorkorderID: workorder1.Id,
			Step:        string(a),
			Image:       msg.Image,

			Test_type: msg.Test_type,
			Status:    "ready",
			Code:      msg.Code,
		}

		_, err := s.eng.Insert(&step)
		// INSERT INTO struct () values ()
		if err != nil {
			session.Rollback()
			return "", errors.Wrapf(err, "store data fail")
		}

	}

	err = session.Commit()
	if err != nil {
		return "", errors.Wrapf(err, "commit fail")
	}

	return workorder1.Code, nil

}

func (s *Service) WorkorderOut(order string, workorderID int64) ([]byte, error) {

	var workorder Workorders
	var ss *xorm.Session
	if order == "" {
		ss = s.eng.Alias("r").Where("r.id = ?", workorderID)
	} else {
		ss = s.eng.Alias("r").Where("r.code = ?", order)
	}

	_, e := ss.Get(&workorder)
	if e != nil {
		return nil, e
	}

	var step []Steps
	ss = s.eng.Alias("r").Where("r.x_workorder_id = ?", workorder.Id)
	e = ss.Find(&step)
	if e != nil {
		return nil, e
	}

	var steps []map[string]interface{}
	for i := 0; i < len(step); i++ {
		stepMap := stringToMap(step[i].Step)
		testType, _ := json.Marshal(stepMap["test_type"])

		if !(string(testType) == `"tightening"`) {
			steps = append(steps, stepMap)
			continue
		}

		map1 := strucToMap(step[i])
		for k, v := range map1 {
			stepMap[k] = v
		}
		image1, _ := s.findStepPicture(step[i].ImageRef)
		stepMap["image"] = image1
		steps = append(steps, stepMap)
	}

	workOrderOut := stringToMap(workorder.Workorder)
	map2 := strucToMap(workorder)
	workOrderOut["steps"] = steps
	for k, v := range map2 {
		workOrderOut[k] = v
	}

	image2, _ := s.findOrderPicture(workorder.Product_code)
	workOrderOut["product_type_image"] = image2

	rr, _ := json.Marshal(workOrderOut)

	return rr, nil
}

func strucToMap(in interface{}) (m map[string]interface{}) {
	j, _ := json.Marshal(in)
	json.Unmarshal(j, &m)
	return
}

func stringToMap(in string) (m map[string]interface{}) {
	json.Unmarshal([]byte(in), &m)
	return
}

func (s *Service) findStepPicture(ref string) (string, error) {

	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.tightening_step_ref = ?", ref).Limit(1)
	_, e := ss.Get(&ro)
	if e != nil {
		return ro.Img, e
	}
	return "", nil
}

func (s *Service) findOrderPicture(ref string) (string, error) {
	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.product_type = ?", ref).Limit(1)
	_, e := ss.Get(&ro)
	if e != nil {
		return ro.ProductTypeImage, e
	}
	return "", nil
}
