package ts002

import "github.com/masami10/rush/services/httpd"

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
}

type IHttpService interface {
	GetHandlerByName(version string) (*httpd.Handler, error)
}
