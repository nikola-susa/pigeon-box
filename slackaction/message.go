package slackaction

import (
	"fmt"
	"github.com/nikola-susa/secret-chat/crypt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"time"
)

func (s *SlackAction) CreateMessageDialog(envelope socketmode.Event) {
	cmd, ok := envelope.Data.(slack.SlashCommand)
	if !ok {
		log.Printf("Ignored %+v\n", envelope)
		return
	}

	userID, err := crypt.HashIDEncodeString(cmd.UserID, s.Config.Crypt.HashSalt, s.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error encoding user ID: %v", err)
		s.AckError(err, fmt.Sprintf("Error encoding user ID: ```%v```", err), cmd.UserID)
		s.Socket.Ack(*envelope.Request)
		return
	}

	channelID, err := crypt.HashIDEncodeString(cmd.ChannelID, s.Config.Crypt.HashSalt, s.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error encoding channel ID: %v", err)
		s.AckError(err, fmt.Sprintf("Error encoding channel ID: ```%v```", err), cmd.UserID)
		s.Socket.Ack(*envelope.Request)
		return
	}

	timeOutUnix := time.Now().Add(time.Minute * 30).Unix()
	timeOut, err := crypt.HashIDEncodeInt(int(timeOutUnix), s.Config.Crypt.HashSalt, s.Config.Crypt.HashLength)
	if err != nil {
		log.Printf("Error encoding timeout: %v", err)
		s.AckError(err, fmt.Sprintf("Error encoding timeout: ```%v```", err), cmd.UserID)
		s.Socket.Ack(*envelope.Request)
		return
	}

	var link string

	link = fmt.Sprintf("%s/msg/new/%s/%s/%s", s.Config.Public.URL, userID, channelID, timeOut)

	blocks := []slack.Block{

		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Create a new secure message", false, false),
			nil, nil,
		),
		slack.NewActionBlock(
			"",
			slack.NewButtonBlockElement("", "", slack.NewTextBlockObject("plain_text", "Create a new message", false, false)).WithStyle("primary").WithURL(link),
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Message will be visible to all members of %s", cmd.ChannelName), false, false),
		),
	}

	msgOptions := slack.MsgOptionBlocks(blocks...)

	_, err = s.Api.PostEphemeral(cmd.ChannelID, cmd.UserID, msgOptions)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}

	s.Socket.Ack(*envelope.Request)
}

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
