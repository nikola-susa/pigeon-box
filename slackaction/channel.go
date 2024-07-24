package slackaction

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"log"
)

func (s *SlackAction) CreateChannelMessage(event *slackevents.MemberJoinedChannelEvent) {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi <@%v> :wave:", event.Inviter), false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "This is a welcome message with list of actions", false, false),
			nil, nil,
		),
	}

	msgOptions := slack.MsgOptionBlocks(blocks...)

	_, _, err := s.Api.PostMessage(event.Channel, msgOptions)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}
}
