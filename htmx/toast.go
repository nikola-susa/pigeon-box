package htmx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SUCCESS = "success"
	INFO    = "info"
	WARNING = "warning"
	ERROR   = "error"
)

type Toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func NewToast(level string, message string) Toast {
	return Toast{Level: level, Message: message}
}

func SuccessToast(w http.ResponseWriter, message string) {
	NewToast(SUCCESS, message).Output(w)
}

func InfoToast(w http.ResponseWriter, message string) {
	NewToast(INFO, message).Output(w)
}

func WarningToast(w http.ResponseWriter, message string) {
	NewToast(WARNING, message).Output(w)
}

func ErrorToast(w http.ResponseWriter, message string) {
	NewToast(ERROR, message).Output(w)
}

func (t Toast) Error() string {
	return fmt.Sprintf("%s", t.Message)
}

func (t Toast) Output(w http.ResponseWriter) {
	t.Message = t.Error()
	eventMap := map[string]Toast{}
	eventMap["showToast"] = t
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", string(jsonData))
}
