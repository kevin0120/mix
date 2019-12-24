package ts002

import (
	"github.com/masami10/rush/utils"
	"gopkg.in/resty.v1"
	"time"
)

func NewMesAPI(cfg *MesApiConifg, diag Diagnostic) (*MesAPI, error) {
	api := MesAPI{
		cfg:  cfg,
		diag: diag,
	}

	api.ensureHttpClient()
	return &api, nil
}

// 中车数字相关接口
type MesAPI struct {
	client *resty.Client
	cfg    *MesApiConifg
	diag   Diagnostic
}

func (s *MesAPI) ensureHttpClient() *resty.Client {
	if s.client != nil {
		return s.client
	}

	if client, err := utils.CreateRetryClient(time.Duration(s.cfg.Timeout), s.cfg.RetryCount); err != nil {
		s.diag.Error("ensureHttpClient", err)
		return nil
	} else {
		s.client = client
		return client
	}
}

func (s *Service) SendNFCData() error {
	//todo: 发送读卡器数据
	//client := s.ensureHttpClient()
	//r := client.R()
	return nil
}
