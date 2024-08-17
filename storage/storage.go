package storage

import (
	"github.com/nikola-susa/pigeon-box/config"
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
		return NewLocal(c.File.Local.BaseDir, c.File.Local.PathPrefix, *c)
	case "aws":
		return NewAWS(c.File.AWS.NamePrefix, *c)
	}
	return Storage(nil)
}
