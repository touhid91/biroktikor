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
const bucket = "minin-bucket"
const expire = 15 * time.Minute

// PresignInput command for handler
type PresignInput struct {
	mimetype  string
	directory string
}

// Presign create presign url
func Presign(com *PresignInput) (string, error) {
	ses, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return "", err
	}

	svc := s3.New(ses)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		ContentType: aws.String(com.mimetype),
		Bucket:      aws.String(bucket),
		Key:         aws.String(com.directory + "/" + uuid.Must(uuid.NewV4()).String() + "." + com.mimetype[6:]),
	})
	return req.Presign(expire)
}