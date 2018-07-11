package controller

type Diagnostic interface {
	Error(msg string, err error)
}

type Controller interface {
	Start()
}

type Protocol interface {
	Parse(msg string) ([]byte, error)
	Write(sn string, buf []byte) error
	AddNewController(cfg Config)
}

type Service struct {
	configs []Config

	protocols map[string]Protocol //进行服务注入, serial_no : Protocol

	diag Diagnostic
}

func NewService(cs Configs, d Diagnostic, pAudi Protocol, pOpenprotocol Protocol) (*Service, error) {
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

		case OPENPROTOCOL:
			pOpenprotocol.AddNewController(c)
			s.protocols[c.SN] = pOpenprotocol

		default:

		}
	}

	return s, nil
}

func (s *Service) Write(serialNo string, buf []byte) error {

	controller := s.protocols[serialNo]

	return controller.Write(serialNo, buf)
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}
