package server

import (
	"fmt"
	"github.com/masami10/rush/command"
	"github.com/masami10/rush/keyvalue"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/broker"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/diagnostic"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/grpc"
	"github.com/masami10/rush/services/hmi"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/openprotocol/vendors"
	"github.com/masami10/rush/services/reader"
	"github.com/masami10/rush/services/scanner"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/transport"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
)

type BuildInfo struct {
	Version  string
	Commit   string
	Branch   string
	Platform string
}

type Diagnostic interface {
	Debug(msg string, ctx ...keyvalue.T)
	Info(msg string, ctx ...keyvalue.T)
	Error(msg string, err error, ctx ...keyvalue.T)
}

// Service represents a service attached to the server.
type Service = utils.ICommonService

type Server struct {
	dataDir  string
	hostname string

	StorageService *storage.Service

	HTTPDService        *httpd.Service
	OdooService         *odoo.Service
	AudiVWService       *audi_vw.Service
	OpenprotocolService *openprotocol.Service

	WSNotifyService *wsnotify.Service
	AiisService     *aiis.Service
	MinioService    *minio.Service

	ScannerService *scanner.Service
	IOService      *io.Service
	ReaderService  *reader.Service

	TighteningDeviceService *tightening_device.Service
	DeviceService           *device.Service
	DispatcherBusService    *dispatcherbus.Service

	BrokerService *broker.Service

	GRPCService *grpc.Service

	TransportService *transport.Service

	config *Config
	// List of services in startup order
	Services []Service

	ServicesByName map[string]int
	err            chan error

	BuildInfo   BuildInfo
	Commander   command.Commander
	DiagService *diagnostic.Service
	Diag        Diagnostic
}

// New returns a new instance of Server built from a config.
func New(c *Config, buildInfo BuildInfo, diagService *diagnostic.Service) (*Server, error) {
	err := c.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %s. To generate a valid configuration file run `rushd config > rush.generated.conf`. ", err)
	}
	d := diagService.NewServerHandler()
	s := &Server{
		config:         c,
		DiagService:    diagService,
		Diag:           d,
		BuildInfo:      buildInfo,
		dataDir:        c.DataDir,
		hostname:       c.Hostname,
		ServicesByName: make(map[string]int),
		err:            make(chan error),
		Commander:      c.Commander,
	}
	if err := s.appendDispatcherBus(); err != nil {
		return nil, errors.Wrap(err, "appendDispatcherBus")
	}

	if err := s.initHTTPDService(); err != nil {
		return nil, errors.Wrap(err, "init httpd service")
	}

	s.appendStorageService()

	s.appendWebsocketService()

	if err := s.initAudiVWProtocolService(); err != nil {
		return nil, errors.Wrap(err, "init Audi/VW service")
	}

	s.appendMinioService()

	s.appendGRPCService()

	s.appendBrokerService()

	if err := s.appendDeviceService(); err != nil {
		return nil, errors.Wrap(err, "appendDeviceService")
	}

	s.appendAudiVWService() //此服务必须在控制器服务后进行append

	s.appendOpenProtocolService()

	s.appendTighteningDeviceService()

	s.appendOdooService()

	s.appendHMIService()

	s.appendScannerService()

	s.appendIOService()

	s.appendReaderService()

	if err := s.appendTransportService(); err != nil {
		return nil, err
	} // append transport service

	if err := s.appendAiisService(); err != nil {
		return nil, err
	}

	s.appendHTTPDService()

	return s, nil
}

func (s *Server) appendTransportService() error {
	c := s.config.Transport
	d := s.DiagService.NewTransportHandler()
	srv := transport.NewService(c, d)

	if err := srv.BindTransportByProvider(s); err != nil {
		s.Diag.Error("BindTransportByProvider", err)
		return err
	}

	s.TransportService = srv

	s.AppendService("transport", srv)
	return nil
}

func (s *Server) GetServiceByName(name string) Service {
	if idx, ok := s.ServicesByName[name]; !ok {
		// Should be unreachable code
		return nil
	} else {
		return s.Services[idx]
	}
}

func (s *Server) AppendService(name string, srv Service) {
	if _, ok := s.ServicesByName[name]; ok {
		// Should be unreachable code
		panic("cannot append service twice")
	}
	i := len(s.Services)
	s.Services = append(s.Services, srv)
	s.ServicesByName[name] = i
}

func (s *Server) appendDispatcherBus() error {
	d := s.DiagService.NewDispatcherBusHandler()
	srv, err := dispatcherbus.NewService(d)

	if err != nil {
		return errors.Wrap(err, "Append dispatcherBus Service Fail")
	}

	s.DispatcherBusService = srv
	s.AppendService("dispatcher_bus", srv)

	return nil
}

func (s *Server) initHTTPDService() error {
	p := s.config.DocPath
	exist, err := utils.FileIsExist(p)
	if err != nil || !exist {
		return fmt.Errorf("API File %s is not exist", p)
	}
	d := s.DiagService.NewHTTPDHandler()

	srv := httpd.NewService(p, s.config.HTTP, s.hostname, d, s.DiagService)

	s.HTTPDService = srv

	return nil
}

func (s *Server) initAudiVWProtocolService() error {
	c := s.config.AudiVW
	d := s.DiagService.NewAudiVWHandler()
	srv := audi_vw.NewService(c, d)

	s.AudiVWService = srv

	return nil
}

func (s *Server) appendAudiVWService() {

	s.AudiVWService.Minio = s.MinioService
	s.AudiVWService.Aiis = s.AiisService
	s.AudiVWService.WS = s.WSNotifyService
	s.AudiVWService.DB = s.StorageService
	s.AudiVWService.Odoo = s.OdooService

	s.AppendService("audi/vw", s.AudiVWService)
}

func (s *Server) appendOpenProtocolService() {

	c := s.config.OpenProtocol
	d := s.DiagService.NewOpenProtocolHandler()
	srv := openprotocol.NewService(c, d, vendors.OpenProtocolVendors, s.StorageService, s.OdooService)

	s.OpenprotocolService = srv

	s.AppendService("openprotocol", srv)
}

func (s *Server) appendHTTPDService() {
	s.AppendService("httpd", s.HTTPDService)
}

func (s *Server) appendMinioService() {
	c := s.config.Minio
	d := s.DiagService.NewMinioHandler()
	srv := minio.NewService(c, d, s.StorageService)

	s.MinioService = srv
	s.AppendService("minio", srv)

}

func (s *Server) appendDeviceService() error {
	c := s.config.Device
	d := s.DiagService.NewDeviceHandler()
	srv, err := device.NewService(c, d, s.DispatcherBusService)
	if err != nil {
		return errors.Wrap(err, "append device service fail")
	}

	s.DeviceService = srv
	s.AppendService("device", srv)

	return nil
}

func (s *Server) appendTighteningDeviceService() error {
	c := s.config.TighteningDevice
	d := s.DiagService.NewTighteningDeviceHandler()
	srv, err := tightening_device.NewService(c, d,
		[]tightening_device.ITighteningProtocol{s.OpenprotocolService, s.AudiVWService}, s.DispatcherBusService, s.DeviceService, s.StorageService)

	if err != nil {
		return errors.Wrap(err, "append tightening_device service fail")
	}

	s.TighteningDeviceService = srv
	s.AppendService("tightening_device", srv)

	return nil
}

func (s *Server) appendAiisService() error {
	c := s.config.Aiis
	d := s.DiagService.NewAiisHandler()
	srv := aiis.NewService(c, d, s.DispatcherBusService, s.StorageService, s.TransportService, s.WSNotifyService)

	s.AiisService = srv
	s.AppendService("aiis", srv)

	return nil
}

func (s *Server) appendOdooService() error {
	c := s.config.Odoo
	d := s.DiagService.NewOdooHandler()
	srv := odoo.NewService(c, d, s.DispatcherBusService, s.StorageService, s.HTTPDService)

	s.OdooService = srv
	s.AppendService("odoo", srv)

	return nil
}

func (s *Server) appendWebsocketService() {
	c := s.config.WSNotify
	d := s.DiagService.NewWebsocketHandler()
	srv := wsnotify.NewService(c, d, s.DispatcherBusService, s.HTTPDService)

	s.WSNotifyService = srv
	s.AppendService("websocket", srv)

}

func (s *Server) appendHMIService() error {
	d := s.DiagService.NewHMIHandler()
	srv := hmi.NewService(d, s.DispatcherBusService, s.WSNotifyService, s.HTTPDService, s.OdooService, s.StorageService)
	s.AppendService("hmi", srv)

	return nil
}

func (s *Server) appendStorageService() {
	c := s.config.Storage
	d := s.DiagService.NewStorageHandler()
	srv := storage.NewService(c, d)

	s.StorageService = srv

	s.AppendService("storage", srv)

}

func (s *Server) appendScannerService() error {
	c := s.config.Scanner
	d := s.DiagService.NewScannerHandler()

	srv := scanner.NewService(c, d, s.DispatcherBusService, s.DeviceService)

	s.ScannerService = srv
	s.AppendService("scanner", srv)

	return nil
}

func (s *Server) appendGRPCService() {
	c := s.config.Grpc
	d := s.DiagService.NewGRPCHandler()

	srv := grpc.NewService(c, d)

	s.GRPCService = srv

	s.AppendService(transport.GRPCTransport, srv)
}

func (s *Server) appendBrokerService() {
	c := s.config.Broker
	d := s.DiagService.NewBrokerHandler()

	srv := broker.NewService(c, d, s.DispatcherBusService)

	s.BrokerService = srv

	s.AppendService(transport.BrokerTransport, srv)

}

func (s *Server) appendIOService() error {
	c := s.config.IO
	d := s.DiagService.NewIOHandler()

	srv := io.NewService(c, d, s.DispatcherBusService, s.DeviceService)

	s.IOService = srv
	s.AppendService("io", srv)

	return nil
}

func (s *Server) appendReaderService() {
	c := s.config.Reader
	d := s.DiagService.NewReaderHandler()

	srv := reader.NewService(c, d, s.DeviceService, s.DispatcherBusService)

	s.ReaderService = srv
	s.AppendService("reader", srv)

}

func (s *Server) Open() error {

	if err := s.startServices(); err != nil {
		s.Close()
		return err
	}

	go s.watchServices()

	return nil
}

func (s *Server) startServices() error {
	for _, service := range s.Services {
		s.Diag.Debug("opening service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
		if err := service.Open(); err != nil {
			return fmt.Errorf("open service %T: %s", service, err)
		}
		s.Diag.Debug("opened service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
	}

	return nil
}

// Watch if something dies
func (s *Server) watchServices() {
	var err error
	select {
	case err = <-s.HTTPDService.Err():
	case err = <-s.AudiVWService.Err():
	}
	s.err <- err
}

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

func (s *Server) Reload() {
	return
}

func (s *Server) Close() error {

	if err := s.HTTPDService.Close(); err != nil {
		s.Diag.Error("error closing httpd service", err)
	}

	for i := len(s.Services) - 1; i >= 0; i-- {
		service := s.Services[i]
		s.Diag.Debug("closing service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
		err := service.Close()
		if err != nil {
			s.Diag.Error("error closing service", err, keyvalue.KV("service", fmt.Sprintf("%T", service)))
		}
		s.Diag.Debug("closed service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
	}

	return nil
}
