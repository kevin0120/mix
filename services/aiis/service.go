package aiis

import (
	"sync/atomic"
	"time"

	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"net/url"
	"strings"
)

type Diagnostic interface {
	Error(msg string, err error)
	PutResultDone()
}

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
}

func NewEndpoint(url string, headers map[string]string, method string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
	}
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	endpoints   []*Endpoint
	httpClient  *resty.Client
	rush_port   string
	db          *storage.Service
	ws          utils.RecConn
	//ws			websocket.Client
	SN string
}

func NewService(c Config, d Diagnostic, rush_port string) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		endpoints: e,
		rush_port: rush_port,
	}
	s.configValue.Store(c)
	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(c.Timeout))
	client.SetContentLength(true)
	// Headers for all request
	client.SetHeaders(c.Headers)
	client.
		SetRetryCount(c.MaxRetry).
		SetRetryWaitTime(time.Duration(c.PushInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client

	entry := strings.Split(s.Config().Urls[0], "://")[1]
	url := url.URL{Scheme: "ws", Host: entry, Path: s.Config().WSResultRoute}
	s.ws = utils.RecConn{}
	s.ws.OnConnected = func() {
		ws_msg := WSMsg{
			Type: WS_REG,
			Data: WSRegist{
				Rush_SN: s.SN,
			},
		}

		str, _ := json.Marshal(ws_msg)
		s.ws.WriteMessage(websocket.TextMessage, str)
	}

	s.ws.Dial(url.String(), nil)

	go s.ResultUploadManager()

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) PutResult(result_id int64, body interface{}) error {

	//var err error
	//for _, endpoint := range s.endpoints {
	//	err = s.putResult(body, fmt.Sprintf(endpoint.url, result_id), endpoint.method)
	//	if err == nil {
	//		// 如果第一次就成功，推出循环
	//		return nil
	//	}
	//}

	ws_msg := WSMsg{
		Type: WS_RESULT,
		Data: WSOpResult{
			ResultID: result_id,
			Result:   body.(AIISResult),
			Port:     s.rush_port,
		},
	}

	str, _ := json.Marshal(ws_msg)
	s.ws.WriteMessage(websocket.TextMessage, str)
	//s.ws.SendText(string(str))
	return nil
}

func (s *Service) putResult(body interface{}, url string, method string) error {
	r := s.httpClient.R().SetBody(body).SetHeader("rush_port", s.rush_port)
	var resp *resty.Response
	var err error

	switch method {
	case "PATCH":
		resp, err = r.Patch(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "PUT":
		resp, err = r.Put(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "POST":
		resp, err = r.Post(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	default:
		return errors.New("Result Put :the Method is wrong")

	}
	s.diag.PutResultDone()
	return nil
}

func (s *Service) ResultUploadManager() error {
	for {

		//results, err := s.db.ListUnuploadResults()
		//if err == nil {
		//	//for _, v := range results {
		//	//	//s.PutResult(v.ResultId, v)
		//	//}
		//}

		time.Sleep(time.Duration(s.Config().ResultUploadInteval))
	}
}
