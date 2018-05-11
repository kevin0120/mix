package cvi_listener

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"io"
	"net/url"
)

const (
	API_PREFIX = "/api/v1"
)

type ApiServer struct {
	cvi3_manager *CVI3Manager
}

func (as *ApiServer) StartService(manager *CVI3Manager) {
	as.cvi3_manager = manager

	m := mux.NewRouter()

	m.HandleFunc(API_PREFIX + "/api-doc", as.APIDoc).Methods("GET")
	m.HandleFunc(API_PREFIX + "/psets", as.PSets).Methods("PUT", "OPTIONS")
	m.HandleFunc(API_PREFIX + "/status", as.Status).Methods("GET")
	m.HandleFunc(API_PREFIX + "/results", as.Results).Methods("GET")

	err := http.ListenAndServe(":8080", m) //设置监听的端口
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func (as *ApiServer)APIDoc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "wefawef")
}

// 设置PSET拧接程序
func (as *ApiServer)PSets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	pset := RequestPSet{}
	err := json.Unmarshal(body, &pset)
	if err != nil {
		// 参数错误
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
	}

	e := as.cvi3_manager.PSet(pset.SN, pset.PSet, pset.Workorder_ID)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)

		rt, _ := json.Marshal(ErrorMsg{e.Error()})
		io.WriteString(w, string(rt))
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

// 取得控制器状态
func (as *ApiServer)Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	u, _ := url.Parse(r.RequestURI)
	q, _ := url.ParseQuery(u.RawQuery)
	sn, exist := q["sn"]

	rts := []ResponseStatus{}

	if exist {
		// 指定控制器的状态
		s := as.cvi3_manager.CVI3_clients[sn[0]].Status
		rts = append(rts, ResponseStatus{sn[0], s})
	} else {
		// 所有状态
		for k, v := range as.cvi3_manager.CVI3_clients {
			rts = append(rts, ResponseStatus{k, v.Status})
		}
	}

	rt, _ := json.Marshal(rts)
	io.WriteString(w, string(rt))
}

// 取得控制器中的结果
func (as *ApiServer)Results(w http.ResponseWriter, r *http.Request) {

}

type ErrorMsg struct {
	Msg string	`json:"msg"`
}

type RequestPSet struct {
	SN string `json:"sn"`
	PSet int `json:"pset"`
	Workorder_ID int `json:"workorder_id"`
}

type ResponseStatus struct {
	SN string `json:"sn"`
	Status string `json:"status"`
}

type ResponseResult struct {
	SN string `json:"sn"`
	Workorder_ID int `json:"workorder_id"`
	CurFile string `json:"cur_file"`
	Result string `json:"result"`
	Dat string `json:"dat"`
	PSet int `json:"pset"`
	PSetDefine struct {
		Strategy string `json:"strategy"`
		Mp float64 `json:"M+"`
		Mm float64 `json:"M-"`
		Ms float64 `json:"MS"`
		Ma float64 `json:"MA"`
		Wp float64 `json:"W+"`
		Wm float64 `json:"W-"`
		Wa float64 `json:"WS"`
	} `json:"pset_define"`

	ResultValue struct {
		Mi float64 `json:"MI"`
		Wi float64 `json:"WI"`
		Ti float64 `json:"TI"`
	} `json:"result_value"`
}