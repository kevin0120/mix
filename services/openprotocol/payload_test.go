package openprotocol

import "testing"

func Test_parseOpenProtocolErrorCode(t *testing.T) {
	err := "111"
	ss := parseOpenProtocolErrorCode(err)
	t.Log(ss)
}
