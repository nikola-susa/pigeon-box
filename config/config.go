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
		Port string `envconfig:"SERVER_PORT" default:"8081"`
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
		MaxSize string `envconfig:"FILE_MAX_SIZE" default:"32"`
		BaseDir string `envconfig:"FILE_BASE_DIR" default:"./bucket"`
		Driver  string `envconfig:"FILE_DRIVER" default:"local"`
		Local   struct {
			PathPrefix string `envconfig:"FILE_LOCAL_PATH_PREFIX" default:"/"`
			BaseDir    string `envconfig:"FILE_LOCAL_BASE_DIR" default:"./bucket"`
		}
		AWS struct {
			AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID" default:""`
			SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" default:""`
			EndpointURL     string `envconfig:"AWS_ENDPOINT_URL_S3" default:""`
			Region          string `envconfig:"AWS_REGION" default:"auto"`
			BucketName      string `envconfig:"AWS_BUCKET_NAME" default:""`
			NamePrefix      string `envconfig:"AWS_NAME_PREFIX" default:""`
		}
	}
	Crypt struct {
		Passphrase string `envconfig:"CRYPT_PASSPHRASE" required:"true"`
		HashSalt   string `envconfig:"CRYPT_HASH_SALT" required:"true"`
		HashLength string `envconfig:"CRYPT_HASH_LENGTH" default:"12"`
	}
}

func NewConfig(file string) (*Config, error) {
	if file != "" {
		absFile, _ := filepath.Abs(file)
		_, err := os.Stat(absFile)
		if err == nil {
			err := godotenv.Overload(file)
			if err != nil {
				return nil, err
			}
		}
	}

	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
