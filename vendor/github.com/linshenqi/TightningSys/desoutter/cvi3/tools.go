package cvi3

import (
	"github.com/satori/go.uuid"
	"encoding/base64"
	"strings"
	"time"
)

func GenerateID() (string) {
	u4, _ := uuid.NewV4()
	return base64.RawURLEncoding.EncodeToString(u4.Bytes())
}

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}