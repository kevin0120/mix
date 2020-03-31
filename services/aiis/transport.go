package aiis

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/transport"

	"github.com/kataras/iris/core/errors"
)

type ServiceStatusHandler func(status *ServiceStatus)
type ResultPatchHandler func(rp *ResultPatch)
type StatusHandler = transport.StatusHandler

type ITransport interface {
	// 上传结果
	SendResult(result *PublishResult) error

	// 设置连接状态回调
	SetStatusHandler(handler StatusHandler) error

	// 设置接收服务状态回调
	SetServiceStatusHandler(handler ServiceStatusHandler)

	// 设置接收结果上传反馈回调
	SetResultPatchHandler(handler ResultPatchHandler) error

	// 启动
	Start() error

	// 停止
	Stop() error

	// 获取状态
	Status() string
}

type BaseTransport struct {
	ITransportService
	handlerServiceStatus ServiceStatusHandler
	handlerResultPatch   ResultPatchHandler
}

func (s *BaseTransport) Start() error {
	return nil
}

func (s *BaseTransport) Stop() error {
	return nil
}

func (s *BaseTransport) Status() string {
	return s.ITransportService.Status()
}

func (s *BaseTransport) SendResult(result *PublishResult) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	data, _ := json.Marshal(TransportPayload{
		Method: TransportMethodResult,
		Data:   result,
	})

	return trans.SendMessage(fmt.Sprintf(SubjectResults, trans.GetID(), result.WorkcenterCode, result.ToolSN), data)
}

func (s *BaseTransport) SetServiceStatusHandler(handler ServiceStatusHandler) {
	s.handlerServiceStatus = handler
}

func (s *BaseTransport) SetResultPatchHandler(handler ResultPatchHandler) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	subject := fmt.Sprintf(SubjectResultsResp, trans.GetID()) //返回指定的客户端, 根据客户端ID标识
	fn := func(msg *transport.Message) ([]byte, error) {
		var payload TransportPayload
		data := msg.Body
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}

		body, _ := json.Marshal(payload.Data)
		var result ResultPatch
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		handler(&result)
		return nil, nil
	}
	return trans.OnMessage(subject, fn)
}

func (s *BaseTransport) SetStatusHandler(handler StatusHandler) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	return trans.SetStatusHandler(handler)
}

func NewAIISBaseTransport(t ITransportService) *BaseTransport {
	s := &BaseTransport{}
	s.ITransportService = t
	return s
}
