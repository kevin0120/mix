package reader

import (
	"encoding/json"
	"fmt"
	"github.com/ebfe/scard"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
	"sync/atomic"
	"time"
)

const (
	SEARCH_ITV = 1 * time.Second
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value

	diag          Diagnostic
	ctx           *scard.Context
	WS            *wsnotify.Service
	DeviceService *device.Service
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
	}

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

	go s.initReader()
	return nil
}

func (s *Service) initReader() {
	for {
		var err error
		s.ctx, err = scard.EstablishContext()
		if err != nil {
			s.diag.Error("init reader failed", err)
			continue
			time.Sleep(1 * time.Second)
		}
		break
	}

	go s.search()

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

func (s *Service) DeviceType(sn string) string {
	return "reader"
}

func (s *Service) Children() []string {
	return []string{}
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
	barcode, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_READER_UID,
		Data: ReaderUID{
			UID: uid,
		},
	})

	s.WS.WSSendReader(string(barcode))
}

func (s *Service) search() {
	var err error

	for {
		s.ctx, err = scard.EstablishContext()
		if err == nil {
			break
		} else {
			time.Sleep(SEARCH_ITV)
		}
	}

	for {
		// List available readers
		ctx, err := scard.EstablishContext()
		if err != nil {
			time.Sleep(SEARCH_ITV)
			continue
		}

		readers, err := ctx.ListReaders()
		if err != nil {
			//s.diag.Debug("reader lost")
			_ = ctx.Release()
			time.Sleep(SEARCH_ITV)
			continue
		}

		for {
			if len(readers) > 0 {
				index, err := s.waitUntilCardPresent(ctx, readers)
				if err != nil {
					break
				}

				// Connect to card
				card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
				if err != nil {
					time.Sleep(SEARCH_ITV)
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

					time.Sleep(SEARCH_ITV)
				}
			}
		}
	}
}
