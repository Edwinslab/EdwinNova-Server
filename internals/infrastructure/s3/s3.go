package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client     *s3.Client
	BucketName string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		Client:     client,
		BucketName: bucket,
	}
}

func (s *S3Storage) UploadFile(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	key string,
) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
		Body:   file,
		ContentType: func() *string {
			v := "application/pdf"
			return &v
		}(),
		ACL: "public-read",
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://%s.s3.amazonaws.com/%s",
		s.BucketName,
		key,
	)

	return url, nil
}


func (s *S3Storage) GenerateSignedURL(
	ctx context.Context,
	key string,
) (string, error) {

	presigner := s3.NewPresignClient(s.Client)

	resp, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}