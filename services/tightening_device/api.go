package tightening_device

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
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
	StepID       int64  `json:"step_id"`
	Job          int    `json:"job"`
	UserID       uint   `json:"user_id"`
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

	if req.WorkorderID > 0 {
		err = tool.TraceSet(fmt.Sprintf("%d", req.WorkorderID))
		if err != nil {
			return err
		}
	}

	return tool.SetJob(req.Job)
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

	if req.WorkorderID > 0 {
		err = tool.TraceSet(fmt.Sprintf("%d", req.WorkorderID))
		if err != nil {
			return err
		}
	}

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
