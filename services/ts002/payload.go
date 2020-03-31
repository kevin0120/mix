package ts002

// ************************************************Mes提供的接口Payload************************************************
const (
	MES_CARDINFO_FAIL = "E"
)

// Mes->Rush 结果响应
type MesCardInfoResp struct {
	ResultStatus string `json:"resultStatus"`
	ResultMsg    string `json:"resultMsg"`
}

// Rush->Mes 发送刷卡数据请求
type MesCardInfoReq struct {
	CardCode string `json:"cardCode"`
}

// Rush->Mes 扭矩结果上传
type MesResultUploadReq struct {
	UUID         string  `json:"UUID"`
	ActualTorque float64 `json:"actualTorque"`
	ActualAngle  float64 `json:"actualAngle"`
	Flag         string  `json:"flag"`
}

// ************************************************Rush提供的接口Payload************************************************

// Rush->Mes 错误响应
type RushErrResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// Mes->Rush 报警控制请求
type RushAlarmReq struct {
	Status string `json:"status" validate:"required"`
}

// Mes->Rush PSet下发请求
type RushPSetReq struct {
	PSet    int    `json:"pset" validate:"required,gte=1"`
	ToolID  string `json:"tool_id" validate:"required"`
	PointID string `json:"point_id" validate:"required"`
}

// Mes->Rush IO输出控制
type RushIOControlReq struct {
	Outputs []int  `json:"outputs" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

// Mes->Rush 读卡器信息
type RushCardInfoReq struct {
	CardCode string `json:"cardCode" validate:"required"`
}

// odoo equipment
type Equipment struct {
	EquipmentSN string   `json:"equipment_sn"`
	Location    Location `json:"location"`
	Type        string   `json:"type"`
}

type Location struct {
	EquipmentSN string `json:"equipment_sn"`
	Input       int    `json:"io_input"`
	Output      int    `json:"io_output"`
}
