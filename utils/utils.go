package utils

import (
	"encoding/base64"
	"github.com/kataras/iris/core/errors"
	"github.com/satori/go.uuid"
	"os"
	"strings"
	"sync"
	"time"
)

func FileIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AppendByteSlice(s []byte, t []byte) []byte {
	zlen := len(s) + len(t)
	z := make([]byte, zlen)
	copy(z, s)
	copy(z[len(s):], t)
	return z
}

func GenerateID() string {
	u4 := uuid.NewV4()
	return base64.RawURLEncoding.EncodeToString(u4.Bytes())
}

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}

func RushRound(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

func ReverseString(raw string) string {
	rt := ""
	for _, v := range raw {
		rt = string(v) + rt
	}
	return rt
}

func ArrayContains(s []int, e int) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}

func WaitGroupTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	if wg == nil {
		return errors.New("Wait Group Is Nil")
	}

	wg.Add(1)
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return errors.New("Timeout") // timed out
	}
}
