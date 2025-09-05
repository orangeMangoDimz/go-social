package main

import "net/http"
import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": VERSION,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
