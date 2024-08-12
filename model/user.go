package model

type User struct {
	ID        *int    `json:"id,omitempty" db:"id"`
	Name      *string `json:"name,omitempty" db:"name"`
	Username  *string `json:"username" db:"username"`
	SlackID   string  `json:"slack_id,omitempty" db:"slack_id"`
	Avatar    *string `json:"avatar,omitempty" db:"avatar"`
	CreatedAt *string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *string `json:"updated_at,omitempty" db:"updated_at"`
}

type RenderUser struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Username string  `json:"username,omitempty"`
	Avatar   string  `json:"avatar,omitempty"`
	SlackID  *string `json:"slack_id,omitempty"`
}
