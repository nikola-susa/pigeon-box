package scheduler

import "time"

func (w *Worker) clearSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = w.Store.DeleteExpiredSessions()
		case <-w.Ctx.Done():
			return
		}
	}
}
