package diagnostic

import
(
	"github.com/influxdata/kapacitor/services/diagnostic"
	"runtime"
)

// Cmd handler

type CmdHandler struct {
	l diagnostic.Logger
}

func (h *CmdHandler) Error(msg string, err error) {
	h.l.Error(msg, diagnostic.Error(err))
}

func (h *CmdHandler) AiisStarting(version, branch, commit string) {
	h.l.Info("aiis starting", diagnostic.String("version", version), diagnostic.String("branch", branch), diagnostic.String("commit", commit))
}

func (h *CmdHandler) GoVersion() {
	h.l.Info("go version", diagnostic.String("version", runtime.Version()))
}

func (h *CmdHandler) Info(msg string) {
	h.l.Info(msg)
}
