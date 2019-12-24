package ts002

import (
	"gopkg.in/resty.v1"
	"time"
)

func NewMesAPI(cfg *MesApiConifg) (*MesAPI, error) {
	client := resty.New()
	client.SetRESTMode()
	client.SetTimeout(time.Duration(cfg.Timeout))
	client.SetContentLength(true)
	client.
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(time.Duration(cfg.RetryInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	return &MesAPI{
		httpClient: client,
		cfg:        cfg,
	}, nil
}

// 中车数字相关接口
type MesAPI struct {
	httpClient *resty.Client
	cfg        *MesApiConifg
}
