package model

type Message struct {
	ID        *int    `json:"id,omitempty" db:"id"`
	ThreadID  int     `json:"thread_id,omitempty" db:"thread_id"`
	UserID    int     `json:"user_id,omitempty" db:"user_id"`
	FileID    *int    `json:"file_id,omitempty" db:"file_id"`
	Text      *string `json:"text,omitempty" db:"text"`
	CreatedAt string  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *string `json:"updated_at,omitempty" db:"updated_at"`
	ExpiresAt *string `json:"expires_at,omitempty" db:"expires_at"`
}

type CreateMessageParams struct {
	UserID   int
	ThreadID int
	Text     string
}

type RenderMessage struct {
	ID        string
	ThreadID  string
	Text      string
	CreatedAt string
	Time      string
	User      RenderUser
	File      RenderFile
	IsAuthor  bool
}
