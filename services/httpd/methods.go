package httpd

import (
	"github.com/kataras/iris"
	"io/ioutil"
)

type Methods struct {
	service	*Service
}

func (m *Methods) getDoc(ctx iris.Context) {
	f, _ := ioutil.ReadFile(m.service.ApiDoc)

	ctx.Header("content-type", "application/json")
	ctx.Write(f)
}
