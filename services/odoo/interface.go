package odoo

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	Release(name string, handler string) error
}
