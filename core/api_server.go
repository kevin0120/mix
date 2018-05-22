package core

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"github.com/masami10/rush/db"
	"encoding/json"
	"github.com/masami10/rush/payload"
	"github.com/kataras/iris/websocket"
	"fmt"
	"sync"
	"strings"
	"strconv"
)

const (
	API_PREFIX = "/api/v1"

	WS_EVENT_STATUS = "status"
	WS_EVENT_RESULT = "result"
	WS_EVENT_WORKORDER = "workorder"
)

func CleanClient(id string) {

}

type WSClient struct {
	ID	string
	Conn websocket.Connection
}

type APIServer struct {
	Port string
	CVI3 *CVI3Service
	DB	*rushdb.DB
	WSClients	map[string]WSClient
	WSMtx	sync.Mutex
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
	var oldresult rushdb.Results
	oldresult, err = apiserver.DB.GetResult(pset.Result_id, 0)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	var result rushdb.Results
	result, err = apiserver.DB.GetResult(pset.Result_id, pset.Count)
	if err != nil {
		// 创建新结果
		nr := rushdb.Results{}
		nr.Workorder_id = oldresult.Workorder_id
		nr.Controller_sn = pset.Controller_SN
		nr.Result = payload.RESULT_NONE
		nr.Count = pset.Count
		nr.Cur_data = ""
		nr.Cur_upload = false
		nr.Result_id = pset.Result_id
		nr.Result_data = ""
		nr.Result_upload = false
		err := apiserver.DB.InsertResults(nr)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
	} else {
		result = oldresult
	}

	// 通过控制器设定程序
	err = apiserver.CVI3.Service.PSet(pset.Controller_SN, pset.PSet, result.Workorder_id, pset.Result_id, pset.Count)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}

// 根据hmi序列号以及vin或knr取得工单
func (apiserver *APIServer) getWorkorder(ctx iris.Context) {
	hmi_sn := ctx.URLParam("hmi_sn")
	vin := ctx.URLParam("vin")
	knr := ctx.URLParam("knr")

	if hmi_sn == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn is required")
		return
	}

	if vin == "" && knr == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("vin or knr is required")
		return
	}

	workorder, err := apiserver.DB.FindWorkorder(hmi_sn, vin, knr)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	resp := payload.Workorder{}
	resp.HMI_sn = workorder.HMI_sn
	resp.PSet = workorder.PSet
	resp.Workorder_id = workorder.Workorder_id
	resp.Vin = workorder.Vin
	resp.Knr = workorder.Knr
	resp.Nut_total = workorder.Nut_total
	resp.Status = workorder.Status
	resp.WorkSheet = workorder.WorkSheet
	json.Unmarshal([]byte(workorder.Result_ids), &resp.Result_ids)

	body, _ := json.Marshal(resp)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

// 创建工单
func (apiserver *APIServer) createWorkorders(ctx iris.Context) {
	var err error
	var workorders []payload.ODOOWorkorder
	err = ctx.ReadJSON(&workorders)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	}

	neworders, e := apiserver.DB.InsertWorkorders(workorders)
	if e != nil {
		//fmt.Printf("%s\n", e.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())

		return
	} else {
		// 推送新工单
		go apiserver.push_new_orders(neworders)

		ctx.StatusCode(iris.StatusCreated)
		return
	}
}

func (apiserver *APIServer) push_new_orders(orders []payload.ODOOWorkorder) {
	for _,v := range orders {
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
	target_results :=  map[int]rushdb.Results{}
	for _, v := range results {
		tr, exist := target_results[v.Result_id]
		if exist {
			// 已存在
			if v.Count > tr.Count {
				target_results[v.Result_id] = v
			}
		} else {
			// 不存在
			target_results[v.Result_id] = v
		}
	}

	for _, v := range target_results {
		odoo_result := payload.ODOOResultSync{}
		stime := strings.Split(v.Update_time.Format("2006-01-02 15:04:05"), " ")
		odoo_result.Control_date = fmt.Sprintf("%sT%s+08:00", stime[0], stime[1])

		odoo_result.CURObjects = []payload.CURObject{}

		results, _ := apiserver.DB.ListResults(v.Result_id)
		for _,vi := range results {
			nr := payload.Result{}
			json.Unmarshal([]byte(vi.Result_data), &nr)
			cur_object := payload.CURObject{}
			cur_object.File = nr.CurFile
			cur_object.OP = nr.Count
			odoo_result.CURObjects = append(odoo_result.CURObjects, cur_object)
		}

		or := payload.Result{}
		json.Unmarshal([]byte(v.Result_data), &or)

		odoo_result.Measure_degree = or.ResultValue.Wi
		odoo_result.Measure_result = strings.ToLower(or.Result)
		odoo_result.Measure_t_don = or.ResultValue.Ti
		odoo_result.Measure_torque = or.ResultValue.Mi
		odoo_result.Op_time = or.Count
		odoo_result.Pset_m_max = or.PSetDefine.Mp
		odoo_result.Pset_m_min = or.PSetDefine.Mm
		odoo_result.Pset_m_target = or.PSetDefine.Ma
		odoo_result.Pset_m_threshold = or.PSetDefine.Ms
		odoo_result.Pset_strategy = or.PSetDefine.Strategy
		odoo_result.Pset_w_max = or.PSetDefine.Wp
		odoo_result.Pset_w_min = or.PSetDefine.Wm
		odoo_result.Pset_w_target = or.PSetDefine.Wa
		odoo_result.ID = v.Result_id

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

func (apiserver *APIServer) getStatus(ctx iris.Context) {
	// 返回控制器状态

	//hmi_sn := ctx.URLParam("hmi_sn")
	//if hmi_sn != "" {
	//	// 指定控制器
	//	apiserver.CVI3.Service.
	//} else {
	//	// 所有控制器
	//}
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
	for k,v := range apiserver.WSClients {
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
	return v,e
}

func (apiserver *APIServer) WSSendStatus(payload string) {
	// 群发控制器状态
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	for _,v := range apiserver.WSClients {
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


func (apiserver *APIServer) StartService() error {

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


	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Put("/psets", apiserver.putPSets)
		v1.Get("/workorder", apiserver.getWorkorder)
		v1.Get("/results", apiserver.getResults)
		v1.Patch("/results/{id:int}", apiserver.patchResult)
		v1.Get("/status", apiserver.getStatus)
		v1.Post("/workorders", apiserver.createWorkorders)
		app.Get("/ws", ws.Handler())
	}

	//websocket.ClientSource
	return app.Run(iris.Addr(apiserver.Port), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
