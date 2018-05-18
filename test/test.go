package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"strconv"
	"sync"
	"bytes"
	"encoding/csv"
	"os"
)

type PR struct {
	Pr_value string `json:"pr_value"`
	Pr_group string `json:"pr_group"`
}

type MO struct {
	Pin_check_code int	`json:"pin_check_code"`
	Equipment_name string	`json:"equipment_name"`
	Factory_name string	`json:"factory_name"`
	Pin int	`json:"pin"`
	Year int `json:"year"`
	Assembly_line string `json:"assembly_line"`
	Model string `json:"model"`
	Vin string `json:"vin"`
	Lnr string `json:"lnr"`
	Date_planned_start string `json:"date_planned_start"`
	Prs []PR `json:"prs"`
}

func create_mo(odoo_url string, mo MO, mtx sync.Mutex) (error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/mrp.productions", odoo_url)
	body, _ := json.Marshal(mo)
	s_body := string(body)
	req, err := http.NewRequest("POST", url, strings.NewReader(s_body))

	//SR1J--V001--C6-2018-6171427=5
	mo_name := fmt.Sprint("%s--V001--%s-%d-%d=%d", mo.Equipment_name, mo.Factory_name, mo.Year, mo.Pin, mo.Pin_check_code)

	if err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		t1 := time.Now()
		resp, _ := client.Do(req)
		elapsed := time.Since(t1)

		result := fmt.Sprintf("创建工单:%s 结果:%d 用时:%d\n", mo_name, resp.StatusCode, elapsed)
		fmt.Printf(result)

		// 写入文件
		s := make([]string,3)
		s[0] = mo_name
		s[1] = string(resp.StatusCode)
		s[2] = string(elapsed)
		go WriteFile(s, mtx)
	}

	return nil
}

func WriteFile(s []string, mtx sync.Mutex) {
	defer mtx.Unlock()

	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	r2.Write(s)
	r2.Flush()

	mtx.Lock()
	fout,err := os.Create("test.csv")
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer fout.Close()


	fout.WriteString(buf.String())
}

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}

func main() {

	FILE_MTX := sync.Mutex{}
	odoo_url := "http://10.1.1.31:8069"
	//odoo_url := "http://127.0.0.1:8069"
	count := 0

	mo := MO{}
	mo.Pin_check_code = 5
	mo.Year = 2018
	mo.Factory_name = "C6"
	mo.Equipment_name = "SR1J"
	mo.Assembly_line = "01"
	mo.Model = "BR24J3"
	mo.Lnr = "0001"


	s_date, s_time := GetDateTime()
	mo.Date_planned_start = fmt.Sprintf("%sT%s+08:00", s_date, s_time)

	for {
		if count == 1200 {
			s_date, s_time := GetDateTime()
			mo.Date_planned_start = fmt.Sprintf("%sT%s+08:00", s_date, s_time)
			count = 0
		}

		s_pin := fmt.Sprintf("%07d", count)
		mo.Pin, _ = strconv.Atoi(s_pin)
		mo.Vin = fmt.Sprintf("LSV2A8CA7JN%06d", count)

		create_mo(odoo_url, mo, FILE_MTX)

		count++
	}
}