package cvi_listener

import (
	"net/http"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/masami10/rush/utils/cvi"
)

type HMI struct {
	URL string
}

func XMl2Result(result cvi.CVI3Result) (ResponseResult) {
	rr := ResponseResult{}

	rr.SN = result.PRC_SST.PAR.SN
	rr.Result = result.PRC_SST.PAR.Result
	if rr.Result == "IO" {
		rr.Result = "OK"
	} else if rr.Result == "NIO" {
		rr.Result = "NOK"
	}

	rr.PSet = result.PRC_SST.PAR.FAS.GRP.TIP.PSet
	rr.Workorder_ID = result.PRC_SST.PAR.Workorder_id
	screw_id := result.PRC_SST.PAR.Screw_id
	rr.CurFile = fmt.Sprintf("%s_%d_%s.json", rr.SN, rr.Workorder_ID, screw_id)
	rr.Dat = fmt.Sprintf("%sT%s+08:00", result.PRC_SST.PAR.FAS.GRP.TIP.Date, result.PRC_SST.PAR.FAS.GRP.TIP.Time)
	rr.PSetDefine.Strategy = result.PRC_SST.PAR.FAS.GRP.TIP.BLC.PRO.Strategy

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
		case "TI":
			if result_values[i].Unit == "s" {
				rr.ResultValue.Ti = result_values[i].Value * 1000
			} else {
				rr.ResultValue.Ti = result_values[i].Value
			}
		}
	}

	return rr
}
func (hmi *HMI) PushResult(result ResponseResult) (error) {

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/results", hmi.URL)
	body, _ := json.Marshal(result)
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		client.Do(req)
	}

	return nil
}

func (hmi *HMI) PushStauts(status ResponseStatus) (error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/status", hmi.URL)
	body, _ := json.Marshal(status)
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		client.Do(req)
	}

	return nil
}