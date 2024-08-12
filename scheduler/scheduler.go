package scheduler

import (
	"context"
	"github.com/nikola-susa/pigeon-box/config"
	"github.com/nikola-susa/pigeon-box/serverevent"
	"github.com/nikola-susa/pigeon-box/storage"
	"github.com/nikola-susa/pigeon-box/store"
)

type Worker struct {
	Store   *store.Store
	Event   *serverevent.Server
	Storage storage.Storage
	Config  config.Config
	Ctx     context.Context
}

func New(store *store.Store, event *serverevent.Server, storage storage.Storage, config *config.Config, ctx context.Context) *Worker {
	return &Worker{
		Store:   store,
		Event:   event,
		Storage: storage,
		Config:  *config,
		Ctx:     ctx,
	}
}

func (w *Worker) Run() {
	go w.clearSessions()
	go w.clearMessages()
	go w.clearThread()
}
