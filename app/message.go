package app

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/nikola-susa/secret-chat/crypt"
	"github.com/nikola-susa/secret-chat/htmx"
	"github.com/nikola-susa/secret-chat/model"
	"github.com/nikola-susa/secret-chat/templates"
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
		htmx.ErrorToast(w, fmt.Sprintf("Error parsing form: %s", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")

	if message == "" {
		htmx.ErrorToast(w, "Cannot create message w/o content or files")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", threadId)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	threadKey, err := crypt.Decrypt(a.Config.Crypt.Passphrase, []byte(thread.Key))

	encryptedMessage, err := crypt.Encrypt(string(threadKey), []byte(message))
	if err != nil {
		log.Printf("Error encrypting message: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := model.CreateMessageParams{
		UserID:   *user.ID,
		ThreadID: *thread.ID,
		Text:     string(encryptedMessage),
	}

	id, err := a.Store.CreateMessage(m)
	if err != nil {
		log.Printf("Error creating message: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var renderFile model.RenderFile

	hashedMessageId, err := crypt.HashIDEncodeInt(*id, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	messageRender := model.RenderMessage{
		ID:        hashedMessageId,
		ThreadID:  r.PathValue("thread_id"),
		Text:      message,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		User: model.RenderUser{
			ID:       strconv.Itoa(*user.ID),
			Name:     *user.Name,
			Username: *user.Username,
			Avatar:   *user.Avatar,
		},
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		File:     renderFile,
		IsAuthor: userId == m.UserID,
	}

	component := templates.ChatBubble(messageRender)
	htmlString, err := templ.ToGoHTML(context.Background(), component)
	if err != nil {
		fmt.Println("Error rendering component:", err)
		return
	}

	a.Event.Broadcast(r.PathValue("thread_id"), []byte(htmlString))

	Resource(w, http.StatusCreated, id)
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

	stringMessage := *m.Text

	if *m.Text != "" && m.Text != nil {
		key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, []byte(thread.Key))
		decryptedMessage, err := crypt.Decrypt(string(key), []byte(*m.Text))

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

	component := templates.EditMessage(message)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleEditCancel(w http.ResponseWriter, r *http.Request) {
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

	stringMessage := *m.Text

	if *m.Text != "" && m.Text != nil {
		key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, []byte(thread.Key))
		decryptedMessage, err := crypt.Decrypt(string(key), []byte(*m.Text))

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

	component := templates.ChatBubble(message)
	err = component.Render(r.Context(), w)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	if message == nil {
		log.Printf("Error message not found: %d", messageId)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	if message.UserID != userId {
		log.Printf("Error user not allowed to delete message: %d", userId)
		http.Error(w, "user not allowed to delete message", http.StatusForbidden)
		return
	}

	err = a.Store.DeleteMessage(messageId)
	if err != nil {
		log.Printf("Error deleting message: %s", err)
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

	w.WriteHeader(http.StatusOK)
}
