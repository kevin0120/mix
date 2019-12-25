package reader

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
}
