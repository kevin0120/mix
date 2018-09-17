package main

// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./test.go

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/storage"
	"gopkg.in/resty.v1"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var g_Count int = 0
var g_TableCount int = 0
var g_odoo_url string = "http://10.1.1.31:8069"
var g_task_num int = 8
var g_result_batch bool = false

//var g_req_mtx sync.Mutex
var g_req_inteval int = 200
var g_genetime_mtx sync.Mutex = sync.Mutex{}

var g_odoo odoo.Service = odoo.Service{}

var g_time string = ""

var g_odoo_result odoo.ODOOResultSync = odoo.ODOOResultSync{}
var g_odoo_result_batch odoo.ODOOResultSync = odoo.ODOOResultSync{}

var g_httpClient_odoo *resty.Client
var g_httpClient_aiis *resty.Client

//var g_result_value = [2]string{"ok", "nok"}
//
//var g_r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ODOOCreateMO(url string, body interface{}) (odoo.ODOOMOCreated, error) {

	created := odoo.ODOOMOCreated{}

	r := g_httpClient_odoo.R().SetBody(body)

	resp, err := r.Post(url)
	if err != nil {
		return created, fmt.Errorf("Create MO Post fail: %s\n", err)
	} else {
		if resp.StatusCode() != http.StatusCreated {
			return created, fmt.Errorf("Create MO Post fail: %s\n", resp.Status())
		} else {
			json.Unmarshal(resp.Body(), &created)
		}
	}

	return created, nil
}

func create_mo(odoo_url string, mo odoo.ODOOMO, mtx sync.Mutex, time_mtx sync.Mutex) error {
	mo.Date_planned_start = g_time
	mo.Pin = GenerateRangeNum(7, time_mtx)
	vin_rand := GenerateRangeNum(6, time_mtx)
	mo.Vin = fmt.Sprintf("LSVCDDCA7JN%d", vin_rand)

	//SR1J--V001--C6-2018-6171427=5
	mo_name := fmt.Sprintf("%s--V001--%s-%d-%d=%d", mo.Equipment_name, mo.Factory_name, mo.Year, mo.Pin, mo.Pin_check_code)
	fmt.Printf("%s\n", mo_name)

	created, err := ODOOCreateMO(fmt.Sprintf("%s/api/v1/api/v1/mrp.productions", g_odoo_url), mo)

	if err == nil {
		go put_result(created, time_mtx)
	}

	return nil
}

func putResult(body interface{}, url string) error {
	r := g_httpClient_aiis.R().SetBody(body)
	var resp *resty.Response
	var err error

	resp, err = r.Put(url)

	if err != nil {
		return fmt.Errorf("Result Put fail: %s ", err.Error())
	} else {
		if resp.StatusCode() != http.StatusNoContent {
			return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
		}
	}
	return nil
}

func resultToAiis(odoo_result *odoo.ODOOResultSync) aiis.AIISResult {
	aiisResult := aiis.AIISResult{}

	if odoo_result.Measure_result == storage.RESULT_OK {
		aiisResult.Final_pass = controller.ODOO_RESULT_PASS
		aiisResult.One_time_pass = controller.ODOO_RESULT_PASS

		aiisResult.QualityState = controller.QUALITY_STATE_PASS
		aiisResult.ExceptionReason = ""

	} else {
		aiisResult.Final_pass = controller.ODOO_RESULT_FAIL
		aiisResult.One_time_pass = controller.ODOO_RESULT_FAIL
		aiisResult.QualityState = controller.QUALITY_STATE_PASS
		aiisResult.ExceptionReason = ""
	}

	aiisResult.ExceptionReason = odoo_result.ExceptionReason

	//aiisResult.Control_date = odoo_result.UpdateTime.Format(time.RFC3339)

	aiisResult.Measure_degree = odoo_result.Measure_degree
	aiisResult.Measure_result = strings.ToLower(odoo_result.Measure_result)
	aiisResult.Measure_t_don = odoo_result.Measure_t_don
	aiisResult.Measure_torque = odoo_result.Measure_torque
	//aiisResult.Op_time = result.Count
	//aiisResult.Pset_m_max = result.PSetDefine.Mp
	//aiisResult.Pset_m_min = result.PSetDefine.Mm
	//aiisResult.Pset_m_target = result.PSetDefine.Ma
	//aiisResult.Pset_m_threshold = result.PSetDefine.Ms
	//aiisResult.Pset_strategy = result.PSetDefine.Strategy
	//aiisResult.Pset_w_max = result.PSetDefine.Wp
	//aiisResult.Pset_w_min = result.PSetDefine.Wm
	//aiisResult.Pset_w_target = result.PSetDefine.Wa
	//aiisResult.Pset_w_threshold = 1
	aiisResult.UserID = 1

	// mo相关
	//aiisResult.MO_AssemblyLine = workorder.MO_AssemblyLine
	//aiisResult.MO_EquipemntName = workorder.MO_EquipemntName
	//aiisResult.MO_FactoryName = workorder.MO_FactoryName
	//aiisResult.MO_Pin = workorder.MO_Pin
	//aiisResult.MO_Pin_check_code = workorder.MO_Pin_check_code
	//aiisResult.MO_Year = workorder.MO_Year
	//aiisResult.MO_Lnr = workorder.MO_Lnr
	//aiisResult.MO_NutNo = r.NutNo
	//aiisResult.MO_Model = workorder.MO_Model

	return aiisResult
}

func put_result(mo odoo.ODOOMOCreated, mtx sync.Mutex) {
	var r string

	for _, v := range mo.Result_IDs {
		RandomResult(mtx)

		err := putResult(resultToAiis(&g_odoo_result), fmt.Sprint("http://127.0.0.1:9092/operation.results/%d", v))

		if err != nil {
			r = fmt.Sprintf("推送结果:%d 结果:%s \n", v, "fail")
		} else {
			r = fmt.Sprintf("推送结果:%d 结果:%s \n", v, "ok")
		}
		fmt.Printf(r)
	}

}

func RandomResult(mtx sync.Mutex) {
	if g_result_batch == true {
		g_odoo_result_batch.Control_date = GenerateTime()
		g_odoo_result_batch.Measure_degree = (rand.Float64() * 180) + 5

		s := GenerateRangeNum(8, mtx)
		if s%2 == 0 {
			g_odoo_result_batch.Measure_result = "ok"
		} else {
			g_odoo_result_batch.Measure_result = "nok"
		}

		g_odoo_result_batch.Op_time = rand.Intn(3) + 1

		//fmt.Printf("%s\n", result.Measure_result)
		g_odoo_result_batch.Measure_t_don = rand.Float64()*100 + 100
		g_odoo_result_batch.Measure_torque = (rand.Float64() * 10) + 10
	} else {
		g_odoo_result.Control_date = GenerateTime()
		g_odoo_result.Measure_degree = (rand.Float64() * 180) + 5

		s := GenerateRangeNum(8, mtx)
		if s%2 == 0 {
			g_odoo_result.Measure_result = "ok"
		} else {
			g_odoo_result.Measure_result = "nok"
		}

		g_odoo_result.Op_time = rand.Intn(3) + 1

		//fmt.Printf("%s\n", result.Measure_result)
		g_odoo_result.Measure_t_don = rand.Float64()*100 + 100
		g_odoo_result.Measure_torque = (rand.Float64() * 10) + 10
	}

}

func GenerateTime() string {
	defer g_genetime_mtx.Unlock()

	g_genetime_mtx.Lock()
	sdate, stime := GetDateTime()
	return fmt.Sprintf("%sT%s+08:00", sdate, stime)
}

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}

func GenerateRangeNum(len int, mtx sync.Mutex) int {
	defer mtx.Unlock()

	mtx.Lock()
	nano := time.Now().Nanosecond()

	t := fmt.Sprintf("%09d", nano)
	//fmt.Printf("%s\n", t)
	s := t[0:len]
	i, _ := strconv.Atoi(s)
	return i
}

func tracefile(str_content string, mtx sync.Mutex) {
	defer mtx.Unlock()
	mtx.Lock()
	fd, _ := os.OpenFile("test.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	fd_content := strings.Join([]string{str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
}

func ReCount(mo odoo.ODOOMO, mtx sync.Mutex) odoo.ODOOMO {
	defer mtx.Unlock()

	mtx.Lock()

	if g_Count >= 1200 {
		g_TableCount++
		s_date, s_time := GetDateTime()
		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", s_date, s_time))
		_t := t.Add(time.Duration(744*g_TableCount) * time.Hour)
		st := strings.Split(_t.Format("2006-01-02 15:04:05"), " ")
		g_time = fmt.Sprintf("%sT%s+08:00", st[0], st[1])

		g_Count = 0
	}

	return mo
}

func AddCount(mtx sync.Mutex) {
	defer mtx.Unlock()

	mtx.Lock()
	g_Count++
}

func RunTask_MO(mo odoo.ODOOMO, count_mtx sync.Mutex, file_mtx sync.Mutex, time_mtx sync.Mutex) {
	for {
		ReCount(mo, count_mtx)

		create_mo(g_odoo_url, mo, file_mtx, time_mtx)

		AddCount(count_mtx)

		time.Sleep(time.Duration(g_req_inteval) * time.Millisecond)
	}
}

//func RunTask_Result(result_id int, result payload.ODOOResult, file_mtx sync.Mutex) {
//	for {
//		put_result(result_id, result, g_odoo_url, file_mtx)
//		time.Sleep(time.Duration(g_req_inteval) * time.Millisecond)
//	}
//}

//--odoo --task --inteval
func main() {

	fmt.Printf("start\n")

	//odoo := flag.String("odoo", "http://10.1.1.31", "--odoo")
	odoo_url := flag.String("odoo", "http://127.0.0.1:8069", "--odoo")
	task_num := flag.Int("task", 4, "--task")
	req_itv := flag.Int("inteval", 100, "--inteval")
	batch := flag.Bool("batch", true, "--batch")
	flag.Parse()

	g_odoo_url = *odoo_url
	g_task_num = *task_num
	g_req_inteval = *req_itv
	g_result_batch = *batch

	fmt.Printf("odoo:%s, task:%d, inteval:%d\n", g_odoo_url, g_task_num, g_req_inteval)

	//g_odoo.Conf.MaxRetry = 3
	//g_odoo.Conf.Urls = []string{}
	//g_odoo.Conf.Urls = append(g_odoo.Conf.Urls, g_odoo_url)
	//g_odoo.Conf.Urls[0] = g_odoo_url

	g_httpClient_odoo = resty.New()
	g_httpClient_odoo.SetRESTMode() // restful mode is default
	g_httpClient_odoo.SetTimeout(time.Duration(5 * time.Second))
	g_httpClient_odoo.SetContentLength(true)
	// Headers for all request
	g_httpClient_odoo.
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(1 * time.Second)).
		SetRetryMaxWaitTime(20 * time.Second)

	g_httpClient_aiis = resty.New()
	g_httpClient_aiis.SetRESTMode() // restful mode is default
	g_httpClient_aiis.SetTimeout(time.Duration(5 * time.Second))
	g_httpClient_aiis.SetContentLength(true)
	// Headers for all request
	g_httpClient_aiis.
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(1 * time.Second)).
		SetRetryMaxWaitTime(20 * time.Second)

	//return
	FILE_MTX := sync.Mutex{}
	TIME_MTX := sync.Mutex{}
	COUNT_MTX := sync.Mutex{}

	//odoo_result.Control_date = result_data.Dat

	if g_result_batch == true {
		g_odoo_result_batch.CURObjects = []aiis.CURObject{}
		cur_object := aiis.CURObject{}
		cur_object.File = "test.json"
		cur_object.OP = 1
		g_odoo_result_batch.CURObjects = append(g_odoo_result.CURObjects, cur_object)

		//odoo_result.Op_time = 1
		g_odoo_result_batch.Pset_m_max = 0.5
		g_odoo_result_batch.Pset_m_min = 0.01
		g_odoo_result_batch.Pset_m_target = 0.1
		g_odoo_result_batch.Pset_m_threshold = 0.2
		g_odoo_result_batch.Pset_strategy = "AW"
		g_odoo_result_batch.Pset_w_max = 190
		g_odoo_result_batch.Pset_w_min = 170
		g_odoo_result_batch.Pset_w_target = 180
	} else {
		g_odoo_result.CURObjects = []aiis.CURObject{}
		cur_object := aiis.CURObject{}
		cur_object.File = "test.json"
		cur_object.OP = 1
		g_odoo_result.CURObjects = append(g_odoo_result.CURObjects, cur_object)

		//odoo_result.Op_time = 1
		g_odoo_result.Pset_m_max = 0.5
		g_odoo_result.Pset_m_min = 0.01
		g_odoo_result.Pset_m_target = 0.1
		g_odoo_result.Pset_m_threshold = 0.2
		g_odoo_result.Pset_strategy = "AW"
		g_odoo_result.Pset_w_max = 190
		g_odoo_result.Pset_w_min = 170
		g_odoo_result.Pset_w_target = 180
	}

	mo := odoo.ODOOMO{}
	mo.Pin_check_code = 5
	mo.Year = 2018
	mo.Factory_name = "C6"
	mo.Equipment_name = "SR1J"
	mo.Assembly_line = "01"
	mo.Model = "BR24J3"
	//mo.Model2 = "model"
	mo.Lnr = "0001"
	mo.Prs = []odoo.ODOOPR{}
	var pr odoo.ODOOPR = odoo.ODOOPR{}
	pr.Pr_group = "GSP"
	pr.Pr_value = "G0C"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "ASW"
	pr.Pr_value = "4x0"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "RAD"
	pr.Pr_value = "C0D"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "REI"
	pr.Pr_value = "HSS"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "BRS"
	pr.Pr_value = "1AC"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "GMO"
	pr.Pr_value = "TJ5"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "EDF"
	pr.Pr_value = "0ER"
	mo.Prs = append(mo.Prs, pr)

	pr.Pr_group = "AED"
	pr.Pr_value = "7JI"
	mo.Prs = append(mo.Prs, pr)

	s_date, s_time := GetDateTime()
	g_time = fmt.Sprintf("%sT%s+08:00", s_date, s_time)

	// 启动任务
	for i := 0; i < g_task_num; i++ {
		go RunTask_MO(mo, COUNT_MTX, FILE_MTX, TIME_MTX)
	}

	for {
		time.Sleep(10000 * time.Millisecond)
	}
}
