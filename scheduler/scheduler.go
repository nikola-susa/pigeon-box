package scheduler

import (
	"context"
	"github.com/nikola-susa/secret-chat/store"
)

type Worker struct {
	Store *store.Store
	Ctx   context.Context
}

func New(store *store.Store, ctx context.Context) *Worker {
	return &Worker{
		Store: store,
		Ctx:   ctx,
	}
}

func (w *Worker) Run() {
	go w.clearSessions()
}
