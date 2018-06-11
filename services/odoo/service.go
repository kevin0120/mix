package odoo

import "gopkg.in/resty.v1"


type Diagnostic interface {
	Error(msg string, err error)
}


type Service struct {
	url  	string
	diag 	Diagnostic
	httpClient  *resty.Client
}



func NewService(c Config, d Diagnostic) *Service{
	return &Service{
		url: 		c.URL,
		diag: 		d,
	}
}

func (s *Service)Open() error {
	return nil
}

func (s *Service)Close() error {
	return nil
}

func (s *Service)CreateMO() error {

}