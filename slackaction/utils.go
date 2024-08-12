package slackaction

import (
	"fmt"
	"github.com/slack-go/slack"
	"os"
)

func (s *SlackAction) AckError(err error, msg string, channel string) {
	_, _, err = s.Api.PostMessage(channel, slack.MsgOptionText(msg, false))
	if err != nil {
		return
	}
}

func (s *SlackAction) GetBotUserID() string {

	authTest, authTestErr := s.Api.AuthTest()
	if authTestErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		if err != nil {
			return ""
		}
		os.Exit(1)
	}

	return authTest.UserID
}

func (s *SlackAction) GetSlackUser(userId string) (*slack.User, error) {
	user, err := s.Api.GetUserInfo(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SlackAction) GetSlackUsersByChannel(channelId string) ([]slack.User, error) {
	members, _, err := s.Api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	var users []slack.User

	for _, member := range members {
		user, err := s.GetSlackUser(member)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (s *SlackAction) GetSlackChannelMembers(channelId string) ([]string, error) {
	members, _, err := s.Api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *SlackAction) GetSlackChannel(channelId string) (*slack.Channel, error) {
	channel, err := s.Api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *SlackAction) GetMPIMembers(channelId string) ([]string, error) {
	members, _, err := s.Api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *SlackAction) GetChannelName(channelId string) (string, error) {
	if channelId[0] == 'C' {
		channel, err := s.GetSlackChannel(channelId)
		if err != nil {
			return "", err
		}
		return channel.Name, nil
	} else if channelId[0] == 'D' {
		members, err := s.GetMPIMembers(channelId)
		if err != nil {
			return "", err
		}
		var userIds string
		for _, member := range members {
			userIds += member + ","
		}
		return userIds, nil

	}

	return "", nil
}
