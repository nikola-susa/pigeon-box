package slackaction

import "github.com/nikola-susa/pigeon-box/model"

func (s *SlackAction) GetOrCreateUser(slackID string) (*int, error) {
	user, _ := s.Store.GetUserBySlackID(slackID)
	if user != nil {
		return user.ID, nil
	}

	slackUser, err := s.GetSlackUser(slackID)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		SlackID:  slackID,
		Name:     &slackUser.RealName,
		Username: &slackUser.Name,
		Avatar:   &slackUser.Profile.ImageOriginal,
	}
	id, err := s.Store.CreateUser(newUser)
	if err != nil {
		return nil, err
	}
	return id, nil
}
