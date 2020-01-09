package broker

import "time"

type DefaultBroker struct {
}

func NewDefaultBroker() *DefaultBroker {
	return &DefaultBroker{}
}

func (s *DefaultBroker) Connect(urls []string) error {
	return nil
}

func (s *DefaultBroker) Address() string {
	return ""
}

func (s *DefaultBroker) Close() error {
	return nil
}

func (s *DefaultBroker) Subscribe(subject string, handler SubscribeHandler) error {
	return nil
}

func (s *DefaultBroker) UnSubscribe(subject string) error {
	return nil
}

func (s *DefaultBroker) Publish(subject string, data []byte) error {
	return nil
}

func (s *DefaultBroker) DoRequest(subject string, data []byte, timeOut time.Duration) (resp []byte, err error) {
	return
}

func (s *DefaultBroker) SetStatusHandler(handler StatusHandler) {}
