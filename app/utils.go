package app

import (
	"encoding/json"
	"github.com/nikola-susa/secret-chat/templates"
	"golang.org/x/net/context"
	"net/http"
)

type resource struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func Resource(w http.ResponseWriter, code int, v interface{}) {
	payload := new(resource)
	payload.Status = "SUCCESS"
	payload.Data = v
	if code >= http.StatusBadRequest {
		payload.Status = "ERROR"
	} else if code >= http.StatusInternalServerError {
		payload.Status = "FAILURE"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderError(w http.ResponseWriter) {
	errorPage := templates.ErrorPage()
	err := errorPage.Render(context.Background(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
