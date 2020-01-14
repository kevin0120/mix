package broker

import (
	"github.com/masami10/rush/services/transport"
	"github.com/nats-io/nats.go"
	"time"
)

var StatusBroker = map[nats.Status]string{
	nats.DISCONNECTED:  "DISCONNECTED",
	nats.CONNECTED:     "CONNECTED",
	nats.CLOSED:        "CLOSED",
	nats.RECONNECTING:  "RECONNECTING",
	nats.CONNECTING:    "CONNECTING",
	nats.DRAINING_SUBS: "DRAINING_SUBS",
	nats.DRAINING_PUBS: "DRAINING_PUBS",
}

type StatusHandler func(string)

type IBrokerProvider interface {
	Address() string
	Connect(urls []string) error
	Close() error
	Subscribe(subject string, handler transport.OnMsgHandler) error
	UnSubscribe(subject string) error
	Publish(subject string, data []byte) error
	DoRequest(subject string, data []byte, timeOut time.Duration) ([]byte, error)
	SetStatusHandler(handler StatusHandler)
}
