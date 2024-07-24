package slackaction

import (
	"github.com/nikola-susa/secret-chat/config"
	"github.com/nikola-susa/secret-chat/serverevent"
	"github.com/nikola-susa/secret-chat/storage"
	"github.com/nikola-susa/secret-chat/store"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackAction struct {
	Api     *slack.Client
	Socket  *socketmode.Client
	Store   *store.Store
	Storage *storage.Storage
	Config  *config.Config
	Event   *serverevent.Server
}

func New(api *slack.Client, socket *socketmode.Client, store *store.Store, storage *storage.Storage, config *config.Config, events *serverevent.Server) *SlackAction {
	return &SlackAction{
		Api:     api,
		Socket:  socket,
		Store:   store,
		Storage: storage,
		Config:  config,
		Event:   events,
	}
}
