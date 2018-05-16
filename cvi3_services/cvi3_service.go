package rush_cvi3

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
	"github.com/masami10/rush/odoo"
)

func (service *CVI3Service) OnStatus(sn string, status string) {
	fmt.Printf("%s:%s\n", sn, status)

	// ws推送状态
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

		_, err = service.DB.UpdateResult(r)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
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

		// 结果推送odoo
		odoo_result := payload.ODOOResult{}
		odoo_result.Control_date = result_data.Dat
		cur_object := payload.CURObject{}
		cur_object.File = result_data.CurFile
		cur_object.OP = result_data.Count
		odoo_result.CURObjects = []payload.CURObject{}
		odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
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

		go service.ODOO.PutResult(result_data.Result_id, odoo_result)


		// ws推送结果到hmi
	}


}

type CVI3Service struct {
	Service	cvi3.CVI3
	Port	string
	Configs []cvi3.CVI3Config
	hmis	map[string]string
	DB		*rushdb.DB
	Storage	*rush_storage.Storage
	ODOO	*rush_odoo.ODOO
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