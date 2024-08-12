package slackaction

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

func (s *SlackAction) CreateMessage(context string, channelID string, userID string, url string) string {

	blocks := []slack.Block{
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("New secure message from <@%v> ", userID), false, false),
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", context, false, false),
			nil, nil,
		),
		slack.NewActionBlock(
			"",
			slack.NewButtonBlockElement("", "", slack.NewTextBlockObject("plain_text", "Read the message", false, false)).WithStyle("primary").WithURL(url),
		),
	}

	msgOptions := slack.MsgOptionBlocks(blocks...)

	msg, _, err := s.Api.PostMessage(channelID, msgOptions)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}

	return msg
}
