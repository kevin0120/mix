package fis

import (
	"github.com/masami10/aiis/toml"
	"time"
)

type Config struct {
	CHRecvMission   string        `yaml:"ch_recv_mission"`
	CHSendResult    string        `yaml:"ch_send_result"`
	CHRecvHeartbeat string        `yaml:"ch_recv_heartbeat"`
	CHSendUrg       string        `yaml:"ch_send_urg"`
	CHRecvUrg       string        `yaml:"ch_recv_urg"`
	EquipmentName   string        `yaml:"equipment_name"`
	SystemType      string        `yaml:"system_type"`
	SoftwareVersion string        `yaml:"software_version"`
	Mode            string        `yaml:"mode"`
	FactoryCode     string        `yaml:"factory_code"`
	PRS             []string      `yaml:"prs"`
	HeartbeatItv    toml.Duration `yaml:"heartbeat_interval"`
	Enable          bool          `yaml:"enable"`
	MissionItv      toml.Duration `yaml:"mission_itv"`
	UrgTimeout      toml.Duration `yaml:"urg_timeout"`
	ADTimeout       toml.Duration `yaml:"ad_timeout"`
}

func NewConfig() Config {
	prs := []string{}
	prs = append(prs, "GSP")

	return Config{
		CHRecvMission:   "02",
		CHSendResult:    "01",
		CHRecvHeartbeat: "03",
		CHSendUrg:       "04",
		CHRecvUrg:       "05",
		EquipmentName:   "SR3J",
		SystemType:      "screw",
		SoftwareVersion: "1.0",
		Mode:            "AUTO",
		FactoryCode:     "01",
		PRS:             prs,
		HeartbeatItv:    toml.Duration(time.Minute * 1),
		Enable:          false,
		MissionItv:      toml.Duration(10 * time.Millisecond),
		UrgTimeout:      toml.Duration(300 * time.Millisecond),
		ADTimeout:       toml.Duration(1 * time.Second),
	}
}

func (c Config) Validate() error {

	return nil
}
