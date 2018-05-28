package core

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"github.com/masami10/rush/db"
	"strings"
	"encoding/xml"
	"strconv"
	"encoding/json"
	"time"
	"github.com/masami10/rush/storage"
	"github.com/masami10/rush/payload"
)

type ControllerConf struct {
	SN string	`yaml:"sn"`
	IP string	`yaml:"ip"`
	Port uint	`yaml:"port"`
}

type CVI3Conf struct {
	Listen int `yaml:"listen"`
	Controllers []ControllerConf `yaml:"controllers"`
}

func (service *CVI3Service) OnStatus(sn string, status string) {
	fmt.Printf("%s:%s\n", sn, status)

	// ws推送状态
	s := payload.WSStatus{}
	s.Status = status
	s.SN = sn
	msg, _ := json.Marshal(s)

	service.APIService.WSSendStatus(string(msg))
}

func (service *CVI3Service) OnRecv(msg string) {
	fmt.Printf("%s\n", msg)

	if strings.Contains(msg, cvi3.XML_RESULT_KEY) {
		result := cvi3.CVI3Result{}
		err := xml.Unmarshal([]byte(msg), &result)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		// 波形文件
		cur_result := cvi3.CVI3CurveFile{}
		cur_result.Result = result.PRC_SST.PAR.Result
		if cur_result.Result == "IO" {
			cur_result.Result = payload.RESULT_OK
		} else if cur_result.Result == "NIO" {
			cur_result.Result = payload.RESULT_NOK
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

		cur_file, _ := json.Marshal(cur_result)

		// 结果数据
		result_data := payload.XMl2Result(result)

		// 本地缓存结果和波形
		r, _ := service.DB.GetResult(result_data.Result_id, result_data.Count)

		r.Result = result_data.Result
		loc, _ := time.LoadLocation("Local")
		dt := fmt.Sprintf("%s %s", result.PRC_SST.PAR.FAS.GRP.TIP.Date, result.PRC_SST.PAR.FAS.GRP.TIP.Time)
		r.Update_time, _ = time.ParseInLocation("2006-01-02 15:04:05", dt, loc)
		rd, _ := json.Marshal(result_data)
		r.Result_data = string(rd)

		r.Cur_data = string(cur_file)
		r.Controller_sn = result_data.Controller_SN
		r.Count = result_data.Count

		// 如果是最终结果，则设置发送标识
		if r.Count >= r.Total_count || r.Result == payload.RESULT_OK{
			r.Need_upload = true
		}

		_, err = service.DB.UpdateResult(r)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		// 找到结果对应的hmi
		var w rushdb.Workorders
		w, err = service.DB.GetWorkorder(result_data.Workorder_ID)
		if err == nil {
			// ws推送结果到hmi
			ws_result := payload.WSResult{}
			ws_result.Result_id = result_data.Result_id
			ws_result.Count = result_data.Count
			ws_result.Result = result_data.Result
			ws_result.MI = result_data.ResultValue.Mi
			ws_result.WI = result_data.ResultValue.Wi
			ws_result.TI = result_data.ResultValue.Ti
			ws_str, _ := json.Marshal(ws_result)
			go service.APIService.WSSendResult(w.HMI_sn, string(ws_str))
		}

		// 保存波形到对象存储
		err = service.Storage.Upload(result_data.CurFile, r.Cur_data)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		} else {
			r.Cur_upload = true
			_, err = service.DB.UpdateResult(r)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}

		// 加入odoo推送队列
		odoo_result := payload.ODOOResult{}
		odoo_result.Control_date = result_data.Dat

		odoo_result.CURObjects = []payload.CURObject{}


		results, _ := service.DB.ListResults(result_data.Result_id)
		for _,v := range results {
			nr := payload.Result{}
			json.Unmarshal([]byte(v.Result_data), &nr)
			cur_object := payload.CURObject{}
			cur_object.File = nr.CurFile
			cur_object.OP = nr.Count
			odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
		}

		odoo_result.Measure_degree = result_data.ResultValue.Wi
		odoo_result.Measure_result = strings.ToLower(result_data.Result)
		odoo_result.Measure_t_don = result_data.ResultValue.Ti
		odoo_result.Measure_torque = result_data.ResultValue.Mi
		odoo_result.Op_time = result_data.Count
		odoo_result.Pset_m_max = result_data.PSetDefine.Mp
		odoo_result.Pset_m_min = result_data.PSetDefine.Mm
		odoo_result.Pset_m_target = result_data.PSetDefine.Ma
		odoo_result.Pset_m_threshold = result_data.PSetDefine.Ms
		odoo_result.Pset_strategy = result_data.PSetDefine.Strategy
		odoo_result.Pset_w_max = result_data.PSetDefine.Wp
		odoo_result.Pset_w_min = result_data.PSetDefine.Wm
		odoo_result.Pset_w_target = result_data.PSetDefine.Wa

		result_put := payload.ODOORsultPut{}
		result_put.ID = result_data.Result_id
		result_put.Result = odoo_result

		service.ODOO.PutStack.Push(result_put)
		//
		//go service.ODOO.PutResult(result_data.Result_id, odoo_result)

		// 追加最新波形
		//go service.ODOO.PatchCurve(result_data.Result_id, result_data.CurFile, result_data.Count)

	}


}

type CVI3Service struct {
	Service	cvi3.CVI3
	Port	string
	Configs []cvi3.CVI3Config
	hmis	map[string]string
	DB		*rushdb.DB
	Storage	*rush_storage.Storage
	ODOO	*ODOO
	APIService *APIServer
}

func (service *CVI3Service) Config(configs []cvi3.CVI3Config) {
	service.Configs = configs
}

func (service *CVI3Service) StartService() error {
	service.Service = cvi3.CVI3{}
	service.Service.Config(service.Configs)
	service.Service.RegisterCallBack(service.OnStatus, service.OnRecv)

	return service.Service.StartService(service.Port)
}