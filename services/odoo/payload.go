package odoo

type RoutingOperation struct {
	OperationID       int64       `json:"id"`
	Job               int         `json:"job"`
	MaxOpTime         int         `json:"max_op_time"`
	Name              string      `json:"name"`
	Img               string      `json:"img"`
	ProductId         int64       `json:"product_id"`
	WorkcenterID      int64       `json:"workcenter_id"`
	ProductType       string      `json:"product_type"`
	WorkcenterCode    string      `json:"workcenter_code"`
	VehicleTypeImg    string      `json:"vehicleTypeImg"`
	Points            interface{} `json:"points"`
	TignteningStepRef string      `json:"tightening_step_ref"`
	ProductTypeImage  string      `json:"product_type_image"`
}

type Maintenance struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Expire string `json:"expire_time"`
}

type TighteningStep struct {
	TighteningPoints []StepComsume `json:"tightening_points"`
}

type StepComsume struct {
	Seq      int    `json:"sequence"`
	GroupSeq int    `json:"group_sequence"`
	NutNo    string `json:"nut_no"`
}
