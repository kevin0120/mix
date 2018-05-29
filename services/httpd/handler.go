package httpd



type Route struct {
	Method      string
	Pattern     string
	HandlerFunc interface{}
}

type Handler struct {
	Version string
	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}



	diag Diagnostic
	// Detailed logging of write path
	// Uses normal logger
	writeTrace bool

	// Log every HTTP access.
	loggingEnabled bool
}

func (h *Handler) AddRoutes(routes []Route) error {
	for _, r := range routes {
		err := h.addRawRoute(r)
		if err != nil {
			return err
		}
	}
	return nil
}
