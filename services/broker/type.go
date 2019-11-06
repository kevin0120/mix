package broker

import (
	"context"
	"crypto/tls"
	"time"
)

type IBrokerProvider interface {
	Address() string
	Connect(urls []string) error
	Close() error
	Subscribe(subject string, handler SubscribeHandler) error
	UnSubscribe(subject string) error
	Publish(subject string, data []byte) error
	DoRequest(subject string, data []byte, timeOut time.Duration) ([]byte, error)
}

type SubscribeHandler func(*brokerMessage) ([]byte, error)

type brokerOptions struct {
	Addrs     []string
	Secure    bool
	TLSConfig *tls.Config
	Context   context.Context
}

type brokerMessage struct {
	Header map[string]string
	Body   []byte
}
