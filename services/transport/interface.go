package transport

import (
	"github.com/micro/go-micro/transport"
)

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
	TransportForceOpen()        //强制打开特定的transport服务
	Transport
}

type Transport = transport.Transport

type Message = transport.Message

type Socket = transport.Socket

type Client = transport.Client

type Listener = transport.Listener

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)