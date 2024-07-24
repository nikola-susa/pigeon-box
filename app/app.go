package app

import (
	"fmt"
	"github.com/nikola-susa/secret-chat/assets"
	"github.com/nikola-susa/secret-chat/config"
	"github.com/nikola-susa/secret-chat/serverevent"
	"github.com/nikola-susa/secret-chat/slackaction"
	"github.com/nikola-susa/secret-chat/storage"
	"github.com/nikola-susa/secret-chat/store"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"net/http"
)

type App struct {
	Config      *config.Config
	Storage     storage.Storage
	Store       *store.Store
	Event       *serverevent.Server
	SlackApi    *slack.Client
	Socket      *socketmode.Client
	SlackAction *slackaction.SlackAction
}

func New(config *config.Config, storage storage.Storage, store *store.Store, event *serverevent.Server, slackApi *slack.Client, socket *socketmode.Client, slackAction *slackaction.SlackAction) *App {
	return &App{
		Config:      config,
		Storage:     storage,
		Store:       store,
		Event:       event,
		SlackApi:    slackApi,
		Socket:      socket,
		SlackAction: slackAction,
	}
}

func (a *App) Serve() error {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/{thread_id}/{session_token}", a.HandleAuth)
	mux.Handle("DELETE /t/{thread_id}/auth", AuthMiddleware(a)(http.HandlerFunc(a.HandleDeleteSession)))

	mux.Handle("GET /t/{thread_id}", AuthMiddleware(a)(http.HandlerFunc(a.HandleRenderThread)))
	mux.Handle("POST /t/{thread_id}/m", AuthMiddleware(a)(http.HandlerFunc(a.HandleCreateNewMessage)))

	mux.Handle("GET /t/{thread_id}/f/{id}", AuthMiddleware(a)(http.HandlerFunc(a.HandleDownloadFile)))
	mux.Handle("POST /t/{thread_id}/f", AuthMiddleware(a)(http.HandlerFunc(a.HandleCreateFileMessage)))

	mux.Handle("GET /t/{thread_id}/m/{message_id}/edit", AuthMiddleware(a)(http.HandlerFunc(a.HandleRenderEdit)))
	mux.Handle("POST /t/{thread_id}/m/{message_id}/cancel", AuthMiddleware(a)(http.HandlerFunc(a.HandleEditCancel)))

	mux.Handle("DELETE /t/{thread_id}/m/{message_id}", AuthMiddleware(a)(http.HandlerFunc(a.HandleDeleteMessage)))

	mux.HandleFunc("DELETE /file/{id}", a.HandleDeleteFile)

	mux.HandleFunc("GET /not-authenticated", func(w http.ResponseWriter, r *http.Request) {
		RenderError(w)
		return
	})
	mux.Handle("GET /t/{stream}/e", AuthMiddleware(a)(a.Event))

	mux.Handle("GET /static/{path...}", http.StripPrefix("/static/", http.FileServerFS(assets.PublicFS)))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port),
		Handler: mux,
	}

	log.Printf("Server listening on %s:%d", a.Config.Server.Host, a.Config.Server.Port)
	return server.ListenAndServe()
}
