package aiis

import "github.com/masami10/rush/services/tightening_device"

type TighteningResultHandler = func(data *tightening_device.TighteningResult)

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
}

func NewEndpoint(url string, headers map[string]string, method string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
	}
}
