package tightening_device

import (
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/storage"
)

// api接口
type Reply struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type JobSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	WorkorderID  int64  `json:"workorder_id"`
	Total        int    `json:"total"`
	StepID       int64  `json:"step_id"`
	Job          int    `json:"job"`
	UserID       int64  `json:"user_id"`
}

func (s *JobSet) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SN or Tool SN is required")
	}

	if s.Job <= 0 {
		return errors.New("Job Should Be Greater Than Zero")
	}

	return nil
}

type PSetBatchSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	PSet         int    `json:"pset"`
	Batch        int    `json:"batch"`
}

type PSetSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	StepID       int64  `json:"step_id"`
	WorkorderID  int64  `json:"workorder_id"`
	UserID       int64  `json:"user_id"`
	PSet         int    `json:"pset"`
	Sequence     uint   `json:"sequence"`
	Count        int    `json:"count"`
	Total        int    `json:"total"`
}

func (s *PSetSet) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SN or Tool SN is required")
	}

	if s.PSet <= 0 {
		return errors.New("pset Should Be Greater Than Zero")
	}

	return nil
}

type ToolControl struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Enable       bool   `json:"enable"`
}

func (s *ToolControl) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SN or Tool SN is required")
	}

	return nil
}

type ToolModeSelect struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Mode         string `json:"mode"`
}

func (s *ToolModeSelect) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SN or Tool SN is required")
	}

	return nil
}

type Api struct {
	service *Service
}

func (s *Api) ToolControl(req *ToolControl) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.service.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.ToolControl(req.Enable)
}

func (s *Api) ToolJobSet(req *JobSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.service.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	_ = s.service.StorageService.UpdateTool(&storage.Guns{
		Serial:      req.ToolSN,
		WorkorderID: req.WorkorderID,
		Total:       req.Total,
		UserID:      req.UserID,
	})

	return tool.SetJob(req.Job)
}

func (s *Api) ToolPSetBatchSet(req *PSetBatchSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	tool, err := s.service.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.SetPSetBatch(req.PSet, req.Batch)
}

func (s *Api) ToolPSetSet(req *PSetSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.service.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	_ = s.service.StorageService.UpdateTool(&storage.Guns{
		Serial:      req.ToolSN,
		WorkorderID: req.WorkorderID,
		Seq:         int(req.Sequence),
		Count:       req.Count,
		UserID:      req.UserID,
		Total:       req.Total,
	})

	return tool.SetPSet(req.PSet)
}

func (s *Api) ToolModeSelect(req *ToolModeSelect) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.service.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.ModeSelect(req.Mode)
}
