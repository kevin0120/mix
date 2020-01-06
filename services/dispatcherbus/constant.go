package dispatcherbus

const (
	// ********************************Device***********************************
	// 设备状态(包括IO，条码枪，读卡器，拧紧控制器，拧紧工具等)发生变化时，会通过此分发器进行状态分发
	DISPATCHER_DEVICE_STATUS = "DISPATCHER_DEVICE_STATUS"

	// ********************************Scanner***********************************
	// 当收到条码数据时(来自条码枪，拧紧控制器条码等)，会通过此分发器进行条码分发
	DISPATCHER_SCANNER_DATA = "DISPATCHER_SCANNER_DATA"

	// ********************************IO***********************************
	// 当收到IO输入输出状态变化时(IO模块或拧紧控制器IO等)，会通过此分发器进行IO状态分发
	DISPATCHER_IO = "DISPATCHER_IO"

	// ********************************Reader***********************************
	// 当收到读卡器数据时，会通过此分发器进行读卡器数据分发
	DISPATCHER_READER_DATA = "DISPATCHER_READER_DATA"

	// ********************************Tightening***********************************
	// 当收到拧紧结果时，会通过此分发器进行拧紧结果分发
	DISPATCHER_RESULT = "DISPATCHER_RESULT"

	// 当收到拧紧曲线时，会通过此分发器进行拧紧曲线分发
	DISPATCHER_CURVE = "DISPATCHER_CURVE"

	// 当检测到新工具时，会通过此分发器进行新工具分发
	DISPATCHER_NEW_TOOL = "DISPATCHER_NEW_TOOL"

	// ********************************Service***********************************
	// 当检测到服务状态变化时(aiis, odoo, 外部系统等)，会通过此分发器进行状态分发
	DISPATCHER_SERVICE_STATUS = "DISPATCHER_SERVICE_STATUS"

	DISPATCHER_RPC_RECV = "DISPATCHER_RPC_RECV"

	// ********************************BROKER***********************************
	// 当Broker(MQ)服务状态变化时， 会通过此分发器进行状态分发
	DISPATCHER_BROKER_STATUS = "DISPATCHER_BROKER_STATUS"

	// ********************************WEBSOCKET***********************************
	// 当收到WebSocket请求时， 会通过此分发器进行请求分发
	DISPATCHER_WS_NOTIFY = "DISPATCHER_WS_NOTIFY"
)
