package typeDef

//为了创建工单所使用的Step模型
type StepTighteningPayload struct {
	TighteningTotal int                     `json:"tightening_total"`
	TighteningPoint []RoutingOperationPoint `json:"tightening_points"`
}

type JobPoint struct {
	Seq                int     `json:"sequence"`
	PSet               int     `json:"pset"`
	X                  float64 `json:"offset_x"`
	Y                  float64 `json:"offset_y"`
	MaxRedoTimes       int     `json:"max_redo_times"`
	GroupSeq           int     `json:"group_sequence"`
	ConsuProductID     int64   `json:"consu_product_id"`
	ToleranceMin       float64 `json:"tolerance_min"`
	ToleranceMax       float64 `json:"tolerance_max"`
	ToleranceMinDegree float64 `json:"tolerance_min_degree"`
	ToleranceMaxDegree float64 `json:"tolerance_max_degree"`
}

type RoutingOperationPoint struct {
	JobPoint
	GunSN        string `json:"gun_sn"`
	ControllerSN string `json:"controller_sn"`
}

type RoutingOperation struct {
	OperationID    int64                   `json:"id"`
	Job            int                     `json:"job"`
	MaxOpTime      int                     `json:"max_op_time"`
	Name           string                  `json:"name"`
	Img            string                  `json:"img"`
	ProductId      int64                   `json:"product_id"`
	ProductType    string                  `json:"product_type"`
	WorkCenterCode string                  `json:"workcenter_code"`
	VehicleTypeImg string                  `json:"vehicleTypeImg"`
	Points         []RoutingOperationPoint `json:"points"`
}
