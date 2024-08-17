package app

import (
	"encoding/json"
	"github.com/nikola-susa/pigeon-box/templates"
	"golang.org/x/net/context"
	"net/http"
	"time"
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

func HTMXRedirect(w http.ResponseWriter, r *http.Request, url string) {
	if IsXHRRequest(r) {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func ConvertTimeToUserRegion(r *http.Request, t string) (*time.Time, error) {

	region := r.Header.Get("X-Timezone")
	if region == "" {
		region = "UTC"
	}

	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(region)
	if err != nil {
		return nil, err
	}

	tt := parsedTime.In(loc)

	return &tt, nil
}

func HTMXEventRedirect(w http.ResponseWriter, r *http.Request, a *App) {

	if IsSSERequest(r) {
		//threadID := r.PathValue("stream")
		//if threadID == "" {
		//	fmt.Print("no sse thread id")
		//	return
		//}
		//
		//userID := r.PathValue("user")
		//if userID == "" {
		//	fmt.Print("no sse user id")
		//	return
		//}
		//
		//eventName := "logout:" + threadID
		//a.Event.Broadcast(threadID, []byte(""), &eventName, &userID, nil)
	}
}

func IsSSERequest(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/event-stream"
}

func IsXHRRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}
