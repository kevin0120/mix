package grpc

import "github.com/masami10/rush/services/transport"

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type RespStruct struct {
	msgId string
	msg   chan *transport.Message
}
