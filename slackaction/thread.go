package slackaction

import (
	"fmt"
	"github.com/nikola-susa/secret-chat/crypt"
	"github.com/nikola-susa/secret-chat/model"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"strconv"
	"time"
)

func (s *SlackAction) CreateThreadDialog(envelope socketmode.Event) {

	s.Socket.Ack(*envelope.Request)

	cmd, ok := envelope.Data.(slack.SlashCommand)
	if !ok {
		return
	}

	defaultName := ""

	if cmd.Text != "" {
		defaultName = cmd.Text
	}

	name := slack.NewTextInput("name", "Thread name", defaultName)
	name.Hint = "Name is not encrypted. It will help slack search the thread later."
	description := slack.NewTextAreaInput("description", "Description", "A new secure thread.")
	description.Optional = true
	description.Hint = "Description is not encrypted. It will help slack search the thread later."

	dialog := slack.Dialog{
		CallbackID:  "create-thread-dialog",
		Title:       "New Thread",
		Elements:    []slack.DialogElement{name, description},
		SubmitLabel: "Submit",
	}

	err := s.Api.OpenDialog(cmd.TriggerID, dialog)
	if err != nil {
		s.AckError(err, "Error creating thread", cmd.UserID)
		return
	}
}

func (s *SlackAction) CreateThread(payload slack.InteractionCallback) {

	userId, err := s.GetOrCreateUser(payload.User.ID)
	if err != nil {
		fmt.Println("create user", err)
		return
	}

	keyStr, _ := crypt.GenerateKey(32)

	key, err := crypt.Encrypt(s.Config.Crypt.Passphrase, []byte(keyStr))
	if err != nil {
		fmt.Println("encrypt", err)
		return
	}

	name := payload.Submission["name"]
	description := "A new secure thread."

	if payload.Submission["description"] != "" {
		description = payload.Submission["description"]
	}

	thread := model.Thread{
		Name:        name,
		Description: &description,
		UserID:      *userId,
		SlackID:     payload.Channel.ID,
		Key:         string(key),
	}

	id, err := s.Store.CreateThread(thread)
	if err != nil {
		fmt.Println("create thread", err)
		return
	}

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s* \n _%s_", name, description), false, false),
			nil, nil,
		),
		slack.NewActionBlock(
			"",
			slack.NewButtonBlockElement("auth-thread", strconv.Itoa(*id), slack.NewTextBlockObject("plain_text", "Access the Thread", true, false)).WithStyle("primary"),
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<!here>, new thread has been created by <@%v>", payload.User.ID), false, false),
		),
	}

	msgOptions := slack.MsgOptionBlocks(blocks...)

	_, messageTimestamp, err := s.Api.PostMessage(payload.Channel.ID, msgOptions, slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Println("post message", err)
		return
	}

	err = s.Store.SetThreadSlackTimestamp(*id, messageTimestamp)
	if err != nil {
		fmt.Println("set thread slack timestamp", err)
		return
	}

	return
}

func (s *SlackAction) AuthThread(payload slack.InteractionCallback) {

	if payload.Type == slack.InteractionTypeBlockActions {

		threadId := payload.ActionCallback.BlockActions[0].Value
		threadIdInt, err := strconv.Atoi(threadId)
		if err != nil {
			fmt.Println("parse thread ID", err)
			return
		}

		thread, err := s.Store.GetThread(threadIdInt)
		if err != nil {
			fmt.Println("get thread", err)
			return
		}

		if thread == nil {
			fmt.Println("thread not found")
			return
		}

		hashedThreadID, err := crypt.HashIDEncodeInt(threadIdInt, s.Config.Crypt.HashSalt, s.Config.Crypt.HashLength)
		if err != nil {
			log.Printf("Error encoding channel ID: %v", err)
			return
		}

		user, err := s.GetOrCreateUser(payload.User.ID)
		if err != nil {
			fmt.Println("create user", err)
			return
		}

		session, err := s.Store.CreateSession(model.Session{
			UserID:    *user,
			ThreadID:  threadIdInt,
			ExpiresAt: time.Now().Add(90 * time.Second),
		})

		hashedSessionID, err := crypt.HashIDEncodeInt(*session, s.Config.Crypt.HashSalt, s.Config.Crypt.HashLength)

		url := fmt.Sprintf("%s/auth/%s/%s", s.Config.Public.URL, hashedThreadID, hashedSessionID)

		blocks := authMessage(url, thread.Name, payload.Channel.ID, 60)

		msgOptions := slack.MsgOptionBlocks(blocks...)

		_, msgId, err := s.Api.PostMessage(payload.User.ID, msgOptions)
		if err != nil {
			s.AckError(err, fmt.Sprintf("Error posting ephemeral message: ```%v```", err), payload.User.ID)
			log.Printf("Failed to post message: %v", err)
		}

		go func() {
			time.Sleep(11 * time.Second)
			_, _, err := s.Api.DeleteMessage(payload.User.ID, msgId)
			if err != nil {
				fmt.Println("msdId", msgId)
				fmt.Println("payload.User.ID", payload.User.ID)
				log.Printf("Failed to delete message: %v", err)
			}
		}()

	}
}

func authMessage(url string, threadName string, channelId string, timeout int) []slack.Block {

	timeoutStr := strconv.Itoa(timeout)

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("You have been granted temporary access to *%s* thread in <#%s>.", threadName, channelId), false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Use this link to <%s|authenticate>", url), false, false),
			nil, nil,
		),
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("This message will self destruct in _%s sec_.", timeoutStr), false, false),
		),
	}
	return blocks
}
