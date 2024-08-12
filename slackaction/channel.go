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
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%v> added Pigeon box :bird: to this channel.", event.Inviter), false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Pigeon box allows you to communicate securely within your organization.", false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Messages and files are encrypted and not accessible by slack or anyone outside of this group.", false, false),
			nil, nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Thread and messages can have expiration dates and you're encouraged to set it.", false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "To create a new thread run `/pigeon` command, or simply", false, false),
			nil, nil,
		),
		slack.NewActionBlock(
			"",
			slack.NewButtonBlockElement("create-thread", "", slack.NewTextBlockObject("plain_text", "Create a Thread", true, false)),
		),
	}

	msgOptions := slack.MsgOptionBlocks(blocks...)

	_, _, err := s.Api.PostMessage(event.Channel, msgOptions)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}
}
