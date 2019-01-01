package diagnostic

import (
	"errors"
	"fmt"
	"github.com/masami10/aiis/keyvalue"
	"log"
	"runtime"
	"time"
)

type logLevel int

const (
	llInvalid logLevel = iota
	llDebug
	llError
	llInfo
)

type StaticLevelHandler struct {
	l     Logger
	level logLevel
}

func (h *StaticLevelHandler) Write(buf []byte) (int, error) {
	switch h.level {
	case llDebug:
		h.l.Debug(string(buf))
	case llError:
		h.l.Error(string(buf))
	case llInfo:
		h.l.Info(string(buf))
	default:
		return 0, errors.New("invalid log level")
	}

	return len(buf), nil
}

func Err(l Logger, msg string, err error, ctx []keyvalue.T) {
	if len(ctx) == 0 {
		l.Error(msg, Error(err))
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Error(msg, Error(err), String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Error(msg, Error(err), String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx)+1) // +1 for error
	fields[0] = Error(err)
	for i := 1; i < len(fields); i++ {
		kv := ctx[i-1]
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Error(msg, fields...)
}

func Info(l Logger, msg string, ctx []keyvalue.T) {
	if len(ctx) == 0 {
		l.Info(msg)
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Info(msg, String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Info(msg, String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx))
	for i, kv := range ctx {
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Info(msg, fields...)
}

func Debug(l Logger, msg string, ctx []keyvalue.T) {
	if len(ctx) == 0 {
		l.Debug(msg)
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Debug(msg, String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Debug(msg, String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx))
	for i, kv := range ctx {
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Debug(msg, fields...)
}

// Cmd handler

type CmdHandler struct {
	l Logger
}

func (h *CmdHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *CmdHandler) AiisStarting(version, branch, commit string) {
	h.l.Info("aiis starting", String("version", version), String("branch", branch), String("commit", commit))
}

func (h *CmdHandler) GoVersion() {
	h.l.Info("go version", String("version", runtime.Version()))
}

func (h *CmdHandler) Info(msg string) {
	h.l.Info(msg)
}

type ServerHandler struct {
	l Logger
}

func (h *ServerHandler) Error(msg string, err error, ctx ...keyvalue.T) {
	Err(h.l, msg, err, ctx)
}

func (h *ServerHandler) Info(msg string, ctx ...keyvalue.T) {
	Info(h.l, msg, ctx)
}

func (h *ServerHandler) Debug(msg string, ctx ...keyvalue.T) {
	Debug(h.l, msg, ctx)
}

// HTTPD handler

type HTTPDHandler struct {
	l Logger
}

func (h *HTTPDHandler) NewHTTPServerErrorLogger() *log.Logger {
	s := &StaticLevelHandler{
		l:     h.l.With(String("service", "httpd_server_errors")),
		level: llError,
	}

	return log.New(s, "", log.LstdFlags)
}

func (h *HTTPDHandler) StartingService() {
	h.l.Info("starting HTTP service")
}

func (h *HTTPDHandler) StoppedService() {
	h.l.Info("closed HTTP service")
}

func (h *HTTPDHandler) ShutdownTimeout() {
	h.l.Error("shutdown timedout, forcefully closing all remaining connections")
}

func (h *HTTPDHandler) AuthenticationEnabled(enabled bool) {
	h.l.Info("authentication", Bool("enabled", enabled))
}

func (h *HTTPDHandler) ListeningOn(addr string, proto string) {
	h.l.Info("listening on", String("addr", addr), String("protocol", proto))
}

func (h *HTTPDHandler) WriteBodyReceived(body string) {
	h.l.Debug("write body received by handler: %s", String("body", body))
}

func (h *HTTPDHandler) HTTP(
	host string,
	username string,
	start time.Time,
	method string,
	uri string,
	proto string,
	status int,
	referer string,
	userAgent string,
	reqID string,
	duration time.Duration,
) {
	h.l.Info("http request",
		String("host", host),
		String("username", username),
		Time("start", start),
		String("method", method),
		String("uri", uri),
		String("protocol", proto),
		Int("status", status),
		String("referer", referer),
		String("user-agent", userAgent),
		String("request-id", reqID),
		Duration("duration", duration),
	)
}

func (h *HTTPDHandler) RecoveryError(
	msg string,
	err string,
	host string,
	username string,
	start time.Time,
	method string,
	uri string,
	proto string,
	status int,
	referer string,
	userAgent string,
	reqID string,
	duration time.Duration,
) {
	h.l.Error(
		msg,
		String("err", err),
		String("host", host),
		String("username", username),
		Time("start", start),
		String("method", method),
		String("uri", uri),
		String("protocol", proto),
		Int("status", status),
		String("referer", referer),
		String("user-agent", userAgent),
		String("request-id", reqID),
		Duration("duration", duration),
	)
}

func (h *HTTPDHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

// PMON Handler

type PmonHandler struct {
	l Logger
}

func (h *PmonHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *PmonHandler) Debug(msg string) {
	h.l.Debug(msg)
}

// Odoo Handler

type OdooHandler struct {
	l Logger
}

func (h *OdooHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

// Storage Handler

type StorageHandler struct {
	l Logger
}

func (h *StorageHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *StorageHandler) OpenEngineSuccess(msg string) {
	h.l.Info(fmt.Sprintf("Open Engine Success: %s ", msg))
}

func (h *StorageHandler) UpdateResultSuccess(id int64) {
	h.l.Debug(fmt.Sprintf("Update result Success: %d ", id))
}

// Rush Handler

type RushHandler struct {
	l Logger
}

func (h *RushHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *RushHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *RushHandler) Debug(msg string) {
	h.l.Debug(msg)
}

// Fis Handler
type FisHandler struct {
	l Logger
}

func (h *FisHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *FisHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type MasterPLCHandler struct {
	l Logger
}

func (h *MasterPLCHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *MasterPLCHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type WsHandler struct {
	l Logger
}

type ChanganHandler struct {
	l Logger
}

func (h *ChanganHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *ChanganHandler) ReciveNewTask(msg string) {
	h.l.Info(fmt.Sprintf("ReciveNewTask from Andon: %s ", msg))
}

func (h *WsHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *WsHandler) Disconnect(id string) {
	h.l.Info("ws Connection disconnected", String("ID", id))
}

func (h *WsHandler) Close() {
	h.l.Info("ws server closing")
}

func (h *WsHandler) OnMessage(msg string) {
	h.l.Debug(msg)
}

func (h *WsHandler) Closed() {
	h.l.Info("ws server closed")
}

type MinioHandler struct {
	l Logger
}

func (h *MinioHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}
