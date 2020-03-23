package reader

import (
	"fmt"
	"github.com/ebfe/scard"
	"github.com/masami10/rush/services/device"
	dispatcherBus "github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
	"github.com/satori/go.uuid"
	"go.uber.org/atomic"
	"time"
)

const (
	SearchItv = 1 * time.Second
)

type Service struct {
	device.BaseDevice
	configValue   atomic.Value
	diag          Diagnostic
	ctx           *scard.Context
	deviceService IDeviceService
	dispatcherBus Dispatcher
}

func NewService(c Config, d Diagnostic, ds IDeviceService, dp Dispatcher) *Service {

	s := &Service{
		diag:          d,
		deviceService: ds,
		dispatcherBus: dp,
	}
	s.SetSerialNumber(uuid.NewV4().String())
	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	if err := s.dispatcherBus.Create(dispatcherBus.DispatcherReaderData, utils.DefaultDispatcherBufLen); err != nil {
		s.diag.Error("Create DispatcherReaderData Failed", err)
	}

	err := s.dispatcherBus.Start(dispatcherBus.DispatcherReaderData)
	if err != nil {
		return err
	}

	go s.search()

	s.deviceService.AddDevice("reader", s)

	return nil
}

func (s *Service) Close() error {
	if s.ctx != nil {
		err := s.ctx.Release()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func (s *Service) DeviceType() string {
	return "reader"
}

func (s *Service) Children() map[string]device.IBaseDevice {
	return map[string]device.IBaseDevice{}
}

func (s *Service) Status() string {
	return "online"
}

func (s *Service) Config() interface{} {
	return nil
}

func (s *Service) Data() interface{} {
	return nil
}

func (s *Service) notifyUID(uid string) {
	// 分发读卡器数据
	s.doDispatch(dispatcherBus.DispatcherReaderData, uid)
}

func (s *Service) search() {
	var err error

	for {
		s.ctx, err = scard.EstablishContext()
		if err == nil {
			break
		} else {
			time.Sleep(SearchItv)
		}
	}

	for {
		// List available readers
		ctx, err := scard.EstablishContext()
		if err != nil {
			time.Sleep(SearchItv)
			continue
		}

		readers, err := ctx.ListReaders()
		if err != nil {
			//s.diag.Debug("reader lost")
			_ = ctx.Release()
			time.Sleep(SearchItv)
			continue
		}

		for {
			if len(readers) > 0 {
				index, err := s.waitUntilCardPresent(ctx, readers)
				if err != nil {
					break
				}

				// connect to card
				card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
				if err != nil {
					time.Sleep(SearchItv)
					continue
				}

				for {
					var cmd = []byte{0xff, 0xca, 0x00, 0x00, 0x00} // SELECT uid

					_, err := card.Status()
					if err != nil {
						_ = card.Disconnect(scard.ResetCard)
						break
					}

					rsp, err := card.Transmit(cmd)
					if err != nil {
						// card lost
						//s.diag.Debug("card lost")
						_ = card.Disconnect(scard.ResetCard)
						break
					}

					uid := fmt.Sprintf("%x", string(rsp))
					s.diag.Debug(fmt.Sprintf("uid:%s", uid))

					// ws notify
					s.notifyUID(uid)

					_ = card.Disconnect(scard.ResetCard)

					time.Sleep(SearchItv)
				}
			}
		}
	}
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}
