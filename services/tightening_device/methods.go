package tightening_device

import (
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/storage"
)

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
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
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
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
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
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
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
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
	}

	return nil
}

func (s *Service) ToolControl(req *ToolControl) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.ToolControl(req.Enable)
}

func (s *Service) ToolJobSet(req *JobSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	_ = s.storageService.UpdateTool(&storage.Tools{
		Serial:             req.ToolSN,
		CurrentWorkorderID: req.WorkorderID,
		Total:              req.Total,
		UserID:             req.UserID,
	})

	return tool.SetJob(req.Job)
}

func (s *Service) ToolPSetBatchSet(req *PSetBatchSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.SetPSetBatch(req.PSet, req.Batch)
}

func (s *Service) ToolPSetSet(req *PSetSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	_ = s.storageService.UpdateTool(&storage.Tools{
		Serial:             req.ToolSN,
		CurrentWorkorderID: req.WorkorderID,
		Seq:                int(req.Sequence),
		Count:              req.Count,
		UserID:             req.UserID,
		Total:              req.Total,
		StepID:             req.StepID,
	})

	return tool.SetPSet(req.PSet)
}

func (s *Service) ToolModeSelect(req *ToolModeSelect) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return tool.ModeSelect(req.Mode)
}

type ToolInfo struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
}

type ToolPSet struct {
	ToolInfo
	PSet int `json:"pset"`
}

type ToolJob struct {
	ToolInfo
	Job int `json:"job"`
}

func (s *Service) GetToolPSetList(req *ToolInfo) ([]int, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return tool.GetPSetList()
}

func (s *Service) GetToolPSetDetail(req *ToolPSet) (*PSetDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return tool.GetPSetDetail(req.PSet)
}

func (s *Service) GetToolJobList(req *ToolInfo) ([]int, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return tool.GetJobList()
}

func (s *Service) GetToolJobDetail(req *ToolJob) (*JobDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return tool.GetJobDetail(req.Job)
}
