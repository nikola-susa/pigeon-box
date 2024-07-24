package storage

import (
	"github.com/nikola-susa/secret-chat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLocalReturnsLocalStorageWithGivenValues(t *testing.T) {
	baseDir := "/tmp"
	pathPrefix := "/prefix"
	conf := config.Config{}

	s := NewLocal(baseDir, pathPrefix, conf)

	assert.Equal(t, baseDir, s.BaseDir)
	assert.Equal(t, pathPrefix, s.PathPrefix)
	assert.Equal(t, conf, s.config)
}

func TestPathReturnsCorrectPath(t *testing.T) {
	s := &LocalStorage{
		BaseDir:    "/tmp",
		PathPrefix: "/prefix",
	}

	path := s.Path("file.txt")

	assert.Equal(t, "/tmp/prefix/file.txt", path)
}

func TestLocalStorage_Upload(t *testing.T) {
	s := &LocalStorage{
		BaseDir:    "../tmp",
		PathPrefix: "/",
	}

	path, err := s.Upload("file.txt", []byte("content"))

	assert.Nil(t, err)
	assert.Equal(t, "../tmp/file.txt.blob", path)
}

func TestLocalStorage_Delete(t *testing.T) {
	s := &LocalStorage{
		BaseDir:    "../tmp",
		PathPrefix: "/",
	}

	err := s.Delete("file.txt.blob")

	assert.Nil(t, err)
}
