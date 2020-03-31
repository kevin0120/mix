package broker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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

func newTestService() *Service {
	diag := &TestDiag{}
	c := newTestConfig()
	d := NewService(c, diag, nil)
	return d
}

func TestNewService(t *testing.T) {
	s := newTestService()
	assert.NotNil(t, s)
}

func testHandler(message *Message) ([]byte, error) {
	fmt.Println(fmt.Sprintf("Msg: %v", message))
	return nil, nil
}

func TestService_Subscribe_Fail(t *testing.T) {
	s := newTestService()
	err := s.Subscribe("ttt", testHandler)
	assert.NotNil(t, err)
}

func TestService_Subscribe_Success(t *testing.T) {
	s := newTestService()
	err := s.Open()
	assert.Nil(t, err)
	err = s.Subscribe("ttt", testHandler)
	assert.Nil(t, err)
}

func TestService_Publish(t *testing.T) {

}
