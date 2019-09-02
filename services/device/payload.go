package device

const (
	STATUS_ONLINE  = "online"
	STATUS_OFFLINE = "offline"
)

const (
	WS_DEVICE_STATUS = "WS_DEVICE_STATUS"
)

type DeviceStatus struct {
	SN       string      `json:"sn"`
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	Children interface{} `json:"children"`
	Config   interface{} `json:"config"`
	Data     interface{} `json:"data"`
}
