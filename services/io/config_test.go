package io

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Validate()
	assert.Nil(t, err)

	cfg.IOS[0].Model = "wef"
	err = cfg.Validate()
	assert.NotNil(t, err)
}
