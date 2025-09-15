package healthHandler

import (
	"net/http"

	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/server/http/protocol"
)

type httpHandler struct {
	config  config.Config
	version string
}

func newHTTPHandler(config config.Config, version string) *httpHandler {
	return &httpHandler{
		config:  config,
		version: version,
	}
}

// healthCheckHandler godoc
//
//	@Summary		Health check endpoint
//	@Description	Returns the current health status, environment, and version of the API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Health status information"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/health [get]
func (h *httpHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     h.config.Env,
		"version": h.version,
	}

	if err := protocol.JsonResponse(w, http.StatusOK, data); err != nil {
		protocol.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
