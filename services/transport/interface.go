package transport

type Diagnostic interface {
	Error(string, error)
	Debug(string)
	Info(string)
	SendMsgSuccess(string)
}

// Service represents a service attached to the server.
type ICommonService interface {
	Open() error
	Close() error
}

type IServer interface {
	GetServiceByName(string) ICommonService
}

type ITransportService interface {
	GetServerAddress() []string //服务器地址，可能存在集群服务器如：etcd， nats等
}
