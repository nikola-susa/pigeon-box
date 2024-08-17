package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/nikola-susa/pigeon-box/crypt"
	"github.com/nikola-susa/pigeon-box/htmx"
	"github.com/nikola-susa/pigeon-box/model"
	"github.com/nikola-susa/pigeon-box/storage"
	"github.com/nikola-susa/pigeon-box/templates"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *App) HandleCreateFileMessage(w http.ResponseWriter, r *http.Request) {

	threadId := r.Context().Value(contextKey("thread_id")).(int)

	userId := r.Context().Value(contextKey("user_id")).(int)

	maxSizeInt, err := strconv.Atoi(a.Config.File.MaxSize)
	if err != nil {
		log.Printf("Error parsing max size: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	maxFileSize := int64(maxSizeInt) << 20

	err = r.ParseMultipartForm(maxFileSize)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	for _, fh := range files {

		if fh.Size > maxFileSize {
			log.Printf("File size exceeds the maximum allowed size: %s", fh.Filename)
			htmx.ErrorToast(w, fmt.Sprintf("File %s size (%s) exceeds the maximum allowed %s", fh.Filename, storage.StringSize(fh.Size), storage.StringSize(maxFileSize)))
			http.Error(w, "File size exceeds the maximum allowed size", http.StatusBadRequest)
			return
		}

		f, err := fh.Open()
		if err != nil {
			log.Printf("Error opening file: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf := bytes.NewBuffer(nil)

		if _, err := buf.ReadFrom(f); err != nil {
			log.Printf("Error reading file: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		threadKey, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)

		byt := buf.Bytes()
		fn := fmt.Sprintf("%d-%s", time.Now().Unix(), fh.Filename)

		eByt, err := crypt.Encrypt(string(threadKey), byt)
		if err != nil {
			log.Printf("Error encrypting file: %s", err)
			htmx.ErrorToast(w, fmt.Sprintf("Error encrypting file"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fp, err := a.Storage.Upload(fn, eByt)
		if err != nil {
			log.Printf("Error uploading file: %s", err)
			htmx.ErrorToast(w, fmt.Sprintf("Error uploading file"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		contentType := fh.Header.Get("Content-Type")

		nf := model.File{
			Name:        fh.Filename,
			Path:        &fp,
			Size:        &fh.Size,
			ContentType: &contentType,
			ThreadID:    &threadId,
			UserID:      &userId,
		}

		id, err := a.Store.CreateFile(nf)
		if err != nil {
			log.Printf("Error creating file: %s", err)
			htmx.ErrorToast(w, fmt.Sprintf("Error creating file"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m := model.CreateMessageParams{
			UserID:   userId,
			ThreadID: threadId,
			Text:     nil,
		}

		messageId, err := a.Store.CreateMessage(m)
		if err != nil {
			log.Printf("Error creating message: %s", err)
			htmx.ErrorToast(w, fmt.Sprintf("Error creating message"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = a.Store.SetMessageFileID(*messageId, *id)
		if err != nil {
			log.Printf("Error setting message file id: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if thread.MessagesExpireAt != nil {
			expiration := time.Now().Add(*thread.MessagesExpireAt)
			err = a.Store.SetMessageExpiresAt(*messageId, expiration)
		}

		user, err := a.Store.GetUser(userId)
		if err != nil {
			log.Printf("Error getting user by slack id: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedFileId, err := crypt.HashIDEncodeInt(*id, a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)
		if err != nil {
			log.Printf("Error hashing file id: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
			Text:               "",
			CreatedAt:          ct.Format(time.RFC3339),
			CreatedAtFormatted: createdAtFormatted,
			User: model.RenderUser{
				ID:       strconv.Itoa(*user.ID),
				Name:     *user.Name,
				Username: *user.Username,
				Avatar:   *user.Avatar,
			},
			Time: time.Now().Format("2006-01-02 15:04:05"),
			File: model.RenderFile{
				ID:            hashedFileId,
				Name:          nf.Name,
				Size:          storage.StringSize(*nf.Size),
				ContentType:   contentType,
				ThreadHash:    r.PathValue("thread_id"),
				ShouldPreview: storage.IsPreview(nf),
			},
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

		_ = f.Close()

	}
}

func (a *App) HandleDownloadFile(w http.ResponseWriter, r *http.Request) {

	download := r.URL.Query().Get("download") == "true"
	fileId, err := crypt.HashIDDecodeInt(r.PathValue("id"), a.Config.Crypt.HashSalt, a.Config.Crypt.HashLength)

	file, err := a.Store.GetFile(fileId)
	if err != nil {
		log.Printf("Error getting file by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if file == nil {
		log.Printf("Error file not found: %d", fileId)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	thread, err := a.Store.GetThread(*file.ThreadID)
	if err != nil {
		log.Printf("Error getting thread by id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if thread == nil {
		log.Printf("Error thread not found: %d", file.ThreadID)
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	threadKey, err := crypt.Decrypt(a.Config.Crypt.Passphrase, thread.Key)
	threadKeyStr := string(threadKey)

	dByt, err := a.Storage.Get(*file.Path, &threadKeyStr)
	if err != nil {
		log.Printf("Error downloading file: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if download {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.Name))
	}

	w.Header().Set("Content-Type", *file.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(dByt)))

	_, err = w.Write(dByt)
	if err != nil {
		log.Printf("Error writing file: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("Error parsing file id: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := a.Store.GetFile(id)
	if err != nil {
		log.Printf("Error getting file by id: %s", err)
		htmx.ErrorToast(w, "Error getting file by id")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if file == nil {
		log.Printf("Error file not found: %d", id)
		htmx.ErrorToast(w, "file not found")
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	err = a.Storage.Delete(*file.Path)
	if err != nil {
		log.Printf("Error deleting file: %s", err)
		htmx.ErrorToast(w, "Error deleting file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.Store.DeleteFile(id)
	if err != nil {
		log.Printf("Error deleting file from db: %s", err)
		htmx.ErrorToast(w, "Error deleting file from db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmx.SuccessToast(w, fmt.Sprintf("File \"%s\" deleted", file.Name))

	Resource(w, http.StatusOK, id)

}
