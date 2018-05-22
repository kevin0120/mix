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
	"strconv"
)

const (
	ODOO_ERR_KEY = "Odoo Server Error"
)

type ODOOConf struct {
	Timeout int		`yaml:"api_timeout"`
	Urls []string	`yaml:"urls"`
}

type ODOO struct {
	URL string
	MaxRetry int
	MasterPC_SN string
	DB *rushdb.DB
	APIService *APIServer
}

func (odoo *ODOO) CreateMO(mo payload.ODOOMO) ([]payload.ODOOMOCreated, string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/mrp.productions", odoo.URL)

	var created []payload.ODOOMOCreated = []payload.ODOOMOCreated{}
	body, em := json.Marshal(mo)
	var final_err error

	if em != nil {
		fmt.Printf("%s\n", em.Error())
		return created, "", em
	}

	s_body := string(body)
	req, err := http.NewRequest("POST", url, strings.NewReader(s_body))
	//mo_name := fmt.Sprintf("%s--V001--%s-%d-%d=%d", mo.Equipment_name, mo.Factory_name, mo.Year, mo.Pin, mo.Pin_check_code)

	if err != nil {
		return created, "", err
	} else {
		req.Header.Set("Content-Type", "application/json")

		for i := 0; i < odoo.MaxRetry; i++ {

			t := time.Now()
			resp, err := client.Do(req)
			elapsed := strconv.FormatFloat(time.Since(t).Seconds(), 'f', -1, 64)

			if err != nil {
				final_err = err
				continue
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusCreated {
				final_err = err
				continue
			} else {
				//fmt.Printf("%s\n", string(body))
				err := json.Unmarshal(body, &created)
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					return created, elapsed, err
				} else {
					return created, elapsed, nil
				}

			}

			time.Sleep(time.Duration(1 + i) * time.Second)
		}

	}

	return created, "", final_err
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

func (odoo *ODOO) GetWorkorder(order_id int) (payload.ODOOWorkorder, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/mrp.workorders/%d", odoo.URL, order_id)

	req, err := http.NewRequest("GET", url, nil)
	var workorder payload.ODOOWorkorder = payload.ODOOWorkorder{}
	var final_err error
	if err != nil {
		return workorder, err
	} else {

		// 循环请求
		for i := 0; i < odoo.MaxRetry; i++ {
			resp, err := client.Do(req)
			if err != nil {
				final_err = err
				continue
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK && !strings.Contains(string(body), ODOO_ERR_KEY) {

				err := json.Unmarshal(body, &workorder)

				if err != nil {
					final_err = err
					continue
				} else {
					return workorder, nil
				}
			} else {
				final_err = errors.New(string(resp.StatusCode))
				continue
			}

			time.Sleep(time.Duration(1 + i) * time.Second)
		}

		return workorder, final_err
	}
}

func (odoo *ODOO) PutResult(id int, result payload.ODOOResult) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/operation.results/%d", odoo.URL, id)
	body, _ := json.Marshal(result)
	s_body := string(body)
	//fmt.Printf("odoo结果:%s\n", s_body)
	var final_err error
	req, err := http.NewRequest("PUT", url, strings.NewReader(s_body))

	if err != nil {
		return "", err
	} else {
		req.Header.Set("Content-Type", "application/json")

		// 循环请求
		for i := 0; i < odoo.MaxRetry; i++ {
			t := time.Now()
			resp, err := client.Do(req)
			elapsed := strconv.FormatFloat(time.Since(t).Seconds(), 'f', -1, 64)
			if err != nil {
				final_err = err
				time.Sleep(time.Duration(1 + i) * time.Second)
				continue
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			// 更新flag
			if resp.StatusCode == http.StatusOK && !strings.Contains(string(body), ODOO_ERR_KEY) {
				//odoo.DB.UpdateResults(id, result.Op_time, true)
				return elapsed, nil
			} else {
				//fmt.Printf("%s\n", string(body))
				final_err = errors.New(string(body))
				time.Sleep(time.Duration(1 + i) * time.Second)
				continue
			}


		}


	}

	return "", final_err
}



func (odoo *ODOO) PatchCurve(result_id int, cur_file string, count int) (error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/operation.results/%d/curves_add", odoo.URL, result_id)
	ca := payload.ODOOCurveAppend{}
	ca.OP = count
	ca.File = cur_file
	body, _ := json.Marshal(ca)
	s_body := string(body)

	fmt.Printf("odoo波形:%s\n", s_body)

	req, err := http.NewRequest("PATCH", url, strings.NewReader(s_body))

	if err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		fmt.Printf("追加odoo波形:%d\n", resp.StatusCode)
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
