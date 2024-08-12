package main

import (
	"context"
	"fmt"
	"github.com/nikola-susa/pigeon-box/app"
	"github.com/nikola-susa/pigeon-box/config"
	"github.com/nikola-susa/pigeon-box/scheduler"
	"github.com/nikola-susa/pigeon-box/serverevent"
	"github.com/nikola-susa/pigeon-box/slackaction"
	"github.com/nikola-susa/pigeon-box/storage"
	"github.com/nikola-susa/pigeon-box/store"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
)

func main() {
	configInit, err := config.NewConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	storeInit, err := store.New(configInit)
	if err != nil {
		log.Fatal("store error: ", err)
	}

	err = storeInit.Migrate()
	if err != nil {
		log.Fatal("store migrate error: ", err)
		return
	}

	storageInit := storage.New(configInit)

	sseInit := serverevent.New()

	slackApi := slack.New(configInit.Slack.BotToken, slack.OptionAppLevelToken(configInit.Slack.AppToken))
	socketMode := socketmode.New(slackApi)

	authTest, authTestErr := slackApi.AuthTest()
	if authTestErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	slackAction := slackaction.New(slackApi, socketMode, storeInit, &storageInit, configInit, sseInit)
	slackAction.Run(authTest.UserID)

	cron := scheduler.New(storeInit, sseInit, storageInit, configInit, context.Background())
	cron.Run()

	a := app.New(configInit, storageInit, storeInit, sseInit, slackApi, socketMode, slackAction)
	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}

}
