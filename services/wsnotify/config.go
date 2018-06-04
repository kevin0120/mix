package wsnotify

import "fmt"

const (
	WS_ROUTE = "/rush/v1/ws"
)

type Config struct {
	Route 			string `yaml:"route"`
	ReadBufferSize  int 	`yaml:"read_buf_size"`
	WriteBufferSize int 	`yaml:"write_buf_size"`
}

func NewConfig() Config {
	c := Config{
		Route: WS_ROUTE,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return c
}

func (c Config) Validate() error {
	route := c.Route
	if route == "" {
		return fmt.Errorf("websocket server route can not be empty")
	}
	if len(route) > 0 && route[0] != '/' {
		return fmt.Errorf("route patterns must begin with a '/' %s", route)
	}
	return nil
}
