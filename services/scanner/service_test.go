package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewSrv(label string) *Service{
	c := NewConfig()
	diag := &TestDiag{}
	c.EntityLabel = label
	s := NewService(c, diag)
	return s
}

func TestNewService(t *testing.T) {
	s := NewSrv("COM5")
	assert.NotNil(t, s)
}

func TestService_search(t *testing.T) {
	s := NewSrv("COM5")
	assert.NotNil(t, s)
	s.search()
}

func TestService_searchWrongLabel(t *testing.T) {
	s := NewSrv("COM4")
	assert.NotNil(t, s)
	s.search()
}
