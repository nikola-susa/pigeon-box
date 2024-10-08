package app

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/nikola-susa/pigeon-box/crypt"
	"github.com/nikola-susa/pigeon-box/htmx"
	"github.com/nikola-susa/pigeon-box/md"
	"github.com/nikola-susa/pigeon-box/model"
	"github.com/nikola-susa/pigeon-box/templates"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *App) HandleCreateNewMessage(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(contextKey("user_id")).(int)

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		htmx.ErrorToast(w, fmt.Sprintf("Error processing a from"))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")

	if message == "" {
		htmx.ErrorToast(w, "Cannot create message w/o content")
		http.Error(w, "No message or files", http.StatusBadRequest)
		return
	}

	user, err := a.Store.GetUser(userId)
	if err != nil {
		log.Printf("Error getting user by slack id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		HTMXRedirect(w, r, "/not-authenticated")
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	threadId, err := crypt.HashIDDecodeInt(r.PathValue("thread_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding thread id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := a.Store.GetThread(threadId)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		HTMXRedirect(w, r, "/not-authenticated")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", threadId)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	threadKey, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
	if err != nil {
		log.Printf("Error decrypting thread key: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedMessage, err := crypt.Encrypt(string(threadKey), []byte(message))
	if err != nil {
		log.Printf("Error encrypting message: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := model.CreateMessageParams{
		UserID:   *user.ID,
		ThreadID: *thread.ID,
		Text:     encryptedMessage,
	}

	id, err := a.Store.CreateMessage(m)
	if err != nil {
		log.Printf("Error creating message: %s", err)
		htmx.ErrorToast(w, "Error creating message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread.MessagesExpireAt != nil {
		expiration := time.Now().Add(*thread.MessagesExpireAt)
		err = a.Store.SetMessageExpiresAt(*id, expiration)
	}

	var renderFile model.RenderFile

	hashedMessageId, err := crypt.HashIDEncodeInt(*id, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	createdAtFormatted := ""
	ct, err := ConvertTimeToUserRegion(r, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Printf("Error converting time to user region: %s", err)
		return
	}
	createdAtFormatted = ct.Format("15:04:05")

	messageRender := model.RenderMessage{
		ID:                 hashedMessageId,
		ThreadID:           r.PathValue("thread_id"),
		Text:               string(md.Parse([]byte(message))),
		CreatedAt:          ct.Format(time.RFC3339),
		CreatedAtFormatted: createdAtFormatted,
		User: model.RenderUser{
			ID:       strconv.Itoa(*user.ID),
			Name:     *user.Name,
			Username: *user.Username,
			Avatar:   *user.Avatar,
		},
		Time:     time.Now().Format(time.RFC3339),
		File:     renderFile,
		IsAuthor: userId == m.UserID,
	}

	component := templates.ChatBubble(messageRender)
	htmlString, err := templ.ToGoHTML(context.Background(), component)
	if err != nil {
		fmt.Println("Error rendering component:", err)
		return
	}

	eventName := "created:" + r.PathValue("thread_id")
	a.Event.Broadcast(r.PathValue("thread_id"), []byte(htmlString), &eventName, nil, nil)

	c := templates.CreateMessageForm(r.PathValue("thread_id"))
	err = c.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleRenderEdit(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(contextKey("user_id")).(int)

	user, err := a.Store.GetUser(userId)
	if err != nil {
		log.Printf("Error getting user by slack id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	messageId, err := crypt.HashIDDecodeInt(r.PathValue("message_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding message id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := a.Store.GetMessage(messageId)

	if err != nil {
		log.Printf("Error getting message by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m == nil {
		log.Printf("Error message not found: %d", messageId)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	if m.UserID != userId {
		log.Printf("Error user not allowed to edit message: %d, author %d, msg %d", userId, m.UserID, m.ID)
		htmx.ErrorToast(w, "User not allowed to edit message")
		http.Error(w, "user not allowed to edit message", http.StatusForbidden)
		return
	}

	thread, err := a.Store.GetThread(m.ThreadID)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", m.ThreadID)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	stringMessage := ""

	if m.Text != nil {
		key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
		decryptedMessage, err := crypt.Decrypt(string(key), *m.Text)

		if err != nil {
			log.Printf("Error decrypting message: %s", err)
			return
		}

		stringMessage = string(decryptedMessage)
	}

	message := model.RenderMessage{
		ID:        r.PathValue("message_id"),
		ThreadID:  r.PathValue("thread_id"),
		Text:      stringMessage,
		CreatedAt: m.CreatedAt,
		User: model.RenderUser{
			ID:       strconv.Itoa(*user.ID),
			Name:     *user.Name,
			Username: *user.Username,
			Avatar:   *user.Avatar,
		},
		Time:     m.CreatedAt,
		IsAuthor: userId == m.UserID,
	}

	htmx.MessageEditedEvent(r.PathValue("message_id")).Output(w)

	component := templates.EditMessageForm(message)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleChatBubbleRender(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(contextKey("user_id")).(int)

	user, err := a.Store.GetUser(userId)
	if err != nil {
		log.Printf("Error getting user by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		http.Error(w, "user not found", http.StatusNotFound)
	}

	messageId, err := crypt.HashIDDecodeInt(r.PathValue("message_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding message id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := a.Store.GetMessage(messageId)

	if err != nil {
		log.Printf("Error getting message by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m == nil {
		log.Printf("Error message not found: %d", messageId)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	if m.UserID != userId {
		log.Printf("Error user not allowed to edit message: %d", userId)
		http.Error(w, "user not allowed to edit message", http.StatusForbidden)
		return
	}

	thread, err := a.Store.GetThread(m.ThreadID)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", m.ThreadID)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	updatedAtFormatted := ""
	if m.UpdatedAt != m.CreatedAt {
		ut, err := ConvertTimeToUserRegion(r, m.UpdatedAt)
		if err != nil {
			log.Printf("Error converting time to user region: %s", err)
			return
		}
		updatedAtFormatted = ut.Format("15:04:05")
		m.UpdatedAt = ut.Format(time.RFC3339)
	}

	createdAtFormatted := ""
	ct, err := ConvertTimeToUserRegion(r, m.CreatedAt)
	if err != nil {
		log.Printf("Error converting time to user region: %s", err)
		return
	}
	createdAtFormatted = ct.Format("15:04:05")
	m.CreatedAt = ct.Format(time.RFC3339)

	stringMessage := ""
	if m.Text != nil {
		key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
		decryptedMessage, err := crypt.Decrypt(string(key), *m.Text)

		if err != nil {
			log.Printf("Error decrypting message: %s", err)
			return
		}

		stringMessage = string(md.Parse(decryptedMessage))
	}

	message := model.RenderMessage{
		ID:                 r.PathValue("message_id"),
		ThreadID:           r.PathValue("thread_id"),
		Text:               stringMessage,
		CreatedAt:          m.CreatedAt,
		CreatedAtFormatted: createdAtFormatted,
		UpdatedAt:          m.UpdatedAt,
		UpdatedAtFormatted: updatedAtFormatted,
		User: model.RenderUser{
			ID:       strconv.Itoa(*user.ID),
			Name:     *user.Name,
			Username: *user.Username,
			Avatar:   *user.Avatar,
		},
		Time:     m.CreatedAt,
		IsAuthor: userId == m.UserID,
	}

	component := templates.ChatBubble(message)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleEditSubmit(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(contextKey("user_id")).(int)

	user, err := a.Store.GetUser(userId)
	if err != nil {
		log.Printf("Error getting user by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		http.Error(w, "user not found", http.StatusNotFound)
	}

	messageText := r.FormValue("message")

	if messageText == "" {
		htmx.ErrorToast(w, "Cannot create message w/o content or files")
		http.Error(w, "No message or files", http.StatusBadRequest)
		return
	}

	messageId, err := crypt.HashIDDecodeInt(r.PathValue("message_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding message id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := a.Store.GetMessage(messageId)

	if err != nil {
		log.Printf("Error getting message by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m == nil {
		log.Printf("Error message not found: %d", messageId)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	if m.UserID != userId {
		log.Printf("Error user not allowed to edit message: %d", userId)
		http.Error(w, "user not allowed to edit message", http.StatusForbidden)
		return
	}

	thread, err := a.Store.GetThread(m.ThreadID)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", m.ThreadID)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	threadKey, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
	if err != nil {
		log.Printf("Error decrypting thread key: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedMessage, err := crypt.Encrypt(string(threadKey), []byte(messageText))
	if err != nil {
		log.Printf("Error encrypting message: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.Store.UpdateMessage(model.UpdateMessageParams{
		ID:   messageId,
		Text: encryptedMessage,
	})
	if err != nil {
		log.Printf("updating message failed: %d", messageId)
		htmx.ErrorToast(w, "Error updating message")
		http.Error(w, "updating message failed", http.StatusInternalServerError)
	}

	eventName := "edited:" + r.PathValue("thread_id")
	a.Event.Broadcast(r.PathValue("thread_id"), []byte(r.PathValue("message_id")), &eventName, nil, nil)

	//htmx.MessageUpdatedEvent(r.PathValue("message_id")).Output(w)

	component := templates.CreateMessageForm(r.PathValue("thread_id"))
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleEditCancel(w http.ResponseWriter, r *http.Request) {
	threadId := r.PathValue("thread_id")

	htmx.MessageEditCancelledEvent(r.PathValue("message_id")).Output(w)

	component := templates.CreateMessageForm(threadId)
	err := component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(contextKey("user_id")).(int)

	messageId, err := crypt.HashIDDecodeInt(r.PathValue("message_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding message id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := a.Store.GetMessage(messageId)
	if err != nil {
		log.Printf("Error getting message by id: %s", err)
		htmx.ErrorToast(w, "Error getting message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	if message == nil {
		log.Printf("Error message not found: %d", messageId)
		htmx.ErrorToast(w, "Error getting message")
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	if message.UserID != userId {
		log.Printf("Error user not allowed to delete message: %d, author: %d, msd: %s", userId, message.UserID, message.ID)
		htmx.ErrorToast(w, "User not allowed to delete message")
		http.Error(w, "user not allowed to delete message", http.StatusForbidden)
		return
	}

	err = a.Store.DeleteMessage(messageId)
	if err != nil {
		log.Printf("Error deleting message: %s", err)
		htmx.ErrorToast(w, "Error deleting the message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if message.FileID != nil {
		file, _ := a.Store.GetFile(*message.FileID)
		if file.ID != nil {
			err := a.Storage.Delete(*file.Path)
			if err != nil {
				return
			}
			err = a.Store.DeleteFile(*message.FileID)
			if err != nil {
				log.Printf("Error deleting file: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	eventName := "deleted:" + r.PathValue("thread_id")
	a.Event.Broadcast(r.PathValue("thread_id"), []byte(r.PathValue("message_id")), &eventName, nil, nil)

	w.WriteHeader(http.StatusOK)
}
