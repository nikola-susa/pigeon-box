package model

type Channel struct {
	ID              *int    `json:"id,omitempty" db:"id"`
	SlackID         string  `json:"slack_id,omitempty" db:"slack_id"`
	TeamID          string  `json:"team_id,omitempty" db:"team_id"`
	Passphrase      *string `json:"passphrase,omitempty" db:"passphrase"`
	PassphraseSetAt *string `json:"passphrase_set_at,omitempty" db:"passphrase_set_at"`
}
