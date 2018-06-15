package rush

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/aiis/services/fis"
	"github.com/masami10/aiis/services/httpd"
	"sync"
	"strings"
)


type Diagnostic interface {
	Error(msg string, err error)
}

type cResult struct {
	r  *OperationResult
	id int64
}

type Service struct {
	HTTPDService   *httpd.Service
	workers        int
	wg             sync.WaitGroup
	chResult       chan cResult
	closing        chan struct{}

	StorageService interface {
		UpdateResults(result *OperationResult, id int64, sent int) error
		UpdateSent(id int64, sent int) error
	}

	Fis *fis.Service

	diag Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable {
		return &Service{
			diag:     d,
			workers:  c.Workers,
			chResult: make(chan cResult, c.Workers),
			closing:  make(chan struct{}),
		}
	}
	return nil
}

func (s *Service) putFisResult(ctx iris.Context) {

	var r OperationResult
	err := ctx.ReadJSON(&r)

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Params from Odoo wrong: %s", err.Error()))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	fis_result := s.OperationToFisResult(&r)
	fis_err := s.Fis.PushResult(&fis_result)
	if fis_err != nil {
		ctx.Writef(fmt.Sprintf("Push fis err: %s", fis_err.Error()))
		ctx.StatusCode(iris.StatusBadRequest)
	} else {
		ctx.StatusCode(iris.StatusOK)
	}

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
		Pattern:     "/operation.results/{result_id:long}",
		HandlerFunc: s.getResultUpdate,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	r_fis := httpd.Route{
		Method:      "PUT",
		Pattern:     "/fis.results",
		HandlerFunc: s.putFisResult,
	}
	s.HTTPDService.Handler[0].AddRoute(r_fis)

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)

		go s.run()

	}
	return nil
}

func (s *Service) run() {
	for {
		select {
		case r := <-s.chResult:
			s.HandleResult(&r)

		case <-s.closing:
			s.wg.Done()
			return
		}
	}
}

func (s *Service) Close() error {

	for i := 0; i < s.workers; i++ {
		s.closing <- struct{}{}
	}

	s.wg.Wait()
	return nil
}

func (s *Service) OperationToFisResult(r *OperationResult) (fis.FisResult) {
	fis_result := fis.FisResult{}
	fis_result.Init()

	fis_result.EquipemntName = r.EquipemntName
	fis_result.FactoryName = r.FactoryName
	fis_result.Year = r.Year
	fis_result.Pin = r.Pin
	fis_result.PinCheckCode = r.Pin_check_code
	fis_result.AssemblyLine = r.AssemblyLine
	fis_result.ResultID = "1"
	fis_result.Lnr = r.Lnr

	value_result := 1

	if strings.ToUpper(r.MeasureResult) == "OK" {
		fis_result.ResultValue = "IO__"
	} else {
		fis_result.ResultValue = "NIO_"
		value_result = 0
	}

	fis_result.Dat = r.ControlDate
	fis_result.SystemType = s.Fis.Config().SystemType
	fis_result.SoftwareVersion = s.Fis.Config().SoftwareVersion
	fis_result.Mode = s.Fis.Config().Mode

	// 扭矩结果
	rv := fis.FisResultValue{}
	rv.Value = r.MeasureTorque
	rv.ID = fis.FIS_ID_NM
	rv.Unit = fis.FIS_UNIT_NM
	rv.Measure = value_result
	fis_result.Values = append(fis_result.Values, rv)

	// 角度结果
	rv.Value = r.MeasureDegree
	rv.ID = fis.FIS_ID_DEG
	rv.Unit = fis.FIS_UNIT_DEG
	rv.Measure = value_result
	fis_result.Values = append(fis_result.Values, rv)

	return fis_result
}

func (s *Service) HandleResult(cr *cResult) {

	// 结果推送fis
	fis_result := s.OperationToFisResult(cr.r)

	sent := 1
	e := s.Fis.PushResult(&fis_result)
	if e != nil {
		sent = 0
		s.diag.Error("push result to fis error", e)
	}

	if len(cr.r.CurObjects) == 0 {
		err := s.StorageService.UpdateSent(cr.id, sent)
		if err != nil {
			s.diag.Error("update sent error", err)
		}
	} else {
		// 结果保存数据库
		err := s.StorageService.UpdateResults(cr.r, cr.id, sent)
		if err != nil {
			s.diag.Error("update result error", err)
		}
	}
}