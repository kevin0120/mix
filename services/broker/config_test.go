package broker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestConfig() Config {
	c := NewConfig()
	c.Enable = true
	c.Provider = "nats"
	c.ConnectUrls = []string{"nats://192.168.4.203:4222","nats://192.168.4.203:4223"}
	return c
}

func TestConfig_Validate(t *testing.T) {
	c := newTestConfig()
	err := c.Validate()
	assert.Nil(t, err)
}
