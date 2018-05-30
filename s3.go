package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
)

// TODO move to conf
const region = "ap-northeast-1"
const bucket = "minin-image"
const expire = 15 * time.Minute

// PresignInput command for handler
type PresignInput struct {
	mime string
	key  string
}

// Presign create presign url
func Presign(com *PresignInput) (string, error) {
	ses, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return "", err
	}

	svc := s3.New(ses)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		ContentType: aws.String(com.mime),
		Bucket:      aws.String(bucket),
		Key:         aws.String(com.key + "/" + uuid.Must(uuid.NewV4()).String() + "." + com.mime[6:]),
	})
	return req.Presign(expire)
}
