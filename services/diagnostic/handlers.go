package diagnostic

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/masami10/aiis/keyvalue"
	"github.com/masami10/aiis/models"
	"github.com/masami10/aiis/uuid"
)

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

// Aiis Handler

type AiisHandler struct {
	l Logger
}

func (h *AiisHandler) Start(task string) {
	h.l.Info("started Aiis")
}

func (h *AiisHandler) Stop(task string) {
	h.l.Info("stopped Aiis")
}


func (h *AiisHandler) Error(msg string, err error, ctx ...keyvalue.T) {
	Err(h.l, msg, err, ctx)
}

func TagPairs(tags models.Tags) []Field {
	ts := []Field{}
	for k, v := range tags {
		ts = append(ts, String(k, v))
	}

	return ts
}

func FieldPairs(tags models.Fields) []Field {
	ts := []Field{}
	for k, v := range tags {
		var el Field
		switch t := v.(type) {
		case int64:
			el = Int64(k, t)
		case string:
			el = String(k, t)
		case float64:
			el = Float64(k, t)
		case bool:
			el = Bool(k, t)
		default:
			el = String(k, fmt.Sprintf("%v", t))
		}
		ts = append(ts, el)
	}

	return ts
}


// Session handler

type SessionHandler struct {
	l Logger
}

func (h *SessionHandler) CreatedLogSession(id uuid.UUID, contentType string, tags []tag) {
	ts := make([]string, len(tags))
	for i, t := range tags {
		ts[i] = t.key + "=" + t.value
	}

	h.l.Info("created log session", Stringer("id", id), String("content-type", contentType), Strings("tags", ts))
}

func (h *SessionHandler) DeletedLogSession(id uuid.UUID, contentType string, tags []tag) {
	ts := make([]string, len(tags))
	for i, t := range tags {
		ts[i] = t.key + "=" + t.value
	}

	h.l.Info("deleted log session", Stringer("id", id), String("content-type", contentType), Strings("tags", ts))
}

// Command handler

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

// Load handler
type LoadHandler struct {
	l Logger
}

func (h *LoadHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *LoadHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *LoadHandler) Loading(el string, file string) {
	h.l.Debug("loading object from file", String("object", el), String("file", file))
}


// Server handler
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