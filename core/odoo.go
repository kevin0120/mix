package core

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"github.com/masami10/rush/db"
	"time"
	"github.com/masami10/rush/payload"
	"strings"
)

type ODOO struct {
	URL string
	MasterPC_SN string
	DB *rushdb.DB
	APIService *APIServer
}

func (odoo *ODOO) ListWorkorders(masterpc_sn string, limit int) ([]payload.ODOOWorkorder, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/mrp.workorders?masterpc=%s&limit=%d", odoo.URL, masterpc_sn, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	} else {
		//req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {

			body, _ := ioutil.ReadAll(resp.Body)
			//ss := string(body)
			//fmt.Printf("%s\n", ss)
			var workorders []payload.ODOOWorkorder
			err := json.Unmarshal(body, &workorders)

			if err != nil {
				return nil, err
			} else {
				return workorders, nil
			}
		} else {
			return nil, errors.New(string(resp.StatusCode))
		}
	}
}

func (odoo *ODOO) PutResult(id int, result payload.ODOOResult) (error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/operation.results/%d", odoo.URL, id)
	body, _ := json.Marshal(result)
	s_body := string(body)
	req, err := http.NewRequest("PUT", url, strings.NewReader(s_body))

	if err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		fmt.Printf("推送odoo结果:%d\n", resp.StatusCode)
	}

	return nil
}

func (odoo *ODOO) RequestWorkerOrders() {
	for {
		workorders, err := odoo.ListWorkorders(odoo.MasterPC_SN, 10)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		neworders, e := odoo.DB.InsertWorkorders(workorders)
		if e != nil {
			fmt.Printf("%s\n", e.Error())
		}

		// 推送新工单
		for _,v := range neworders {
			order_str, _ := json.Marshal(v)
			odoo.APIService.WSSendWorkorder(v.HMI.UUID, string(order_str))
		}

		time.Sleep(5 * time.Second)
	}
}
