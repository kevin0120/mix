package changan

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/wsnotify"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	WS           *wsnotify.Service
	HTTPDService *httpd.Service
	diag         Diagnostic
	configValue  atomic.Value
}

func NewService(d Diagnostic, c Config, httpd *httpd.Service, ws *wsnotify.Service) *Service {
	if c.Enable {
		s := &Service{
			diag:         d,
			WS:           ws,
			HTTPDService: httpd,
		}

		s.configValue.Store(c)
		return s
	}

	return nil
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/andon-test",
		HandlerFunc: s.andonTest,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) andonTest(ctx iris.Context) {
	andon_msg := AndonMsg{}
	err := ctx.ReadJSON(&andon_msg)

	if err != nil {
		ctx.Writef(fmt.Sprintf("param error: %s", err.Error()))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	str_data, _ := json.Marshal(andon_msg.Data)

	switch andon_msg.MsgType {
	case MSG_TASK:
		tasks := []AndonTask{}
		err = json.Unmarshal(str_data, &tasks)
		if err != nil {
			ctx.Writef(fmt.Sprintf("tasks error: %s", err.Error()))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		for _, v := range tasks {
			t, _ := json.Marshal(v)
			s.WS.WSSendTask(v.Workcenter, string(t))

			//fmt.Printf("send task -- workcenter:%s payload:%s\n", v.Workcenter, string(t))
		}
	}
}
