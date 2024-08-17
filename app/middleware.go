package app

import (
	"fmt"
	"github.com/nikola-susa/pigeon-box/crypt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type contextKey string

func AuthMiddleware(a *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			c, err := r.Cookie("session_token")
			if err != nil {
				log.Println("cookie error", err)
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			if c == nil {
				log.Println("cookie is nil")
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			token, err := crypt.HashIDDecodeInt(c.Value, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
			if err != nil {
				log.Println("hash decode error", err)
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			session, err := a.Store.GetSession(token)
			if err != nil {
				log.Println("get session error", err)
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			if session == nil {
				log.Println("session not found")
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			err = a.Store.UpdateSessionExpiresAt(session.ID, time.Now().Add(24*time.Hour))
			if err != nil {
				log.Println("update session error", err)
				HTMXRedirect(w, r, "/not-authenticated")
				return
			}

			hashedThreadID, err := crypt.HashIDEncodeInt(session.ThreadID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    c.Value,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     fmt.Sprintf("/t/%s", hashedThreadID),
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
			})

			ctx := context.WithValue(r.Context(), contextKey("user_id"), session.UserID)
			ctx = context.WithValue(ctx, contextKey("thread_id"), session.ThreadID)
			ctx = context.WithValue(ctx, contextKey("thread_hash"), hashedThreadID)
			ctx = context.WithValue(ctx, contextKey("session_id"), session.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func EventMiddleware(a *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			c, err := r.Cookie("session_token")
			if err != nil {
				HTMXEventRedirect(w, r, a)
				return
			}

			token, err := crypt.HashIDDecodeInt(c.Value, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
			if err != nil {
				HTMXEventRedirect(w, r, a)
				return
			}

			session, err := a.Store.GetSession(token)
			if err != nil {
				fmt.Println("token", token)
				HTMXEventRedirect(w, r, a)
				return
			}

			if session == nil {
				HTMXEventRedirect(w, r, a)
				return
			}

			hashedThreadID, err := crypt.HashIDEncodeInt(session.ThreadID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

			ctx := context.WithValue(r.Context(), contextKey("user_id"), session.UserID)
			ctx = context.WithValue(ctx, contextKey("thread_id"), session.ThreadID)
			ctx = context.WithValue(ctx, contextKey("thread_hash"), hashedThreadID)
			ctx = context.WithValue(ctx, contextKey("session_id"), session.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
