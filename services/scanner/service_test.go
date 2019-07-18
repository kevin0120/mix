package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewSrv() *Service{
	c := NewConfig()
	diag := &TestDiag{}
	c.EntityLabel = "COM5"
	s := NewService(c, diag)
	return s
}

func TestNewService(t *testing.T) {
	s := NewSrv()
	assert.NotNil(t, s)
}

func TestService_search(t *testing.T) {
	s := NewSrv()
	assert.NotNil(t, s)
	s.search()
}
