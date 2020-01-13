package grpc



type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}
