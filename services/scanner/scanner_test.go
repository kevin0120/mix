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
	fmt.Printf("info：%s\n", msg)
}

func (td *TestDiag) Debug(msg string) {
	fmt.Printf("debug：%s\n", msg)
}

func (td *TestDiag) Error(msg string, err error) {
	fmt.Printf("msg: %s, error: %s\n", msg, err.Error())
}

func TestNewDevice(t *testing.T) {
	diag := &TestDiag{}
	d := NewDevice("COM5", diag)
	assert.NotNil(t, d)
}
