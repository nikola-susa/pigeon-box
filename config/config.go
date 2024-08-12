package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
)

type Config struct {
	Version string `envconfig:"VERSION" default:"0.1"`
	Public  struct {
		URL string `envconfig:"PUBLIC_URL" default:"http://localhost:8081"`
	}
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"localhost"`
		Port int    `envconfig:"SERVER_PORT" default:"8081"`
	}
	Slack struct {
		AppToken          string `envconfig:"SLACK_APP_TOKEN" required:"true"`
		BotToken          string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
		VerificationToken string `envconfig:"SLACK_BOT_VERIFICATION_TOKEN" required:"true"`
	}
	Database struct {
		URL string `envconfig:"DATABASE_URL" default:"file:./local.db"`
	}
	File struct {
		BaseDir string `envconfig:"FILE_BASE_DIR" default:"./bucket"`
		Driver  string `envconfig:"FILE_DRIVER" default:"local"`
		MaxSize int64  `envconfig:"FILE_MAX_SIZE" default:"32"`
	}
	Crypt struct {
		Passphrase string `envconfig:"CRYPT_PASSPHRASE" required:"true"`
		HashSalt   string `envconfig:"CRYPT_HASH_SALT" required:"true"`
		HashLength int    `envconfig:"CRYPT_HASH_LENGTH" default:"12"`
	}
}

func NewConfig(file string) (*Config, error) {
	if file != "" {
		err := godotenv.Overload(file)
		if err != nil {
			return nil, err
		}
	}

	absFile, _ := filepath.Abs(file)
	_, err := os.Stat(absFile)
	fileNotExists := os.IsNotExist(err)
	if fileNotExists {
		return nil, err
	}

	_ = godotenv.Overload(absFile)

	config := new(Config)
	err = envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
