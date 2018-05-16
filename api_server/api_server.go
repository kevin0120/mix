package rush_api_server

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"github.com/masami10/rush/db"
	"encoding/json"
	"github.com/masami10/rush/cvi3_services"
	"github.com/masami10/rush/payload"
	"github.com/kataras/iris/websocket"
	"fmt"
	"sync"
)

const (
	API_PREFIX = "/api/v1"
)

type APIServer struct {
	Port string
	CVI3 *rush_cvi3.CVI3Service
	DB	*rushdb.DB
	WSClients	map[string]websocket.Connection
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
	var result rushdb.Results
	result, err = apiserver.DB.GetResult(pset.Result_id, pset.Count)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
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

	workorder, err := apiserver.DB.GetWorkorder(hmi_sn, vin, knr)
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

func (apiserver *APIServer) getResults(ctx iris.Context) {
	// 根据查询参数返回结果
}

func (apiserver *APIServer) getStatus(ctx iris.Context) {
	// 返回控制器状态

	//apiserver.CVI3.Service.
}

func (apiserver *APIServer) onWSConn(c websocket.Connection) {
	c.On("chat", func(msg string) {
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		//c.Emit("chat", "wef")
	})
}

func (apiserver *APIServer) AddClient(sn string, c websocket.Connection) {
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	apiserver.WSClients[sn] = c
}

func (apiserver *APIServer) GetClient(sn string) websocket.Connection {
	defer apiserver.WSMtx.Unlock()

	apiserver.WSMtx.Lock()
	return apiserver.WSClients[sn]
}

func (apiserver *APIServer) WSSend(sn string, payload string) {
	client := apiserver.GetClient(sn)
	if client != nil {
		client.Emit("chat", payload)
	}
}

func (apiserver *APIServer) StartService() error {

	apiserver.WSClients = map[string]websocket.Connection{}

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "content-type", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Screen"},
		AllowCredentials: true,
	})

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Put("/psets", apiserver.putPSets)
		v1.Get("/workorder", apiserver.getWorkorder)
		v1.Get("/results", apiserver.getResults)
		v1.Get("/status", apiserver.getStatus)
	}


	ws := websocket.New(websocket.Config{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	})

	ws.OnConnection(apiserver.onWSConn)
	app.Get("/ws", ws.Handler())
	//websocket.ClientSource
	return app.Run(iris.Addr(apiserver.Port), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
