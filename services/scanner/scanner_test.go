package scanner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	//"time"
)

type TestDiag struct {
}

func (td *TestDiag) Info(msg string) {
	fmt.Printf("info：%s", msg)
}

func (td *TestDiag) Debug(msg string) {
	fmt.Printf("debug：%s", msg)
}

func (td *TestDiag) Error(msg string, err error) {
	fmt.Printf("error: %s", err.Error())
}

func TestNewDevice(t *testing.T) {
	diag := &TestDiag{}
	d := NewDevice("COM5", diag)
	assert.NotNil(t, d)
}
