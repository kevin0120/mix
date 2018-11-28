package changan

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TighteningResults struct {
	gorm.Model
	WorkcenterCode string
	ControllerSn   string
	ToolSn         string
	Result         string
	Exception      string
	Torque         float64
	Angle          float64
	Spent          int
	UpdateTime     time.Time
	Strategy       string
	TorqueMax      float64
	TorqueMin      float64
	TorqueTarget   float64
	AngleMax       float64
	AngleMin       float64
	AngleTarget    float64
	Batch          string
	Mode           string
	Vin            string
	Cartype        string
	TighteningId   int64
}
