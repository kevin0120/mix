package aiis

const (
	TRANSPORT_TYPE_BROKER = "broker"
	TRANSPORT_TYPE_GRPC   = "grpc"
)

const (
	TRANSPORT_METHOD_RESULT         = "method_result"
	TRANSPORT_METHOD_RESULT_PATCH   = "method_result_patch"
	TRANSPORT_METHOD_SERVICE_STATUS = "method_service_status"
	TRANSPORT_METHOD_WORKORDER      = "method_workorder"
	TRANSPORT_METHOD_ORDER_START    = "method_order_start"
	TRANSPORT_METHOD_ORDER_FINISH   = "method_order_finish"
)

type ServiceStatusHandler func(status ServiceStatus)
type ResultPatchHandler func(rp ResultPatch)
type StatusHandler func(status string)

type ITransport interface {
	// 上传结果
	SendResult(result *AIISResult) error

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
	ITransport
	handlerServiceStatus ServiceStatusHandler
	handlerResultPatch   ResultPatchHandler
	handlerStatus        StatusHandler
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

func (s *BaseTransport) SendResult(result *AIISResult) error {
	return nil
}
