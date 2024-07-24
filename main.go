package main

import (
	"context"
	"fmt"
	"github.com/nikola-susa/secret-chat/app"
	"github.com/nikola-susa/secret-chat/config"
	"github.com/nikola-susa/secret-chat/scheduler"
	"github.com/nikola-susa/secret-chat/serverevent"
	"github.com/nikola-susa/secret-chat/slackaction"
	"github.com/nikola-susa/secret-chat/storage"
	"github.com/nikola-susa/secret-chat/store"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
)

func main() {
	c, err := config.NewConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	storeInit, err := store.New(c)
	if err != nil {
		log.Fatal("store error: ", err)
	}

	storageInit := storage.New(c)

	events := serverevent.New()

	slackApi := slack.New(
		c.Slack.BotToken,
		slack.OptionAppLevelToken(c.Slack.AppToken),
		//slack.OptionDebug(true),
		//slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode := socketmode.New(
		slackApi,
		//socketmode.OptionDebug(true),
		//socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := slackApi.AuthTest()
	if authTestErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	selfUserId := authTest.UserID

	slackAction := slackaction.New(slackApi, socketMode, storeInit, &storageInit, c, events)

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:

				socketMode.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)

				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if event.User != selfUserId &&
							strings.Contains(strings.ToLower(event.Text), "hello") {
							_, _, err := slackApi.PostMessage(
								event.Channel,
								slack.MsgOptionText(
									fmt.Sprintf(":wave: Hi there, <@%v>!", event.User),
									false,
								),
							)
							if err != nil {
								log.Printf("Failed to reply: %v", err)
							}
						}
					case *slackevents.MemberJoinedChannelEvent:
						if event.User == selfUserId {
							slackAction.CreateChannelMessage(event)
						}
					default:
						socketMode.Debugf("Skipped: %v", event)
					}
				default:
					socketMode.Debugf("unsupported Events API eventPayload received")
				}

			case socketmode.EventTypeInteractive:

				callback, ok := envelope.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", envelope)
					continue
				}

				socketMode.Debugf("Interaction received: %+v", callback)

				if callback.Type == slack.InteractionTypeDialogSubmission {
					if callback.CallbackID == "create-thread-dialog" {
						slackAction.CreateThread(callback)

						socketMode.Ack(*envelope.Request)
					}
				} else if callback.Type == slack.InteractionTypeBlockActions {

					fmt.Println("callback.ActionCallback.BlockActions", callback.ActionCallback.BlockActions)

					for _, action := range callback.ActionCallback.BlockActions {
						if action.ActionID == "auth-thread" {
							slackAction.AuthThread(callback)
						}
					}

					socketMode.Ack(*envelope.Request)
				}

			case socketmode.EventTypeSlashCommand:

				cmd, ok := envelope.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", envelope)

					continue
				}

				socketMode.Debugf("Slash command received: %+v", cmd)
				socketMode.Debugf("Slash command text: %+v", cmd.Text)

				if cmd.Command == "/echo" {
					slackAction.CreateThreadDialog(envelope)
				}
			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	go func() {
		err = socketMode.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	scheduler.New(storeInit, context.Background())

	a := app.New(c, storageInit, storeInit, events, slackApi, socketMode, slackAction)
	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}

}
