package transport

import (
	"github.com/masami10/rush/utils"
	"github.com/micro/go-micro/transport"
	"time"
)

type Diagnostic interface {
	Error(string, error)
	Debug(string)
	Info(string)
	SendMsgSuccess(string)
}


type IServer interface {
	GetServiceByName(string) utils.ICommonService
}

type ITransportService interface {
	GetServerAddress() []string //服务器地址，可能存在集群服务器如：etcd， nats等
	TransportForceOpen() error  //强制打开特定的transport服务
	OnMessage(subject string, handler OnMsgHandler) error
	SendMessage(subject string, data []byte) error
	Request(subject string, data []byte, timeOut time.Duration) ([]byte, error)
	Status() string
}

type Transport = transport.Transport

type Message = transport.Message

type Socket = transport.Socket

type Client = transport.Client

type Listener = transport.Listener

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)

type OnMsgHandler func(*Message) ([]byte, error)
