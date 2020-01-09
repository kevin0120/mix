package aiis

const (
	TransportTypeBroker = "broker"
	TransportTypeGrpc   = "grpc"
)

const (
	TransportStatusOnline  = "online"
	TransportStatusOffline = "offline"
)

const (
	TransportMethodResult        = "method_result"
	TransportMethodResultPatch   = "method_result_patch"
	TransportMethodServiceStatus = "method_service_status"
	//TransportMethodWorkorder     = "method_workorder"
	//TransportMethodOrderStart    = "method_order_start"
	//TransportMethodOrderFinish   = "method_order_finish"
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
