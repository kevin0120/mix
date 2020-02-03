package device

type Status struct {
	SN       string      `json:"sn"`
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	Children interface{} `json:"children"`
	Config   interface{} `json:"config"`
	Data     interface{} `json:"data"`
}
