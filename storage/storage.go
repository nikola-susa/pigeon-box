package storage

import (
	"github.com/nikola-susa/secret-chat/config"
)

type Storage interface {
	Path(path string) string
	Get(path string, passphrase *string) ([]byte, error)
	GetString(path string, passphrase *string) (string, error)
	Upload(path string, data []byte) (string, error)
	Delete(path string) error
}

func New(c *config.Config) Storage {
	driver := c.File.Driver
	switch driver {
	case "local":
		return NewLocal(c.File.BaseDir, "", *c)
	}
	return Storage(nil)
}
