package core

import (
	"fmt"
	"github.com/linshenqi/TightningSys/desoutter/cvi3"
	"github.com/masami10/rush/db"
	"strings"
	"encoding/xml"
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

func (service *CVI3Service) HandleResult(result payload.ControllerResult) (error) {
	var err error
	r, err := service.DB.GetResult(result.Result_id, 0)
	if err != nil {
		return err
	}

	// 保存结果
	loc, _ := time.LoadLocation("Local")
	times := strings.Split(result.Dat, " ")
	dt := fmt.Sprintf("%s %s", times[0], times[1])
	r.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", dt, loc)
	r.Result = result.Result
	r.Count = result.Count
	r.HasUpload = false
	r.ControllerSN = result.Controller_SN
	s_value, _ := json.Marshal(result.ResultValue)
	s_pset, _ := json.Marshal(result.PSetDefine)

	r.ResultValue = string(s_value)
	r.PSetDefine = string(s_pset)

	_, err = service.DB.UpdateResult(r)
	if err != nil {
		fmt.Printf("HandleResult err:%s\n", err.Error())
		return nil
	}

	workorder, err := service.DB.GetWorkorder(result.Workorder_ID)
	if err == nil {
		// 结果推送hmi
		ws_result := payload.WSResult{}
		ws_result.Result_id = result.Result_id
		ws_result.Count = result.Count
		ws_result.Result = result.Result
		ws_result.MI = result.ResultValue.Mi
		ws_result.WI = result.ResultValue.Wi
		ws_result.TI = result.ResultValue.Ti
		ws_str, _ := json.Marshal(ws_result)
		go service.APIService.WSSendResult(workorder.HMISN, string(ws_str))
	}

	if r.Count >= int(workorder.MaxRedoTimes) || r.Result == payload.RESULT_OK {
		// 结果推送AIIS

		odoo_result := payload.ODOOResult{}
		odoo_result.Control_date = fmt.Sprintf("%sT%s+08:00", times[0], times[1])

		odoo_result.Measure_degree = result.ResultValue.Wi
		odoo_result.Measure_result = strings.ToLower(result.Result)
		odoo_result.Measure_t_don = result.ResultValue.Ti
		odoo_result.Measure_torque = result.ResultValue.Mi
		odoo_result.Op_time = result.Count
		odoo_result.Pset_m_max = result.PSetDefine.Mp
		odoo_result.Pset_m_min = result.PSetDefine.Mm
		odoo_result.Pset_m_target = result.PSetDefine.Ma
		odoo_result.Pset_m_threshold = result.PSetDefine.Ms
		odoo_result.Pset_strategy = result.PSetDefine.Strategy
		odoo_result.Pset_w_max = result.PSetDefine.Wp
		odoo_result.Pset_w_min = result.PSetDefine.Wm
		odoo_result.Pset_w_target = result.PSetDefine.Wa

		curves, err := service.DB.ListCurves(result.Result_id)
		if err != nil {
			return err
		}
		for _, v := range curves {
			curobject := payload.CURObject{}
			curobject.OP = v.Count
			curobject.File = v.CurveFile
			odoo_result.CURObjects = append(odoo_result.CURObjects, curobject)
		}

		_, err = service.ODOO.PutResult(result.Result_id, odoo_result)
		if err == nil {
			// 发送成功
			r.HasUpload = true
			_, err := service.DB.UpdateResult(r)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (service *CVI3Service) HandleCurve(curve payload.ControllerCurve) (error) {
	// 保存波形到数据库
	c := rushdb.Curves{}
	c.ResultID = curve.ResultID
	c.CurveData = curve.CurveData
	c.CurveFile = curve.CurveFile
	c.Count = curve.Count
	c.HasUpload = false

	exist, err := service.DB.CurveExist(c)
	if err != nil {
		return err
	} else {
		if exist {
			_, err := service.DB.UpdateCurve(c)
			if err != nil {
				return err
			}
		} else {
			err := service.DB.InsertCurve(c)
			if err != nil {
				return err
			}
		}
	}

	// 保存波形到对象存储
	err = service.Storage.Upload(curve.CurveFile, curve.CurveData)
	if err != nil {
		return err
	} else {
		c.HasUpload = true
		_, err = service.DB.UpdateCurve(c)
		if err != nil {
			return err
		}
	}

	return nil

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

		// 结果数据
		result_data := payload.XML2Result(result)

		// 波形文件
		curve_file := payload.XML2Curve(result)


		curve := payload.ControllerCurve{}
		s_curvedata, _ := json.Marshal(curve_file)
		curve.CurveData = string(s_curvedata)
		curve.Count = result_data.Count
		curve.CurveFile = result_data.CurFile
		curve.ResultID = result_data.Result_id

		e := service.HandleCurve(curve)
		if  e == nil {
			go service.HandleResult(result_data)
		} else {
			fmt.Printf("OnRecv err:%s\n", e.Error())
		}

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