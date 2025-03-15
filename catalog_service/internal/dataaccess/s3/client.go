package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	"github.com/minhhoanq/ecommerce/common/logger"
)

type Client interface {
	CreateBucketIfNotExist(ctx context.Context) error
	UploadFile(ctx context.Context, fileName string, fileData []byte) error
	GetFile(ctx context.Context, fileName string) ([]byte, error)
}

type S3Client struct {
	client *s3.Client
	bucket string
	l      logger.Interface
}

func NewClient(cfg configs.Config, l logger.Interface) (Client, error) {
	s3Config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.S3Region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(s3Config)
	return &S3Client{
		client: client,
		bucket: cfg.S3Bucket,
		l:      l,
	}, nil
}

func (s *S3Client) CreateBucketIfNotExist(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})

	if err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			s.l.Info("bucket does not exist, creating...")
			_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: &s.bucket,
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}

			s.l.Info("bucket create successfully")
			return nil
		}

		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	s.l.Info("bucket already exists")
	return nil
}

func (s *S3Client) UploadFile(ctx context.Context, fileName string, fileData []byte) error {
	fmt.Println("process upload image")
	fmt.Println("image: ")

	fileReader := bytes.NewReader(fileData)
	_, err := s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(fileName),
			Body:   fileReader,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *S3Client) GetFile(ctx context.Context, fileName string) ([]byte, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	defer resp.Body.Close()

	// read data from response
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return buf.Bytes(), nil
}
