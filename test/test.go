package main

// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./test.go

import (
	"fmt"
	"strings"
	"time"
	"sync"
	"os"
	"strconv"
	"flag"
	"github.com/masami10/rush/payload"
	"github.com/masami10/rush/core"
	"math/rand"
)

var g_Count int = 0
var g_TableCount int = 0
var g_odoo_url string = "http://10.1.1.31:8069"
var g_task_num int = 8
//var g_req_mtx sync.Mutex
var g_req_inteval int = 200
var g_genetime_mtx sync.Mutex = sync.Mutex{}

var g_odoo core.ODOO = core.ODOO{}

var g_time string = ""

//var g_result_value = [2]string{"ok", "nok"}
//
//var g_r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func create_mo(odoo_url string, mo payload.ODOOMO, odoo_result payload.ODOOResult, mtx sync.Mutex, time_mtx sync.Mutex) (error) {
	mo.Date_planned_start = g_time
	mo.Pin = GenerateRangeNum(7, time_mtx)
	vin_rand := GenerateRangeNum(6, time_mtx)
	mo.Vin = fmt.Sprintf("LSV2A8CA7JN%d", vin_rand)

	//SR1J--V001--C6-2018-6171427=5
	mo_name := fmt.Sprintf("%s--V001--%s-%d-%d=%d", mo.Equipment_name, mo.Factory_name, mo.Year, mo.Pin, mo.Pin_check_code)


	_, t, err := g_odoo.CreateMO(mo)


	if err != nil {
		result := fmt.Sprintf("创建工单:%s 用时:%s 结果:%s \n", mo_name, t, "failed")
		fmt.Printf(result)
	} else {
		// 写入文件
		go tracefile(t, mtx)

		// 推送结果
		//if len(created)  > 0 {
		//	go put_result(created[0], odoo_result, time_mtx)
		//}


		result := fmt.Sprintf("创建工单:%s 用时:%s 结果:%s \n", mo_name, t, "ok")
		fmt.Printf(result)
	}

	return nil
}

func put_result(mo payload.ODOOMOCreated, result payload.ODOOResult, mtx sync.Mutex) {
	var r string
	for _, v := range mo.Result_IDs {
		result = RandomResult(result, mtx)

		t, err := g_odoo.PutResult(v, result)

		if err != nil {
			r = fmt.Sprintf("推送结果:%d 用时:%s 结果:%s \n", v, t, "fail")
		} else {
			r = fmt.Sprintf("推送结果:%d 用时:%s 结果:%s \n", v, t, "ok")
		}
		fmt.Printf(r)
	}
}

func RandomResult(result payload.ODOOResult, mtx sync.Mutex) (payload.ODOOResult) {
	result.Control_date = GenerateTime()
	result.Measure_degree = (rand.Float64()*180) + 5

	s := GenerateRangeNum(8, mtx)
	if s % 2 == 0{
		result.Measure_result = "ok"
		result.Op_time = 1
	} else {
		result.Measure_result = "nok"
		result.Op_time = 2
	}

	if string(s)[0] == '3' {
		result.Op_time = 3
	}

	//fmt.Printf("%s\n", result.Measure_result)
	result.Measure_t_don = rand.Float64() * 100 + 100
	result.Measure_torque = (rand.Float64()*10) + 10

	return result
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

	t := fmt.Sprintf("%09d",nano)
	//fmt.Printf("%s\n", t)
	s := t[0:len]
	i,_ := strconv.Atoi(s)
	return i
}

func tracefile(str_content string, mtx sync.Mutex)  {
	defer mtx.Unlock()
	mtx.Lock()
	fd,_:=os.OpenFile("test.csv",os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	defer fd.Close()
	fd_content:=strings.Join([]string{str_content,"\n"},"")
	buf:=[]byte(fd_content)
	fd.Write(buf)
}

func ReCount(mo payload.ODOOMO, mtx sync.Mutex) payload.ODOOMO {
	defer mtx.Unlock()

	mtx.Lock()

	if g_Count >= 1200 {
		g_TableCount++
		s_date, s_time := GetDateTime()
		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", s_date, s_time))
		_t := t.Add(time.Duration(744 * g_TableCount) * time.Hour)
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

func RunTask_MO(mo payload.ODOOMO, odoo_result payload.ODOOResult, count_mtx sync.Mutex, file_mtx sync.Mutex, time_mtx sync.Mutex) {
	for {
		ReCount(mo, count_mtx)

		create_mo(g_odoo_url, mo, odoo_result, file_mtx, time_mtx)

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

	odoo := flag.String("odoo", "http://10.1.1.31", "--odoo")
	//odoo := flag.String("odoo", "http://127.0.0.1:8069", "--odoo")
	task_num := flag.Int("task", 4, "--task")
	req_itv := flag.Int("inteval", 100, "--inteval")
	flag.Parse()

	g_odoo_url = *odoo
	g_task_num = *task_num
	g_req_inteval = *req_itv

	fmt.Printf("odoo:%s, task:%d, inteval:%d\n", g_odoo_url, g_task_num, g_req_inteval)

	g_odoo.MaxRetry = 3
	g_odoo.URL = g_odoo_url

	//return
	FILE_MTX := sync.Mutex{}
	TIME_MTX := sync.Mutex{}
	COUNT_MTX := sync.Mutex{}


	odoo_result := payload.ODOOResult{}
	//odoo_result.Control_date = result_data.Dat
	odoo_result.CURObjects = []payload.CURObject{}
	cur_object := payload.CURObject{}
	cur_object.File = "test.json"
	cur_object.OP = 1
	odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)

	//odoo_result.Op_time = 1
	odoo_result.Pset_m_max = 0.5
	odoo_result.Pset_m_min = 0.01
	odoo_result.Pset_m_target = 0.1
	odoo_result.Pset_m_threshold = 0.2
	odoo_result.Pset_strategy = "AW"
	odoo_result.Pset_w_max = 190
	odoo_result.Pset_w_min = 170
	odoo_result.Pset_w_target = 180


	mo := payload.ODOOMO{}
	mo.Pin_check_code = 5
	mo.Year = 2018
	mo.Factory_name = "C6"
	mo.Equipment_name = "SR1J"
	mo.Assembly_line = "01"
	mo.Model = "BR24J3"
	//mo.Model2 = "model"
	mo.Lnr = "0001"
	mo.Prs = []payload.ODOOPR{}
	var pr payload.ODOOPR = payload.ODOOPR{}
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
		go RunTask_MO(mo, odoo_result, COUNT_MTX, FILE_MTX, TIME_MTX)

		//if g_test_result {
		//	go RunTask_Result(1, odoo_result, FILE_MTX)
		//}

	}

	for {
		time.Sleep(5 * time.Second)
	}
}