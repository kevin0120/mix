package dispatcherbus

//todo: 通过封装getattr, setattr 减少这种赋值代码
const (
	// ********************************Device***********************************
	// 设备状态(包括IO，条码枪，读卡器，拧紧控制器，拧紧工具等)发生变化时，会通过此分发器进行状态分发
	DISPATCHER_DEVICE_STATUS = "DEVICE_STATUS"

	// ********************************Scanner***********************************
	// 当收到条码数据时(来自条码枪，拧紧控制器条码等)，会通过此分发器进行条码分发
	DISPATCHER_SCANNER_DATA = "DISPATCHER_SCANNER_DATA"

	// ********************************IO***********************************
	// 当收到IO输入输出状态变化时(IO模块或拧紧控制器IO等)，会通过此分发器进行IO状态分发
	DISPATCH_IO = "DISPATCH_IO"

	// ********************************Reader***********************************
	// 当收到读卡器数据时，会通过此分发器进行读卡器数据分发
	DISPATCHER_READER_DATA = "DISPATCHER_READER_DATA"

	// ********************************Tightening***********************************
	// 当收到拧紧结果时，会通过此分发器进行拧紧结果分发
	DISPATCH_RESULT = "DISPATCH_RESULT"

	// 当收到拧紧曲线时，会通过此分发器进行拧紧曲线分发
	DISPATCH_CURVE = "DISPATCH_CURVE"

	// 当检测到新工具时，会通过此分发器进行新工具分发
	DISPATCH_NEW_TOOL = "DISPATCH_NEW_TOOL"

	// ********************************AIIS***********************************
	DISPATCH_ODOO_STATUS = "DISPATCH_ODOO_STATUS"
	DISPATCH_AIIS_STATUS = "DISPATCH_AIIS_STATUS"

	DISPATCH_RPC_STATUS = "DISPATCH_RPC_STATUS"
	DISPATCH_RPC_RECV   = "DISPATCH_RPC_RECV"

	// ********************************BROKER***********************************
	// 当Broker(MQ)服务状态变化时， 会通过此分发器进行状态分发
	DISPATCH_BROKER_STATUS = "DISPATCH_BROKER_STATUS"

	// ********************************WEBSOCKET***********************************
	// 当收到WebSocket请求时， 会通过此分发器进行请求分发
	DISPATCH_WS_NOTIFY = "DISPATCH_WS_NOTIFY"
)
