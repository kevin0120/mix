package ts002

import (
	"encoding/json"
	"github.com/kataras/iris"
)

// Mes->Rush 报警控制请求
func (s *Service) putAlarmReq(ctx iris.Context) {
	var req RushAlarmReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	err = s.alarmControl(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

// Mes->Rush PSet下发请求
func (s *Service) putPSetReq(ctx iris.Context) {
	var req RushPSetReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	err = s.psetControl(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

// Mes->Rush IO输出控制
func (s *Service) putIOReq(ctx iris.Context) {
	var req RushIOControlReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	err = s.ioControl(&req)
	if err != nil {
		ErrResponse(ctx, iris.StatusBadRequest, string(iris.StatusBadRequest), err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func ErrResponse(ctx iris.Context, respCode int, errCode string, errMsg string) {
	ctx.StatusCode(respCode)
	errResp := RushErrResp{
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
	}

	body, _ := json.Marshal(errResp)
	_, _ = ctx.Write(body)
}
