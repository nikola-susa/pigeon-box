package slackaction

import "github.com/slack-go/slack"

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
