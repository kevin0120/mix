package ts002


type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
}
