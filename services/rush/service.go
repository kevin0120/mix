package rush

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/httpd"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	HTTPDService *httpd.Service

	StorageService interface {
		UpdateResults(result *OperationResult, id int64) error
	}

	diag Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable {
		return &Service{
			diag: d,
		}
	}
	return nil
}

func (s *Service) getResultUpdate(ctx iris.Context) {
	resultId, err := ctx.Params().GetInt64("result_id")

	if err != nil {
		ctx.Writef("error while trying to parse resultId parameter")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	var r OperationResult
	err = ctx.ReadJSON(&r)

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Params from Rush wrong: %s", err))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	err = s.StorageService.UpdateResults(&r, resultId)

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Update to Database wrong: %s", err))
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) Open() error {

	r := httpd.Route{
		Method:      "PUT",
		Pattern:     "/aiis/v1/operation.results/{result_id:long}",
		HandlerFunc: s.getResultUpdate,
	}
	s.HTTPDService.Handler[0].AddRoute(r)
	return nil
}

func (s *Service) Close() error {
	return nil
}
