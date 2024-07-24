package slackaction

import (
	"fmt"
	"os"
)

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
