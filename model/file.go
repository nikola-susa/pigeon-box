package model

type File struct {
	ID          *int    `json:"id,omitempty" db:"id"`
	ThreadID    *int    `json:"thread_id,omitempty" db:"thread_id"`
	UserID      *int    `json:"user_id,omitempty" db:"user_id"`
	Name        string  `json:"name,omitempty" db:"name"`
	Path        *string `json:"path,omitempty" db:"path"`
	ContentType *string `json:"content_type,omitempty" db:"content_type"`
	Size        *int64  `json:"size,omitempty" db:"size"`
	CreatedAt   *string `json:"created_at,omitempty" db:"created_at"`
}

type RenderFile struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Size        string `json:"size,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	ThreadHash  string `json:"thread_hash,omitempty"`
}
