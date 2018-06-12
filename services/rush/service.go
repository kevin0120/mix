package rush

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/httpd"
	"sync"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type cResult struct {
	r 		*OperationResult
	id 		int64
}

type Service struct {
	HTTPDService *httpd.Service
	workers			int
	wg 				sync.WaitGroup
	chResult		chan cResult
	closing 		chan struct{}
	StorageService interface {
		UpdateResults(result *OperationResult, id int64) error
	}

	diag 		Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable {
		return &Service{
			diag: 			d,
			workers: 		c.Workers,
			chResult:		make(chan cResult, c.Workers),
			closing: 		make(chan struct{}),
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
	cr := cResult{r: &r, id: resultId}

	s.chResult <- cr

	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) Open() error {

	r := httpd.Route{
		Method:      "PUT",
		Pattern:     "/aiis/v1/operation.results/{result_id:long}",
		HandlerFunc: s.getResultUpdate,
	}
	s.HTTPDService.Handler[0].AddRoute(r)
	for i := 0; i < s.workers; i ++ {
		s.wg.Add(1)

		go s.run()

	}
	return nil
}

func (s *Service) run() {
	for {
		select {
		case r := <- s.chResult:
			err :=s.StorageService.UpdateResults(r.r, r.id)
			if err != nil {
				s.diag.Error("update result error ",err)
			}

		case <-s.closing:
			s.wg.Done()
			return
		}
	}
}

func (s *Service) Close() error {

	for i := 0; i < s.workers; i ++ {
		s.closing<- struct{}{}
	}

	s.wg.Wait()
	return nil
}
