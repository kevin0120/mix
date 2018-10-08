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
	Torque         float32
	Angle          float32
	Spent          int
	UpdateTime     time.Time
	Strategy       string
	TorqueMax      float32
	TorqueMin      float32
	TorqueTarget   float32
	AngleMax       float32
	AngleMin       float32
	AngleTarget    float32
	Batch          string
	Mode           string
	Vin            string
	Cartype        string
	TighteningId   int64
}
