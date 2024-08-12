package slackaction

import (
	"fmt"
	"github.com/nikola-susa/pigeon-box/config"
	"github.com/nikola-susa/pigeon-box/serverevent"
	"github.com/nikola-susa/pigeon-box/storage"
	"github.com/nikola-susa/pigeon-box/store"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
)

type SlackAction struct {
	Api     *slack.Client
	Socket  *socketmode.Client
	Store   *store.Store
	Storage *storage.Storage
	Config  *config.Config
	Event   *serverevent.Server
}

func New(api *slack.Client, socket *socketmode.Client, store *store.Store, storage *storage.Storage, config *config.Config, events *serverevent.Server) *SlackAction {
	return &SlackAction{
		Api:     api,
		Socket:  socket,
		Store:   store,
		Storage: storage,
		Config:  config,
		Event:   events,
	}
}

func (s *SlackAction) Run(selfUserId string) {
	go func() {
		for envelope := range s.Socket.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:

				s.Socket.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)

				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MemberJoinedChannelEvent:
						if event.User == selfUserId {
							s.CreateChannelMessage(event)
						}
					default:
						s.Socket.Debugf("Skipped: %v", event)
					}
				default:
					s.Socket.Debugf("unsupported Events API eventPayload received")
				}

			case socketmode.EventTypeInteractive:

				callback, ok := envelope.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", envelope)
					continue
				}

				s.Socket.Debugf("Interaction received: %+v", callback)

				if callback.Type == slack.InteractionTypeDialogSubmission {
					if callback.CallbackID == "create-thread-dialog" {
						s.CreateThread(callback)

						s.Socket.Ack(*envelope.Request)
					}
				} else if callback.Type == slack.InteractionTypeBlockActions {

					for _, action := range callback.ActionCallback.BlockActions {
						if action.ActionID == "auth-thread" {
							s.AuthThread(callback)
						} else if action.ActionID == "create-thread" {
							s.CreateThreadDialogAction(envelope)
						}
					}

					s.Socket.Ack(*envelope.Request)
				}

			case socketmode.EventTypeSlashCommand:

				cmd, ok := envelope.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", envelope)
					continue
				}

				if cmd.Command == "/pigeon" {
					s.CreateThreadDialogCommand(envelope)
				}
			default:
				s.Socket.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	go func() {
		err := s.Socket.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
