package aiis

import (
	"sync/atomic"
	"time"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"strings"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
	PutResultDone()
}

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
}

func NewEndpoint(url string, headers map[string]string, method string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
	}
}

type OnOdooStatus func(status string)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	endpoints   []*Endpoint
	httpClient  *resty.Client
	rush_port   string
	DB          *storage.Service
	ws          utils.RecConn
	//ws			websocket.Client
	//LocalWSServer *wsnotify.Service
	//Odoo 		*odoo.Service
	OnOdooStatus OnOdooStatus
	SN           string
	rpc          GRPCClient
}

func NewService(c Config, d Diagnostic, rush_port string) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		endpoints: e,
		rush_port: rush_port,
		rpc:       GRPCClient{},
	}
	s.rpc.RPCRecv = s.OnRPCRecv
	s.rpc.srv = s
	s.configValue.Store(c)
	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(c.Timeout))
	client.SetContentLength(true)
	// Headers for all request
	client.SetHeaders(c.Headers)
	client.
		SetRetryCount(c.MaxRetry).
		SetRetryWaitTime(time.Duration(c.PushInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client

	//entry := strings.Split(s.Config().Urls[0], "://")[1]
	//url := url.URL{Scheme: "ws", Host: entry, Path: s.Config().WSResultRoute}
	//s.ws = utils.RecConn{}
	//s.ws.OnConnected = func() {
	//	ws_msg := WSMsg{
	//		Type: WS_REG,
	//		Data: WSRegist{
	//			Rush_SN: s.SN,
	//		},
	//	}
	//
	//	str, _ := json.Marshal(ws_msg)
	//	s.ws.WriteMessage(websocket.TextMessage, str)
	//}
	//
	//s.ws.Dial(url.String(), nil)

	go s.ResultUploadManager()

	s.rpc.Start()

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) OnRPCRecv(payload string) {
	rpcPayload := RPCPayload{}
	json.Unmarshal([]byte(payload), &rpcPayload)
	str_data, _ := json.Marshal(rpcPayload.Data)

	switch rpcPayload.Type {
	case TYPE_RESULT:
		rp := ResultPatch{}
		json.Unmarshal(str_data, &rp)
		err := s.DB.UpdateResultByCount(rp.ID, 0, rp.HasUpload)
		if err == nil {
			s.diag.Debug(fmt.Sprintf("结果上传成功 ID:%d", rp.ID))
		} else {
			s.diag.Error(fmt.Sprintf("结果上传失败 ID:%d", rp.ID), err)
		}
		break

	case TYPE_ODOO_STATUS:
		status := ODOOStatus{}
		json.Unmarshal(str_data, &status)

		if s.OnOdooStatus != nil {
			s.OnOdooStatus(status.Status)
		}
		break
	}

}

func (s *Service) PutResult(result_id int64, body interface{}) error {

	//var err error
	//for _, endpoint := range s.endpoints {
	//	err = s.putResult(body, fmt.Sprintf(endpoint.url, result_id), endpoint.method)
	//	if err == nil {
	//		// 如果第一次就成功，推出循环
	//		return nil
	//	}
	//}

	//ws_msg := WSMsg{
	//	Type: WS_RESULT,
	//	Data: WSOpResult{
	//		ResultID: result_id,
	//		Result:   body.(AIISResult),
	//		Port:     s.rush_port,
	//	},
	//}
	//
	//str, _ := json.Marshal(ws_msg)
	//s.ws.WriteMessage(websocket.TextMessage, str)
	//s.ws.SendText(string(str))

	result := WSOpResult{
		ResultID: result_id,
		Result:   body.(AIISResult),
		Port:     s.rush_port,
	}

	str, _ := json.Marshal(result)
	err := s.rpc.RPCSend(string(str))
	if err != nil {
		fmt.Printf("grpc err: %s\n", err.Error())
	}

	return err
}

func (s *Service) putResult(body interface{}, url string, method string) error {
	r := s.httpClient.R().SetBody(body).SetHeader("rush_port", s.rush_port)
	var resp *resty.Response
	var err error

	switch method {
	case "PATCH":
		resp, err = r.Patch(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "PUT":
		resp, err = r.Put(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "POST":
		resp, err = r.Post(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	default:
		return errors.New("Result Put :the Method is wrong")

	}
	s.diag.PutResultDone()
	return nil
}

func (s *Service) ResultToAiisResult(result *storage.Results, curve interface{}) (AIISResult, error) {
	aiisResult := AIISResult{}

	resultValue := ResultValue{}
	json.Unmarshal([]byte(result.ResultValue), &resultValue)

	psetDefine := PSetDefine{}
	json.Unmarshal([]byte(result.PSetDefine), &psetDefine)

	dbWorkorder, err := s.DB.GetWorkorder(result.WorkorderID, true)
	if err != nil {
		return aiisResult, err
	}

	if curve != nil {
		aiisResult.CURObjects = append(aiisResult.CURObjects, *curve.(*CURObject))
	}

	curves, err := s.DB.ListCurvesByResult(result.Id)
	if err == nil {
		aiisCurve := CURObject{}
		for _, v := range curves {
			if curve != nil {
				if curve.(*CURObject).File == v.CurveFile {
					continue
				}
			}

			aiisCurve.OP = v.Count
			aiisCurve.File = v.CurveFile
			aiisResult.CURObjects = append(aiisResult.CURObjects, aiisCurve)
		}
	}

	aiisResult.ID = result.Id
	aiisResult.Control_date = result.UpdateTime.Format(time.RFC3339)
	aiisResult.Measure_degree = resultValue.Wi
	aiisResult.Measure_result = strings.ToLower(result.Result)
	aiisResult.Measure_t_don = resultValue.Ti
	aiisResult.Measure_torque = resultValue.Mi
	aiisResult.Op_time = result.Count
	aiisResult.Pset_m_max = psetDefine.Mp
	aiisResult.Pset_m_min = psetDefine.Mm
	aiisResult.Pset_m_target = psetDefine.Ma
	aiisResult.Pset_m_threshold = psetDefine.Ms
	aiisResult.Pset_strategy = psetDefine.Strategy
	aiisResult.Pset_w_max = psetDefine.Wp
	aiisResult.Pset_w_min = psetDefine.Wm
	aiisResult.Pset_w_target = psetDefine.Wa
	aiisResult.Pset_w_threshold = 1
	aiisResult.UserID = result.UserID
	aiisResult.Seq = result.Seq

	aiisResult.MO_AssemblyLine = dbWorkorder.MO_AssemblyLine
	aiisResult.MO_EquipemntName = dbWorkorder.MO_EquipemntName
	aiisResult.MO_FactoryName = dbWorkorder.MO_FactoryName
	aiisResult.MO_Pin = dbWorkorder.MO_Pin
	aiisResult.MO_Pin_check_code = dbWorkorder.MO_Pin_check_code
	aiisResult.MO_Year = dbWorkorder.MO_Year
	aiisResult.MO_Lnr = dbWorkorder.MO_Lnr
	aiisResult.MO_NutNo = result.NutNo
	aiisResult.MO_Model = dbWorkorder.MO_Model
	aiisResult.Batch = result.Batch
	aiisResult.Vin = dbWorkorder.Vin

	aiisResult.Mode = dbWorkorder.Mode
	aiisResult.TighteningId, _ = strconv.ParseInt(result.TighteningID, 10, 64)
	aiisResult.Lacking = "normal"

	gun, err := s.DB.GetGun(result.GunSN)
	if err != nil {
		gun.GunID = 0
	}

	aiisResult.GunID = gun.GunID
	aiisResult.WorkcenterID = dbWorkorder.WorkcenterID
	aiisResult.ProductID = dbWorkorder.ProductID
	aiisResult.NutID = result.ConsuProductID

	aiisResult.WorkcenterCode = dbWorkorder.WorkcenterCode
	aiisResult.ToolSN = result.GunSN
	aiisResult.ControllerSN = result.ControllerSN

	aiisResult.Job = fmt.Sprintf("%d", dbWorkorder.JobID)

	if result.Result == storage.RESULT_OK {
		aiisResult.Final_pass = ODOO_RESULT_PASS
		if result.Count == 1 {
			aiisResult.One_time_pass = ODOO_RESULT_PASS
		} else {
			aiisResult.One_time_pass = ODOO_RESULT_FAIL
		}

		if (resultValue.Mi >= result.ToleranceMin && resultValue.Mi <= result.ToleranceMax) &&
			(resultValue.Wi >= result.ToleranceMinDegree && resultValue.Wi <= result.ToleranceMaxDegree) {
			aiisResult.QualityState = QUALITY_STATE_PASS
			aiisResult.ExceptionReason = ""
		} else {
			aiisResult.QualityState = QUALITY_STATE_EX
			aiisResult.ExceptionReason = QUALITY_STATE_EX
		}

	} else {
		aiisResult.Final_pass = ODOO_RESULT_FAIL
		aiisResult.One_time_pass = ODOO_RESULT_FAIL

		if (resultValue.Mi >= result.ToleranceMin && resultValue.Mi <= result.ToleranceMax) &&
			(resultValue.Wi >= result.ToleranceMinDegree && resultValue.Wi <= result.ToleranceMaxDegree) {

			aiisResult.QualityState = QUALITY_STATE_EX
			aiisResult.ExceptionReason = QUALITY_STATE_EX
		} else {
			aiisResult.QualityState = QUALITY_STATE_FAIL
			aiisResult.ExceptionReason = ""
		}
	}

	return aiisResult, nil
}

func (s *Service) ResultUploadManager() error {
	for {

		results, err := s.DB.ListUnuploadResults()
		if err == nil {
			for _, v := range results {
				aiisResult, err := s.ResultToAiisResult(&v, nil)
				if err == nil {
					s.PutResult(v.ResultId, aiisResult)
				}
			}
		}

		time.Sleep(time.Duration(s.Config().ResultUploadInteval))
	}
}
