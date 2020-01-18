package storage

import (
	"encoding/json"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/

func (s *Service) WorkorderSync(work *Workorders) (string, error) {

	err := s.validator.Struct(work)
	if err != nil {
		return "", errors.Wrapf(err, "loss workorder-steps information")
	}

	session := s.eng.NewSession().ForUpdate()
	defer session.Close()
	session.Begin()

	roll, err := s.DeleteWorkAndStep(session, work.Code, work.UniqueNum)
	if roll {
		session.Rollback()
		return "", errors.Wrapf(err, "本地已有更新版本号对应的工单")
	}

	_, err = session.Insert(work)
	if err != nil {
		session.Rollback()
		return "", errors.Wrapf(err, "store data fail")
	}

	for i := 0; i < len(work.Steps); i++ {
		work.Steps[i].WorkorderID = work.Id
		_, err := session.Insert(work.Steps[i])

		if err != nil {
			session.Rollback()
			return "", errors.Wrapf(err, "store data fail")
		}

	}

	err = session.Commit()
	if err != nil {
		return "", errors.Wrapf(err, "commit fail")
	}

	return work.Code, nil
}

func (s *Service) WorkorderIn(in []byte) (string, error) {

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

	workorder1 := Workorders{
		Code:                work.Code,
		TrackCode:           work.TrackCode,
		ProductCode:         work.ProductCode,
		DatePlannedStart:    work.DatePlannedStart,
		DatePlannedComplete: work.DatePlannedComplete,
		UniqueNum:           work.UniqueNum,
		Workorder:           string(wp),

		Status: "todo",
	}

	var workorderMap map[string]interface{}
	var step []map[string]interface{}

	err = json.Unmarshal(in, &workorderMap)

	if _, exist := workorderMap["steps"]; exist {
		steps, _ := json.Marshal(workorderMap["steps"])
		err = json.Unmarshal(steps, &step)
	}

	for i := 0; i < len(step); i++ {
		stepString, _ := json.Marshal(step[i])
		var msg Steps
		var stepText StepTextPayload
		var stepTightening StepTighteningPayload
		var sp []byte
		err = json.Unmarshal(stepString, &msg)

		if msg.TestType == "tightening" {
			err = json.Unmarshal(stepString, &stepTightening)
			sp, _ = json.Marshal(stepTightening)
		} else {
			err = json.Unmarshal(stepString, &stepText)
			sp, _ = json.Marshal(stepText)
		}

		step := Steps{
			ImageRef:       msg.ImageRef,
			TestType:       msg.TestType,
			Code:           msg.Code,
			Sequence:       msg.Sequence,
			FailureMessage: msg.FailureMessage,
			Desc:           msg.Desc,
			Skippable:      msg.Skippable,
			Undoable:       msg.Undoable,
			Data:           msg.Data,

			Text:           msg.Text,
			ToleranceMin:   msg.ToleranceMin,
			ToleranceMax:   msg.ToleranceMax,
			Target:         msg.Target,
			ConsumeProduct: msg.ConsumeProduct,
			WorkorderID:    workorder1.Id,
			Step:           string(sp),

			Status: "ready",
		}

		workorder1.Steps = append(workorder1.Steps, step)
	}
	return s.WorkorderSync(&workorder1)
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

		if step[i].TestType != "tightening" {
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

	image2, _ := s.findOrderPicture(workorder.ProductCode)
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

func (s *Service) GetRoutingOperationViaProductTypeCode(ProductType string) (*RoutingOperations, error) {
	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.product_type = ?", ProductType).Limit(1)
	_, e := ss.Get(&ro)
	if e != nil {
		return nil, errors.Errorf("Operation For Product Type Code: %s Is Not Existed", ProductType)
	}
	return &ro, nil
}
