package fis

import (
	"github.com/masami10/aiis/toml"
	"time"
)

type Config struct {
	CHRecvMission   string 			`yaml:"ch_recv_mission"`
	CHSendResult    string 			`yaml:"ch_send_result"`
	CHRecvHeartbeat string 			`yaml:"ch_recv_heartbeat"`
	SystemType      string 			`yaml:"system_type"`
	SoftwareVersion string 			`yaml:"software_version"`
	Mode            string 			`yaml:"mode"`
	FactoryCode		string 			`yaml:"factory_code"`
	PRS				[]string		`yaml:"prs"`
	HeartbeatItv	toml.Duration	`yaml:"heartbeat_interval"`
}

func NewConfig() Config {
	prs := []string{}
	prs = append(prs, "GSP")

	return Config{
		CHRecvMission:   "02",
		CHSendResult:    "01",
		CHRecvHeartbeat: "03",
		SystemType:      "screw",
		SoftwareVersion: "1.0",
		Mode:            "AUTO",
		FactoryCode:	 "01",
		PRS:			 prs,
		HeartbeatItv:	 toml.Duration(time.Minute * 1),
	}
}

func (c Config) Validate() error {

	return nil
}
