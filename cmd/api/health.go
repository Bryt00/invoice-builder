package main

import (
	"net/http"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"status": "available",
		"system_info": Envelope{
			"environment": app.config.env,
			"version":     "1.0.0",
		},
	}

	err := app.apiHandler.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
