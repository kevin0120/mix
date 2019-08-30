package openprotocol

import "github.com/masami10/rush/services/tightening_device"

type TighteningTool struct {
	tightening_device.TighteningTool
}

// 工具使能控制
func (s *TighteningTool) ToolControl(enable bool) error { return nil }

// 设置PSet
func (s *TighteningTool) SetPSet(pset int) error { return nil }

// 设置Job
func (s *TighteningTool) SetJob(job int) error { return nil }
