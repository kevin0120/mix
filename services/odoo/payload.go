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
	Points            interface{} `json:"points"`
	TignteningStepRef string      `json:"tightening_step_ref"`
	ProductTypeImage  string      `json:"product_type_image"`
}

type Maintenance struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Expire string `json:"expire_time"`
}

type OdooWorkorder struct {
}

type OdooStep struct {
}

type OdooUser struct {
	Status     string `json:"status"`
	Name       string `json:"name"`
	ImageSmall string `json:"image_small"`
	Login      string `json:"login"`
	UUID       string `json:"uuid"`
}
