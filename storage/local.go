package storage

import (
	"bytes"
	"fmt"
	"github.com/nikola-susa/secret-chat/config"
	"github.com/nikola-susa/secret-chat/crypt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct {
	BaseDir    string
	PathPrefix string
	MkdirPerm  os.FileMode
	WritePerm  os.FileMode
	ExpireTime time.Duration
	config     config.Config
}

func NewLocal(baseDir string, pathPrefix string, conf config.Config) *LocalStorage {
	if pathPrefix == "" {
		pathPrefix = "/"
	}
	s := &LocalStorage{
		BaseDir:    baseDir,
		PathPrefix: pathPrefix,
		MkdirPerm:  0755,
		WritePerm:  0644,
		config:     conf,
	}
	return s
}

func (s *LocalStorage) Path(path string) string {
	return filepath.Join(s.BaseDir, s.PathPrefix, path)
}

func (s *LocalStorage) Get(path string, passphrase *string) ([]byte, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(path, ".blob") {
		dByt, err := crypt.Decrypt(*passphrase, data)
		if err != nil {
			return nil, err
		}
		data = dByt
	}

	return data, nil
}

func (s *LocalStorage) GetString(path string, passphrase *string) (string, error) {
	data, err := s.Get(path, passphrase)
	if err != nil {
		return "", err
	}

	str, err := DecodeBytes(data)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (s *LocalStorage) Upload(path string, data []byte) (string, error) {
	fp := s.Path(path + ".blob")

	if err := os.MkdirAll(filepath.Dir(fp), s.MkdirPerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	o, err := os.Create(fp)
	if err != nil {
		return "", err
	}

	defer func(o *os.File) {
		err := o.Close()
		if err != nil {
			return
		}
	}(o)

	_, err = io.Copy(o, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	return fp, nil
}

func (s *LocalStorage) Delete(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
