package core

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/db"
	"github.com/masami10/rush/payload"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

const (
	API_PREFIX = "/api/v1"

	WS_EVENT_STATUS    = "status"
	WS_EVENT_RESULT    = "result"
	WS_EVENT_WORKORDER = "workorder"
)

type WSClient struct {
	ID   string
	Conn websocket.Connection
}

type APIServer struct {
	Port      string
	CVI3      *CVI3Service
	DB        *rushdb.DB
	WSClients map[string]WSClient
	WSMtx     sync.Mutex
	DocPath   string
}

func (apiserver *APIServer) putPSets(ctx iris.Context) {

	var err error
	var pset payload.PSet
	err = ctx.ReadJSON(&pset)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	if pset.Controller_SN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if pset.PSet == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("pset is required")
		return
	}

	if pset.Count == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count is required")
		return
	}

	if pset.Result_id == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result_id is required")
		return
	}

	// 检测count
	if pset.Count < 1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("tightning count should be greater than 0")
		return
	}

	// 检测结果id
	result, err := apiserver.DB.GetResult(pset.Result_id, 0)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	// 通过控制器设定程序
	err = apiserver.CVI3.Service.PSet(pset.Controller_SN, pset.PSet, result.WorkorderID, pset.Result_id, pset.Count)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

// 根据hmi序列号以及vin或knr取得工单
func (apiserver *APIServer) getWorkorder(ctx iris.Context) {
	hmi_sn := ctx.URLParam("hmi_sn")
	vin_or_knr := ctx.URLParam("vin_or_knr")

	if hmi_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn is required")
		return
	}

	if vin_or_knr == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("vin_or_knr is required")
		return
	}

	var vin = ""
	var knr = ""

	if strings.Contains(vin_or_knr, payload.KNR_KEY) {
		knr = vin_or_knr
	} else {
		vin = vin_or_knr
	}

	workorder, err := apiserver.DB.FindWorkorder(hmi_sn, vin, knr)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	resp := payload.Workorder{}
	resp.HMI_sn = workorder.HMISN
	resp.PSet = workorder.PSet
	resp.Workorder_id = workorder.WorkorderID
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.Nut_total = workorder.NutTotal
	resp.Status = workorder.Status
	resp.WorkSheet = workorder.WorkSheet
	json.Unmarshal([]byte(workorder.ResultIDs), &resp.Result_ids)

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

// 创建工单
func (apiserver *APIServer) postWorkorders(ctx iris.Context) {
	var err error
	var workorders []payload.ODOOWorkorder
	err = ctx.ReadJSON(&workorders)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	_, e := apiserver.DB.InsertWorkorders(workorders)
	if e != nil {
		//fmt.Printf("%s\n", e.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	} else {
		// 推送新工单
		//go apiserver.push_new_orders(neworders)

		ctx.StatusCode(iris.StatusCreated)
		return
	}
}

//api文档
func (apiserver *APIServer) getDoc(ctx iris.Context) {
	f, _ := ioutil.ReadFile(apiserver.DocPath)

	ctx.Header("content-type", "application/json")
	ctx.Write(f)
}

func (apiserver *APIServer) push_new_orders(orders []payload.ODOOWorkorder) {
	for _, v := range orders {
		order_str, _ := json.Marshal(v)
		apiserver.WSSendWorkorder(v.HMI.UUID, string(order_str))
	}
}

func (apiserver *APIServer) getResults(ctx iris.Context) {
	// 根据查询参数返回结果
	has_upload := ctx.URLParam("has_upload")
	if has_upload == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload is required")
		return
	}

	result := ctx.URLParams()["result"]

	if result == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("result is required")
		return
	}
	re_list := strings.Split(result, ",")

	bool_has_upload, err := strconv.ParseBool(has_upload)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("has_upload value error")
		return
	}

	//list_result := []string{}
	//e := json.Unmarshal([]byte(result), &list_result)
	//if e != nil {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("result value error")
	//	return
	//}

	resp := []payload.ODOOResultSync{}
	results, _ := apiserver.DB.FindResults(bool_has_upload, re_list)
	target_results := map[int]rushdb.Results{}
	for _, v := range results {
		tr, exist := target_results[v.ResultId]
		if exist {
			// 已存在
			if v.Count > tr.Count {
				target_results[v.ResultId] = v
			}
		} else {
			// 不存在
			target_results[v.ResultId] = v
		}
	}

	for _, v := range target_results {
		odoo_result := payload.ODOOResultSync{}
		stime := strings.Split(v.UpdateTime.Format("2006-01-02 15:04:05"), " ")
		odoo_result.Control_date = fmt.Sprintf("%sT%s+08:00", stime[0], stime[1])

		odoo_result.CURObjects = []payload.CURObject{}

		curves, err := apiserver.DB.ListCurves(v.ResultId)
		if err != nil {
			for _, c := range curves {
				cur_object := payload.CURObject{}
				cur_object.File = c.CurveFile
				cur_object.OP = c.Count
				odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
			}
		}

		r := payload.ResultValue{}
		json.Unmarshal([]byte(v.ResultValue), &r)

		pset := payload.PSetDefine{}
		json.Unmarshal([]byte(v.PSetDefine), &pset)

		odoo_result.Measure_degree = r.Wi
		odoo_result.Measure_result = strings.ToLower(v.Result)
		odoo_result.Measure_t_don = r.Ti
		odoo_result.Measure_torque = r.Mi
		odoo_result.Op_time = v.Count
		odoo_result.Pset_m_max = pset.Mp
		odoo_result.Pset_m_min = pset.Mm
		odoo_result.Pset_m_target = pset.Ma
		odoo_result.Pset_m_threshold = pset.Ms
		odoo_result.Pset_strategy = pset.Strategy
		odoo_result.Pset_w_max = pset.Wp
		odoo_result.Pset_w_min = pset.Wm
		odoo_result.Pset_w_target = pset.Wa
		odoo_result.ID = v.ResultId

		resp = append(resp, odoo_result)
	}

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)

}

func (apiserver *APIServer) patchResult(ctx iris.Context) {
	id, err := strconv.Atoi(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	var e error
	var up payload.ResultPatch
	e = ctx.ReadJSON(&up)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}

	e = apiserver.DB.UpdateResults(id, 0, up.HasUpload)
	if e != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(e.Error())
		return
	}
}

func (apiserver *APIServer) getHealthz(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNoContent)
	return
}

func (apiserver *APIServer) getStatus(ctx iris.Context) {
	// 返回控制器状态

	sn := ctx.URLParam("controller_sn")
	status, err := apiserver.CVI3.Service.GetControllersStatus(sn)

	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(err.Error())
		return
	} else {
		body, _ := json.Marshal(status)
		ctx.Header("content-type", "application/json")
		ctx.Write(body)
	}
}

func (apiserver *APIServer) AddClient(sn string, c WSClient) {
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	apiserver.WSClients[sn] = c
}

func (apiserver *APIServer) RemoveClient(id string) {
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	var key string = ""
	for k, v := range apiserver.WSClients {
		if v.ID == id {
			key = k
			break
		}
	}

	if key != "" {
		delete(apiserver.WSClients, key)
	}
}

func (apiserver *APIServer) GetClient(sn string) (WSClient, bool) {
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	v, e := apiserver.WSClients[sn]
	return v, e
}

func (apiserver *APIServer) WSSendStatus(payload string) {
	// 群发控制器状态
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	for _, v := range apiserver.WSClients {
		v.Conn.Emit(WS_EVENT_STATUS, payload)
	}
}

func (apiserver *APIServer) WSSendResult(sn string, payload string) {
	// 将拧接结果发送给指定的控制器
	client, exist := apiserver.GetClient(sn)
	if exist {
		client.Conn.Emit(WS_EVENT_RESULT, payload)
	}
}

func (apiserver *APIServer) WSSendWorkorder(sn string, payload string) {
	// 将工单推送给制定控制器
	client, exist := apiserver.GetClient(sn)
	if exist {
		client.Conn.Emit(WS_EVENT_WORKORDER, payload)
	}
}

func (apiserver *APIServer) onWSConn(c websocket.Connection) {

	api := apiserver

	c.OnMessage(func(data []byte) {
		// 接受客户端链接
		fmt.Printf("recv from %s: %s\n", c.ID(), string(data))

		var reg_msg payload.WSRegistMsg
		var reg_str []byte
		reg := payload.WSRegist{}
		err := json.Unmarshal(data, &reg)
		if err != nil {
			reg_msg.Msg = "regist msg error"
			reg_str, _ = json.Marshal(reg_msg)
			c.EmitMessage(reg_str)
		}

		_, exist := apiserver.GetClient(reg.HMI_SN)
		if exist {
			reg_msg.Msg = fmt.Sprintf("client with sn:%s already exists", reg.HMI_SN)
			reg_str, _ = json.Marshal(reg_msg)
			c.EmitMessage(reg_str)
		} else {
			// 将客户端加入列表
			client := WSClient{}
			client.Conn = c
			client.ID = c.ID()
			apiserver.AddClient(reg.HMI_SN, client)
			reg_msg.Msg = "OK"
			reg_str, _ = json.Marshal(reg_msg)
			c.EmitMessage(reg_str)
		}
	})

	c.OnDisconnect(func() {
		fmt.Printf("Connection with id: %s has been disconnected!\n", c.ID())
		// 删除客户端
		api.RemoveClient(c.ID())
	})

}

func (apiserver *APIServer) StartService(doc_path string) error {

	apiserver.DocPath = doc_path
	apiserver.WSClients = map[string]WSClient{}
	apiserver.WSMtx = sync.Mutex{}

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "content-type", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Screen"},
		AllowCredentials: true,
	})

	ws := websocket.New(websocket.Config{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	})

	ws.OnConnection(apiserver.onWSConn)

	v1 := app.Party(API_PREFIX, crs).AllowMethods(iris.MethodOptions)
	{
		v1.Put("/psets", apiserver.putPSets)
		v1.Get("/workorder", apiserver.getWorkorder)
		v1.Get("/results", apiserver.getResults)
		v1.Patch("/results/{id:int}", apiserver.patchResult)
		v1.Get("/controller-status", apiserver.getStatus)
		v1.Get("/healthz", apiserver.getHealthz)
		v1.Post("/workorders", apiserver.postWorkorders)
		v1.Get("/doc", apiserver.getDoc)
		app.Get("/ws", ws.Handler())
	}

	//websocket.ClientSource
	return app.Run(iris.Addr(apiserver.Port), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
