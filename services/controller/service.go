package controller

type Diagnostic interface {
	Error(msg string, err error)
}

type Controller interface {
	Start()
}

type Protocol interface {
	Parse(buf []byte) ([]byte, error)
	Write(sn string, buf []byte) error
	AddNewController(cfg Config)
}

type Service struct {
	configs []Config

	protocols map[string]Protocol //进行服务注入, serial_no : Protocol

	diag Diagnostic
}

func NewService(cs Configs, d Diagnostic, pAudi Protocol) (*Service, error) {
	s := &Service{
		configs:   cs,
		diag:      d,
		protocols: map[string]Protocol{},
	}

	for _, c := range cs {
		switch c.Protocol {
		case AUDIPROTOCOL:
			pAudi.AddNewController(c)
			s.protocols[c.SN] = pAudi

		default:

		}
	}

	return s, nil
}

func (s *Service) Write(serial_no string, buf []byte) error {

	controller := s.protocols[serial_no]

	return controller.Write(serial_no, buf)
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}
