package app

import (
	"fmt"
	"github.com/nikola-susa/secret-chat/crypt"
	"log"
	"net/http"
	"time"
)

func (a *App) HandleAuth(w http.ResponseWriter, r *http.Request) {
	threadId, err := crypt.HashIDDecodeInt(r.PathValue("thread_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error parsing thread id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionToken, err := crypt.HashIDDecodeInt(r.PathValue("session_token"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding session token: %s", err)
		RenderError(w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := a.Store.GetSession(sessionToken)
	if err != nil {
		log.Printf("Error getting session by token: %s", err)
		RenderError(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session == nil {
		log.Printf("Error session not found: %s", sessionToken)
		RenderError(w)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		log.Printf("Error session expired: %s", sessionToken)
		RenderError(w)
		http.Error(w, "session expired", http.StatusUnauthorized)
		return
	}

	if session.ThreadID != threadId {
		log.Printf("Error session thread id mismatch: %d != %d", session.ThreadID, threadId)
		RenderError(w)
		http.Error(w, "session thread id mismatch", http.StatusUnauthorized)
		return
	}

	err = a.Store.UpdateSessionExpiresAt(session.ID, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		RenderError(w)
		return
	}

	user, err := a.Store.GetUser(session.UserID)
	if err != nil {
		RenderError(w)
		log.Printf("Error getting user by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		RenderError(w)
		log.Printf("Error user not found: %d", session.UserID)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	thread, err := a.Store.GetThread(threadId)
	if err != nil {
		RenderError(w)
		log.Printf("Error getting thread by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		RenderError(w)
		log.Printf("Error thread not found: %d", threadId)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	hashedThreadID, err := crypt.HashIDEncodeInt(threadId, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   r.PathValue("session_token"),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    fmt.Sprintf("/t/%s", hashedThreadID),
	})

	redirectURL := fmt.Sprintf("/t/%s", hashedThreadID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (a *App) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {

	sessionId := r.Context().Value(contextKey("session_id")).(int)

	err := a.Store.DeleteSession(sessionId)
	if err != nil {
		log.Printf("Error deleting session: %s", err)
		RenderError(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-24 * time.Hour),
		Path:    fmt.Sprintf("/t/%s", r.PathValue("thread_id")),
	})

	w.Header().Set("HX-Redirect", "/not-authenticated")

	w.WriteHeader(http.StatusOK)
}
