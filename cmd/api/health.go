package main

import (
	"net/http"
)

// healthCheckHandler returns the health status of the API
//
//	@Summary		Health check endpoint
//	@Description	Returns the current health status, environment, and version of the API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Health status information"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": VERSION,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
