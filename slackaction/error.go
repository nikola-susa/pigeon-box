package slackaction

import "github.com/slack-go/slack"

func (s *SlackAction) AckError(err error, msg string, channel string) {
	_, _, err = s.Api.PostMessage(channel, slack.MsgOptionText(msg, false))
	if err != nil {
		return
	}

}
