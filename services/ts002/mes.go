package ts002

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"sync/atomic"
	"time"
)

var mapMESStatusIO = map[string]uint16{
	"on":  io.OUTPUT_STATUS_ON,
	"off": io.OUTPUT_STATUS_OFF,
}

func NewMesAPI(cfg MesApiConfig, diag Diagnostic) (*MesAPI, error) {
	api := MesAPI{
		diag: diag,
	}

	api.cfg.Store(cfg)

	api.ensureHttpClient()
	return &api, nil
}

// 中车数字相关接口
type MesAPI struct {
	client *resty.Client
	cfg    atomic.Value
	diag   Diagnostic
}

func (s *MesAPI) Config() MesApiConfig {
	return s.cfg.Load().(MesApiConfig)
}

func (s *MesAPI) healthCheck() {
	c := s.Config()
	cc := s.ensureHttpClient()
	r := cc.R()
	url := fmt.Sprintf("%s%s", c.APIUrl, c.EndpointHealthzCheck)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for ; ; {
		select {
		case <-ticker.C:
			if resp, err := r.Get(url); err != nil {
				s.diag.Error("MES healzCheck", err)
			} else {
				if resp.StatusCode() != iris.StatusNoContent {
					e := errors.New("MES healzCheck Fail")
					s.diag.Error("healzCheck", e)
				}
			}
		}
	}
}

func (s *MesAPI) ensureHttpClient() *resty.Client {
	c := s.Config()
	if s.client != nil {
		return s.client
	}

	if client, err := utils.CreateRetryClient(time.Duration(c.Timeout), c.RetryCount); err != nil {
		s.diag.Error("ensureHttpClient", err)
		return nil
	} else {
		client.SetHeader("Project", "TS002")
		client.SetHeader("Service", "MES")
		s.client = client
		return client
	}
}

func (s *MesAPI) sendNFCData(data string) error {
	//todo: 发送读卡器数据
	c := s.Config()
	client := s.ensureHttpClient()
	r := client.R()

	url := fmt.Sprintf("%s%s", c.APIUrl, c.EndpointCardInfo)
	payload := RushCardInfoReq{CardCode: data}
	if resp, err := r.SetBody(payload).Post(url); err != nil {
		s.diag.Error("sendNFCData Error", err)
		return err
	} else if resp.StatusCode() != iris.StatusOK {
		err := errors.New(resp.String())
		s.diag.Error("sendNFCData Error", err)
		return err
	}
	return nil
}
