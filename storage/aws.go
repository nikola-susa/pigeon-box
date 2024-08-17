package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nikola-susa/pigeon-box/config"
	"github.com/nikola-susa/pigeon-box/crypt"
	"io"
	"path/filepath"
	"strings"
)

type AWSStorage struct {
	BucketName string
	PathPrefix string
	Config     config.Config
	S3Client   *s3.Client
}

func NewAWS(pathPrefix string, conf config.Config) *AWSStorage {

	if conf.File.AWS.BucketName == "" {
		fmt.Print("AWS_BUCKET_NAME is required when using AWS storage driver")
		return nil
	}

	if conf.File.AWS.AccessKeyID == "" {
		fmt.Print("AWS_ACCESS_KEY_ID is required when using AWS storage driver")
		return nil
	}

	if conf.File.AWS.SecretAccessKey == "" {
		fmt.Print("AWS_SECRET_ACCESS_KEY is required when using AWS storage driver")
		return nil
	}

	if conf.File.AWS.EndpointURL == "" {
		fmt.Print("AWS_ENDPOINT_URL_S3 is required when using AWS storage driver")
		return nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(conf.File.AWS.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.File.AWS.AccessKeyID, conf.File.AWS.SecretAccessKey, "")),
	)

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.File.AWS.EndpointURL)
	})

	if err != nil {
		fmt.Println("unable to load SDK config, " + err.Error())
		return nil
	}

	s := &AWSStorage{
		BucketName: conf.File.AWS.BucketName,
		PathPrefix: pathPrefix,
		Config:     conf,
		S3Client:   client,
	}
	return s
}

func (s *AWSStorage) Path(path string) string {
	return filepath.Join(s.PathPrefix, path)
}

func (s *AWSStorage) Get(path string, passphrase *string) ([]byte, error) {
	key := path
	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(result.Body)

	data, err := io.ReadAll(result.Body)
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

func (s *AWSStorage) GetString(path string, passphrase *string) (string, error) {
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

func (s *AWSStorage) Upload(path string, data []byte) (string, error) {
	key := s.Path(path + ".blob")

	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *AWSStorage) Delete(path string) error {
	key := path

	_, err := s.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	return nil
}
