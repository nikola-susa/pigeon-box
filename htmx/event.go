package htmx

import (
	"encoding/json"
	"net/http"
)

const (
	MESSAGE_CREATED        = "MessageCreated"
	MESSAGE_UPDATED        = "MessageUpdated"
	MESSAGE_EDITED         = "MessageEdited"
	MESSAGE_EDIT_CANCELLED = "MessageEditCancelled"
	MESSAGE_DELETED        = "MessageDeleted"

	FILE_CREATED = "FileCreated"
	FILE_UPDATED = "FileUpdated"
	FILE_DELETED = "FileDeleted"

	THREAD_DELETED = "ThreadDeleted"
)

type Event struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

func NewEvent(t string, target string) Event {
	return Event{Type: t, Target: target}
}

func MessageCreatedEvent(target string) Event {
	return NewEvent(MESSAGE_CREATED, target)
}

func MessageUpdatedEvent(target string) Event {
	return NewEvent(MESSAGE_UPDATED, target)
}

func MessageEditedEvent(target string) Event {
	return NewEvent(MESSAGE_EDITED, target)
}

func MessageEditCancelledEvent(target string) Event {
	return NewEvent(MESSAGE_EDIT_CANCELLED, target)
}

func MessageDeletedEvent(target string) Event {
	return NewEvent(MESSAGE_DELETED, target)
}

func FileCreatedEvent(target string) Event {
	return NewEvent(FILE_CREATED, target)
}

func FileDeletedEvent(target string) Event {
	return NewEvent(FILE_DELETED, target)
}

func FileUpdatedEvent(target string) Event {
	return NewEvent(FILE_UPDATED, target)
}

func ThreadDeletedEvent(target string) Event {
	return NewEvent(THREAD_DELETED, target)
}

func (e Event) Output(w http.ResponseWriter) {
	key := e.Type + e.Target

	eventMap := map[string]Event{}
	eventMap[key] = e
	eventMap[e.Type] = e
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", string(jsonData))
}
