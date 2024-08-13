package scheduler

import (
	"fmt"
	"github.com/nikola-susa/pigeon-box/crypt"
	"time"
)

func (w *Worker) clearSessions() {
	ticker := time.NewTicker(3 * time.Minute)
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

func (w *Worker) clearMessages() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.collectDeletedExpiredMessages()
		case <-w.Ctx.Done():
			return
		}
	}
}

func (w *Worker) collectDeletedExpiredMessages() {

	messages, err := w.Store.GetExpiredMessages()
	if err != nil {
		fmt.Printf("error getting expired messages: %s", err)
		return
	}
	if messages == nil {
		return
	}

	for _, message := range messages {

		threadIdHashed, err := crypt.HashIDEncodeInt(message.ThreadID, w.Config.Crypt.HashSalt, w.Config.Crypt.HashLength)
		messageIdHashed, err := crypt.HashIDEncodeInt(*message.ID, w.Config.Crypt.HashSalt, w.Config.Crypt.HashLength)

		eventName := "deleted:" + threadIdHashed
		w.Event.Broadcast(threadIdHashed, []byte(messageIdHashed), &eventName, nil, nil)

		err = w.Store.DeleteMessage(*message.ID)
		if err != nil {
			fmt.Println("error deleting message: ", err)
			return
		}

		if message.FileID != nil {
			file, _ := w.Store.GetFile(*message.FileID)
			if file.ID != nil {
				err := w.Storage.Delete(*file.Path)
				if err != nil {
					return
				}
				err = w.Store.DeleteFile(*message.FileID)
				if err != nil {
					return
				}
			}
		}
	}

}

func (w *Worker) clearThread() {
	ticker := time.NewTicker(45 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.collectExpiredThreads()
		case <-w.Ctx.Done():
			return
		}
	}
}

func (w *Worker) collectExpiredThreads() {
	threads, err := w.Store.GetExpiredThreads()
	if err != nil {
		fmt.Printf("error getting expired threads: %s", err)
		return
	}
	if threads == nil {
		return
	}

	for _, thread := range threads {
		threadIdHashed, err := crypt.HashIDEncodeInt(*thread.ID, w.Config.Crypt.HashSalt, w.Config.Crypt.HashLength)

		eventName := "deleted:" + threadIdHashed
		w.Event.Broadcast(threadIdHashed, nil, &eventName, nil, nil)

		err = w.Store.DeleteThread(*thread.ID)
		if err != nil {
			fmt.Println("error deleting thread: ", err)
			return
		}
	}
}
