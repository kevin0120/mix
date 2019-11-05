package broker

import (
	"context"
	"crypto/tls"
)

type IBrokerProvider interface {
	Address() string
	Connect(urls []string) error
	Close() error
	Subscribe(subject string, handler SubscribeHandler) error
	UnSubscribe(subject string) error
	Publish(subject string, data []byte) error
}

type SubscribeHandler func(*brokerMessage)

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
