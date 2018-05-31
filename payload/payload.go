package payload

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"strconv"
	"container/list"
	"sync"
	"strings"
)

const (
	RESULT_NONE = "NONE"
	RESULT_OK = "OK"
	RESULT_NOK = "NOK"
)

const (
	RESULT_STAGE_INIT = "init"
	RESULT_STAGE_FINAL = "final"
)

const (
	KNR_KEY = "-"
)

type PSet struct {
	Controller_SN string `json:"controller_sn"`
	PSet int `json:"pset"`
	Result_id int `json:"result_id"`
	Count int `json:"count"`
}

type ResultPatch struct {
	HasUpload bool `json:"has_upload"`
}

type Workorder struct {
	Workorder_id	int			`json:"workorder_id"`
	HMI_sn			string		`json:"hmi_sn"`
	PSet			int			`json:"pset"`
	Nut_total		float64		`json:"nut_total"`
	Vin				string		`json:"vin"`
	Knr				string		`json:"knr"`
	Status			string		`json:"status"`
	Result_ids		[]int		`json:"result_ids"`
	WorkSheet		string		`json:"work_sheet"`
}

type ControllerResult struct {
	Result_id int	`json:"result_id"`
	Controller_SN string `json:"controller_sn"`
	Workorder_ID int `json:"workorder_id"`
	CurFile string `json:"cur_file"`
	Result string `json:"result"`
	Dat string `json:"dat"`
	PSet int `json:"pset"`
	Count int `json:"count"`
	PSetDefine PSetDefine `json:"pset_define"`

	ResultValue ResultValue `json:"result_value"`
}

type PSetDefine struct {
	Strategy string `json:"strategy"`
	Mp float64 `json:"M+"`
	Mm float64 `json:"M-"`
	Ms float64 `json:"MS"`
	Ma float64 `json:"MA"`
	Wp float64 `json:"W+"`
	Wm float64 `json:"W-"`
	Wa float64 `json:"WS"`
}

type ResultValue struct {
	Mi float64 `json:"MI"`
	Wi float64 `json:"WI"`
	Ti float64 `json:"TI"`
}

type ControllerCurve struct {
	ResultID int
	CurveFile string
	CurveData string
	Count int
}

type ControllerCurveFile struct {
	Result string	`json:"result"`
	CUR_M []float64 `json:"cur_m"`
	CUR_W []float64 `json:"cur_w"`
}

func XML2Curve (result cvi3.CVI3Result) (ControllerCurveFile) {
	cur_result := ControllerCurveFile{}
	cur_result.Result = result.PRC_SST.PAR.Result
	if cur_result.Result == "IO" {
		cur_result.Result = RESULT_OK
	} else if cur_result.Result == "NIO" {
		cur_result.Result = RESULT_NOK
	}

	cur_ms := strings.Split(result.PRC_SST.PAR.FAS.GRP.TIP.BLC.CUR.SMP.CUR_M, " ")
	for i := range cur_ms {
		v, _ := strconv.ParseFloat(cur_ms[i], 64)
		cur_result.CUR_M = append(cur_result.CUR_M, v)
	}

	cur_ws := strings.Split(result.PRC_SST.PAR.FAS.GRP.TIP.BLC.CUR.SMP.CUR_W, " ")
	for i := range cur_ws {
		v, _ := strconv.ParseFloat(cur_ws[i], 64)
		cur_result.CUR_W = append(cur_result.CUR_W, v)
	}

	return cur_result
}

func XML2Result(result cvi3.CVI3Result) (ControllerResult) {
	rr := ControllerResult{}

	rr.Controller_SN = result.PRC_SST.PAR.SN
	rr.Result = result.PRC_SST.PAR.Result
	if rr.Result == "IO" {
		rr.Result = RESULT_OK
	} else if rr.Result == "NIO" {
		rr.Result = RESULT_NOK
	}

	rr.PSet = result.PRC_SST.PAR.FAS.GRP.TIP.PSet
	rr.Workorder_ID = result.PRC_SST.PAR.Workorder_id
	rr.Dat = fmt.Sprintf("%s %s", result.PRC_SST.PAR.FAS.GRP.TIP.Date, result.PRC_SST.PAR.FAS.GRP.TIP.Time)
	result_id := result.PRC_SST.PAR.Result_id
	rr.Result_id, _ = strconv.Atoi(result_id)
	rr.CurFile = fmt.Sprintf("%s_%d_%s_%s.json", rr.Controller_SN, rr.Workorder_ID, result_id, cvi3.GenerateID())
	rr.PSetDefine.Strategy = result.PRC_SST.PAR.FAS.GRP.TIP.BLC.PRO.Strategy
	rr.Count = result.PRC_SST.PAR.Count


	result_values := result.PRC_SST.PAR.FAS.GRP.TIP.BLC.PRO.Values
	for i := range result_values {
		switch result_values[i].Name {
		case "M+":
			rr.PSetDefine.Mp = result_values[i].Value
		case "M-":
			rr.PSetDefine.Mm = result_values[i].Value
		case "MS":
			rr.PSetDefine.Ms = result_values[i].Value
		case "MA":
			rr.PSetDefine.Ma = result_values[i].Value
		case "W+":
			rr.PSetDefine.Wp = result_values[i].Value
		case "W-":
			rr.PSetDefine.Wm = result_values[i].Value
		case "WA":
			rr.PSetDefine.Wa = result_values[i].Value
		case "MI":
			rr.ResultValue.Mi = result_values[i].Value
		case "WI":
			rr.ResultValue.Wi = result_values[i].Value
		case "tI":
			if result_values[i].Unit == "s" {
				rr.ResultValue.Ti = result_values[i].Value * 1000
			} else {
				rr.ResultValue.Ti = result_values[i].Value
			}
		}
	}

	return rr
}

type ODOOMOCreated struct {
	ID int				`json:"id"`
	KNR string			`json:"knr"`
	VIN string			`json:"vin"`
	ProductID int		`json:"product_id"`
	Result_IDs []int	`json:"result_ids"`
	AssembleLine int 	`json:"assembly_line_id"`
}

type ODOOWorkorder struct {
	ID	int				`json:"id"`
	Status string		`json:"status"`
	NutTotal float64	`json:"nut_total"`
	PSet string			`json:"pset"`
	Max_redo_times	int	`json:"max_redo_times"`
	Max_op_time	int		`json:"max_op_time"`

	HMI struct {
		ID 		int		`json:"id"`
		UUID string		`json:"uuid"`
	} `json:"hmi"`

	Result_IDs []int	`json:"result_ids"`
	Worksheet struct {
		Content string	`json:"content"`
		Points []struct {
			X int32 `json:"x_offset"`
			Y int32 `json:"y_offset"`
		} `json:"points"`

	}	`json:"worksheet"`
	KNR string			`json:"knr"`
	VIN string			`json:"vin"`
}

type CURObject struct {
	File string  `json:"file"`
	OP int  `json:"op"`
}

type ODOOCurveAppend struct {
	File string `json:"file"`
	OP int `json:"op"`
}

type ODOOResultSync struct {
	ID	int	`json:"id"`
	Pset_m_threshold	float64 `json:"pset_m_threshold"`
	Pset_m_max	float64 `json:"pset_m_max"`
	Pset_m_min	float64 `json:"pset_m_min"`
	Control_date	string `json:"control_date"`
	Pset_w_max	float64 `json:"pset_w_max"`
	Pset_strategy	string `json:"pset_strategy"`
	Pset_m_target	float64 `json:"pset_m_target"`
	Measure_degree	float64 `json:"measure_degree"`
	Measure_t_don	float64 `json:"measure_t_don"`
	Measure_torque	float64 `json:"measure_torque"`
	Measure_result	string `json:"measure_result"`
	Op_time	int `json:"op_time"`
	Pset_w_min	float64 `json:"pset_w_min"`
	Pset_w_target	float64 `json:"pset_w_target"`


	CURObjects	[]CURObject	`json:"cur_objects"`
}

type ODOOResult struct {
	Pset_m_threshold	float64 `json:"pset_m_threshold"`
	Pset_m_max	float64 `json:"pset_m_max"`
	Pset_m_min	float64 `json:"pset_m_min"`
	Control_date	string `json:"control_date"`
	Pset_w_max	float64 `json:"pset_w_max"`
	Pset_strategy	string `json:"pset_strategy"`
	Pset_m_target	float64 `json:"pset_m_target"`
	Measure_degree	float64 `json:"measure_degree"`
	Measure_t_don	float64 `json:"measure_t_don"`
	Measure_torque	float64 `json:"measure_torque"`
	Measure_result	string `json:"measure_result"`
	Op_time	int `json:"op_time"`
	Pset_w_min	float64 `json:"pset_w_min"`
	Pset_w_target	float64 `json:"pset_w_target"`
	Final_pass	string	`json:"final_pass"`
	One_time_pass	string	`json:"one_time_pass"`

	CURObjects	[]CURObject	`json:"cur_objects"`
}

type ODOOPR struct {
	Pr_value string `json:"pr_value"`
	Pr_group string `json:"pr_group"`
}

type ODOOMO struct {
	Pin_check_code int	`json:"pin_check_code"`
	Equipment_name string	`json:"equipment_name"`
	Factory_name string	`json:"factory_name"`
	Pin int	`json:"pin"`
	Year int `json:"year"`
	Assembly_line string `json:"assembly_line"`
	Model string `json:"model"`
	//Model2 string `json:"model"`
	Vin string `json:"vin"`
	Lnr string `json:"lnr"`
	Date_planned_start string `json:"date_planned_start"`
	Prs []ODOOPR `json:"prs"`
}

type WSStatus struct {
	SN string  `json:"controller_sn"`
	Status string  `json:"status"`
}

type WSResult struct {
	Result_id int	`json:"result_id"`
	Count int	`json:"count"`
	Result string  `json:"result"`
	MI float64  `json:"mi"`
	WI float64  `json:"wi"`
	TI float64  `json:"ti"`
}

type WSRegist struct {
	HMI_SN string	`json:"hmi_sn"`
}

type WSRegistMsg struct {
	Msg string	`json:"msg"`
}

type ODOORsultPut struct {
	ID int
	Result ODOOResult
}

type SafeStack struct {
	List *list.List
	Mtx sync.Mutex
}

func (stack *SafeStack) Init() {
	stack.List = list.New()
	stack.Mtx = sync.Mutex{}
}

func (stack *SafeStack) Push(value interface{}) {
	defer stack.Mtx.Unlock()

	stack.Mtx.Lock()
	stack.List.PushBack(value)
}

func (stack *SafeStack) Pop() interface{} {
	defer stack.Mtx.Unlock()

	stack.Mtx.Lock()

	e := stack.List.Back()
	if e != nil {
		stack.List.Remove(e)
		return e.Value
	}

	return nil
}