package aiis

import (
	"encoding/json"
	"fmt"

	"github.com/kataras/iris/core/errors"
)

type ServiceStatusHandler func(status ServiceStatus)
type ResultPatchHandler func(rp ResultPatch)
type StatusHandler func(status string)

type ITransport interface {
	// 上传结果
	SendResult(result *PublishResult) error

	// 设置连接状态回调
	SetStatusHandler(handler StatusHandler)

	// 设置接收服务状态回调
	SetServiceStatusHandler(handler ServiceStatusHandler)

	// 设置接收结果上传反馈回调
	SetResultPatchHandler(handler ResultPatchHandler)

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
	handlerStatus        StatusHandler
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

	return trans.SendMessage(fmt.Sprintf(SubjectResults, result.ToolSN), data)
}

func (s *BaseTransport) SetServiceStatusHandler(handler ServiceStatusHandler) {
	s.handlerServiceStatus = handler
}

func (s *BaseTransport) SetResultPatchHandler(handler ResultPatchHandler) {
	s.handlerResultPatch = handler
}

func (s *BaseTransport) SetStatusHandler(handler StatusHandler) {
	s.handlerStatus = handler
}

func NewAIISBaseTransport(t ITransportService) *BaseTransport {
	s := &BaseTransport{}
	s.ITransportService = t
	return s
}
