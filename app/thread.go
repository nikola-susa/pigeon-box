package app

import (
	"fmt"
	"github.com/nikola-susa/pigeon-box/crypt"
	"github.com/nikola-susa/pigeon-box/md"
	"github.com/nikola-susa/pigeon-box/model"
	"github.com/nikola-susa/pigeon-box/storage"
	"github.com/nikola-susa/pigeon-box/templates"
	"log"
	"net/http"
	"strconv"
	"time"
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

		HTMXRedirect(w, r, fmt.Sprintf("/error/%d", http.StatusInternalServerError))
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", threadId)
		HTMXRedirect(w, r, fmt.Sprintf("/error/%d", http.StatusNotFound))
		return
	}

	user, err := a.Store.GetUser(userId)
	if err != nil {
		log.Printf("Error getting user by slack id: %s", err)
		HTMXRedirect(w, r, fmt.Sprintf("/error/%d", http.StatusInternalServerError))
		return
	}

	if user == nil {
		log.Printf("Error user not found: %d", userId)
		HTMXRedirect(w, r, fmt.Sprintf("/error/%d", http.StatusNotFound))
		return
	}

	hashedUserId, err := crypt.HashIDEncodeInt(userId, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	renderUser := model.RenderUser{
		ID:       hashedUserId,
		Name:     *user.Name,
		Username: *user.Username,
		Avatar:   *user.Avatar,
		SlackID:  &user.SlackID,
	}

	expiresAt := ""
	if thread.ExpiresAt != nil {
		expr := *thread.ExpiresAt / 3600000000000
		if expr == 1 {
			expiresAt = "1 hour"
		} else {
			expiresAt = fmt.Sprintf("%d hours", expr)
		}
	}

	msgExpiresAt := ""
	if thread.MessagesExpireAt != nil {
		expr := *thread.MessagesExpireAt / 3600000000000
		if expr == 1 {
			msgExpiresAt = "1 hour"
		} else {
			msgExpiresAt = fmt.Sprintf("%d hours", expr)
		}
	}

	hashedAuthorId, err := crypt.HashIDEncodeInt(thread.UserID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	renderThread := model.RenderThread{
		ID:               strconv.Itoa(*thread.ID),
		Name:             thread.Name,
		Description:      thread.Description,
		AuthorID:         hashedAuthorId,
		IsAuthor:         userId == thread.UserID,
		ExpiresAt:        expiresAt,
		MessagesExpireAt: msgExpiresAt,
		Version:          a.Config.Version,
	}

	component := templates.Home(renderThread, r.PathValue("thread_id"), renderUser)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleGetMessages(w http.ResponseWriter, r *http.Request) {

	var lastId *int
	lastIdParam := r.URL.Query().Get("last_id")

	if lastIdParam != "" {
		lid, err := crypt.HashIDDecodeInt(lastIdParam, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
		if err != nil {
			log.Printf("Error decoding last id: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lastId = &lid
	}

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

	rawMessages, err := a.Store.GetMessagesByThread(threadId, lastId)
	if err != nil {
		RenderError(w)
		log.Printf("Error getting messages by thread: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rawMessages == nil {
		return
	}

	var lastMessageId *int
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

		stringMessage := ""

		if m.Text != nil && m.FileID == nil {
			key, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
			decryptedMessage, err := crypt.Decrypt(string(key), *m.Text)

			if err != nil {
				log.Printf("Error decrypting message: %s", err)
				continue
			}

			stringMessage = string(md.Parse(decryptedMessage))

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
				ID:            hashedFileId,
				Name:          file.Name,
				Size:          storage.StringSize(*file.Size),
				ContentType:   *file.ContentType,
				ThreadHash:    r.PathValue("thread_id"),
				ShouldPreview: storage.IsPreview(*file),
			}
		}

		hashedMessageId, err := crypt.HashIDEncodeInt(*m.ID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

		updatedAtFormatted := ""
		if m.UpdatedAt != m.CreatedAt {
			ut, err := ConvertTimeToUserRegion(r, m.UpdatedAt)
			if err != nil {
				log.Printf("Error converting time to user region: %s", err)
				continue
			}
			updatedAtFormatted = ut.Format("15:04:05")
			m.UpdatedAt = ut.Format(time.RFC3339)
		}

		createdAtFormatted := ""
		ct, err := ConvertTimeToUserRegion(r, m.CreatedAt)
		if err != nil {
			log.Printf("Error converting time to user region: %s", err)
			continue
		}
		createdAtFormatted = ct.Format("15:04:05")
		m.CreatedAt = ct.Format(time.RFC3339)

		hashedUserId, err := crypt.HashIDEncodeInt(*user.ID, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
		if err != nil {
			log.Printf("Error hashing user id: %s", err)
			continue
		}

		messages = append(messages, model.RenderMessage{
			ID:                 hashedMessageId,
			ThreadID:           r.PathValue("thread_id"),
			Text:               stringMessage,
			CreatedAt:          m.CreatedAt,
			CreatedAtFormatted: createdAtFormatted,
			UpdatedAt:          m.UpdatedAt,
			UpdatedAtFormatted: updatedAtFormatted,
			User: model.RenderUser{
				ID:       hashedUserId,
				Name:     *user.Name,
				Username: *user.Username,
				Avatar:   *user.Avatar,
			},
			Time:     m.CreatedAt,
			File:     renderFile,
			IsAuthor: userId == m.UserID,
		})

		lastMessageId = m.ID
	}

	lastMessageIdHashed, err := crypt.HashIDEncodeInt(*lastMessageId, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	component := templates.ChatList(messages, lastMessageIdHashed, r.PathValue("thread_id"))
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (a *App) HandleThreadDelete(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	files, err := a.Store.GetFilesByThread(threadId)
	if err != nil {
		log.Printf("Error getting files by thread: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if file.ID != nil {
			err := a.Storage.Delete(*file.Path)
			if err != nil {
				continue
			}
		}
	}

	err = a.Store.DeleteThread(threadId)
	if err != nil {
		log.Printf("Error deleting thread: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _, err = a.SlackApi.DeleteMessage(thread.SlackID, *thread.SlackTimestamp)
	if err != nil {
		log.Printf("Error deleting slack message: %s", err)
	}

	HTMXRedirect(w, r, "/not-authenticated")
}

func (a *App) HandleThreadSlackDetails(w http.ResponseWriter, r *http.Request) {

	threadId, err := crypt.HashIDDecodeInt(r.PathValue("thread_id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error decoding thread id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := a.Store.GetThread(threadId)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		return
	}

	if thread == nil {
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	var channelName string
	var channelDM string
	if thread.SlackID[0] == 'C' {
		channel, err := a.SlackAction.GetSlackChannel(thread.SlackID)
		if err != nil {
			log.Printf("Error getting slack channel by id: %s", err)
			return
		}
		channelName = channel.Name
	}

	if thread.SlackID[0] == 'D' {
		members, err := a.SlackAction.GetMPIMembers(thread.SlackID)
		if err != nil {
			log.Printf("Error getting slack channel members: %s", err)
			return
		}
		for i, member := range members {
			user, err := a.SlackAction.GetSlackUser(member)
			if err != nil {
				log.Printf("Error getting user by slack id: %s", err)
				return
			}
			if user == nil {
				log.Printf("Error user not found: %s", member)
				return
			}
			channelDM += user.RealName
			if i < len(members)-1 {
				channelDM += ", "
			}
		}
	}

	component := templates.SlackDetails(channelName, channelDM)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}

}

func (a *App) HandleThreadSlackWorkspace(w http.ResponseWriter, r *http.Request) {

	var workspaceName string
	workspace, err := a.SlackApi.GetTeamInfo()
	if err != nil {
		log.Printf("Error getting slack workspace by id: %s", err)
		return
	}
	workspaceName = workspace.Name

	component := templates.SlackWorkspace(workspaceName)
	err = component.Render(r.Context(), w)
	if err != nil {
		return
	}

}
