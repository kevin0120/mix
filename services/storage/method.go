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
	session.Begin()
	var work Workorders
	var workPayload WorkorderPayload
	err := json.Unmarshal(in, &work)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(in, &workPayload)
	if err != nil {
		return "", err
	}

	wp, err := json.Marshal(workPayload)
	//dt := time.Now()
	//if data.controllerResult.Dat != "" {
	//	loc, _ := time.LoadLocation("Local")
	//	dt, _ = time.ParseInLocation("2006-01-02 15:04:05", data.controllerResult.Dat, loc)
	//}
	//
	//dbResult.UpdateTime = dt.UTC()
	//

	workorder1 := Workorders{
		Workorder:    string(wp),
		Code:         work.Code,
		Track_code:   work.Track_code,
		Product_code: work.Product_code,
		//Workcenter:            work.Workcenter,
		Status:                "todo",
		Date_planned_start:    work.Date_planned_start,
		Date_planned_complete: work.Date_planned_start,
	}
	//
	s.DeleteWorkAndStep(session, work.Code)

	//var hh map[string]interface{}
	//
	//err = json.Unmarshal(wor, &hh)
	_, err = session.Insert(&workorder1)
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

	steps, _ := json.Marshal(workorderMap["steps"])

	err = json.Unmarshal(steps, &step)

	for i := 0; i < len(step); i++ {
		stepString, _ := json.Marshal(step[i])
		var msg Steps
		var stepText StepTextPayload
		var stepTightening StepTighteningPayload
		var sp []byte
		err = json.Unmarshal(stepString, &msg)

		if msg.Testtype == "tightening" {
			err = json.Unmarshal(stepString, &stepTightening)
			sp, _ = json.Marshal(stepTightening)
		} else {
			err = json.Unmarshal(stepString, &stepText)
			sp, _ = json.Marshal(stepText)
		}

		step := Steps{
			WorkorderID:    workorder1.Id,
			Step:           string(sp),
			ImageRef:       msg.ImageRef,
			Testtype:       msg.Testtype,
			Status:         "ready",
			Code:           msg.Code,
			Sequence:       msg.Sequence,
			FailureMessage: msg.FailureMessage,
			Desc:           msg.Desc,
			Skippable:      msg.Skippable,
			Undoable:       msg.Undoable,
			Data:           msg.Data,
		}

		_, err := session.Insert(&step)
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

func (s *Service) WorkorderOut(order string, workorderID int64) (interface{}, error) {

	var workorder Workorders
	var ss *xorm.Session
	if order == "" {
		ss = s.eng.Alias("r").Where("r.id = ?", workorderID)
	} else {
		ss = s.eng.Alias("r").Where("r.code = ?", order)
	}

	bool, e := ss.Get(&workorder)
	if e != nil || !bool {
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
		stepOut := strucToMap(step[i])
		stepOut["payload"] = stepMap

		if step[i].Testtype != "tightening" {
			delete(stepOut, "tightening_image_by_step_code")
			steps = append(steps, stepOut)
			continue
		}
		image1, _ := s.findStepPicture(step[i].ImageRef)
		stepOut["image"] = image1
		delete(stepOut, "tightening_image_by_step_code")
		steps = append(steps, stepOut)
	}

	map2 := stringToMap(workorder.Workorder)
	workOrderOut := strucToMap(workorder)
	workOrderOut["steps"] = steps
	workOrderOut["payload"] = map2

	image2, _ := s.findOrderPicture(workorder.Product_code)
	workOrderOut["product_type_image"] = image2
	//delete(workOrderOut,"product_code")
	//rr, _ := json.Marshal(workOrderOut)

	return workOrderOut, nil
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
		return "", e
	}
	return ro.Img, nil
}

func (s *Service) findOrderPicture(ref string) (string, error) {
	ro, err := s.GetRoutingOperationViaProductTypeCode(ref)
	if err != nil {
		err := errors.Errorf("GetRoutingOperationViaProductTypeCode: %s Fail", ref)
		s.diag.Error("findOrderPicture", err)
		return "", err
	}
	return ro.ProductTypeImage, nil
}

func (s *Service)GetRoutingOperationViaProductTypeCode(ProductType string) (*RoutingOperations, error) {
	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.product_type = ?", ref).Limit(1)
	_, e := ss.Get(&ro)
	if e != nil {
		return nil, errors.Errorf("Operation For Product Type Code: %s Is Not Existed", ProductType)
	}
	return &ro, nil
}
