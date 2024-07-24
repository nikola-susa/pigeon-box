package app

import (
	"github.com/nikola-susa/secret-chat/crypt"
	"github.com/nikola-susa/secret-chat/model"
	"github.com/nikola-susa/secret-chat/storage"
	"github.com/nikola-susa/secret-chat/templates"
	"log"
	"net/http"
	"strconv"
)

func (a *App) HandleRenderThread(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(contextKey("user_id")).(int)

	threadId, err := crypt.HashIDDecodeInt(r.PathValue("thread_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding thread id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := a.Store.GetThread(threadId)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)

		RenderError(w)
		return
	}

	if thread == nil {
		RenderError(w)
		log.Printf("Error thread not found: %d", threadId)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	user, err := a.Store.GetUser(userId)
	if err != nil {
		RenderError(w)
		log.Printf("Error getting user by slack id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	rawMessages, err := a.Store.GetMessagesByThread(threadId)
	if err != nil {
		RenderError(w)
		log.Printf("Error getting messages by thread: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []model.RenderMessage
	for _, m := range rawMessages {
		user, err := a.Store.GetUser(m.UserID)
		if err != nil {
			log.Printf("Error getting user by slack id: %s", err)
			continue
		}

		if user == nil {
			log.Printf("Error user not found: %s", m.UserID)
			continue
		}

		stringMessage := *m.Text

		if *m.Text != "" && m.Text != nil {
			key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, []byte(thread.Key))
			decryptedMessage, err := crypt.Decrypt(string(key), []byte(*m.Text))

			if err != nil {
				log.Printf("Error decrypting message: %s", err)
				continue
			}

			stringMessage = string(decryptedMessage)
		}

		var renderFile model.RenderFile

		if m.FileID != nil {
			file, err := a.Store.GetFile(*m.FileID)
			if err != nil {
				log.Printf("Error getting file by id: %s", err)
				continue
			}

			if file == nil {
				log.Printf("Error file not found: %d", *m.FileID)
				continue
			}

			hashedFileId, err := crypt.HashIDEncodeInt(*file.ID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
			if err != nil {
				log.Printf("Error hashing file id: %s", err)
				continue
			}

			renderFile = model.RenderFile{
				ID:          hashedFileId,
				Name:        file.Name,
				Size:        storage.StringSize(*file.Size),
				ContentType: *file.ContentType,
				ThreadHash:  r.PathValue("thread_id"),
			}
		}

		hashedMessageId, err := crypt.HashIDEncodeInt(*m.ID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

		messages = append(messages, model.RenderMessage{
			ID:        hashedMessageId,
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
			File:     renderFile,
			IsAuthor: userId == m.UserID,
		})

	}

	renderUser := model.RenderUser{
		ID:       strconv.Itoa(*user.ID),
		Name:     *user.Name,
		Username: *user.Username,
		Avatar:   *user.Avatar,
	}

	component := templates.Home(*thread, r.PathValue("thread_id"), messages, renderUser)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}
