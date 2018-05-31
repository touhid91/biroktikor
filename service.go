package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
)

func Presign(command *PresignInput) (*PresignOutput, error) {
	region := "ap-northeast-1"
	bucket := "minin-image"
	expire := 15 * time.Minute
	dir := "default"
	id := uuid.Must(uuid.NewV4()).String()
	name := id

	// TODO Preclause for unsupported mime types
	if !strings.HasPrefix(command.Mime, "image/") {
		return nil, BadArgError{
			Arg:     "mime",
			ErrCode: "unsupported_mime_type",
		}
	}

	if "" != command.Key {
		dir = command.Key
	}

	if "" != command.Meta.Title {
		name = command.Meta.Title
	}

	owner := strconv.FormatFloat(command.Meta.OwnerID, 'f', 6, 64)
	ses, serr := session.NewSession(&aws.Config{Region: aws.String(region)})
	if serr != nil {
		return nil, serr
	}

	svc := s3.New(ses)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		ContentType: aws.String(command.Mime),
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s.%s", dir, id, command.Mime[6:])),
		Metadata: map[string]*string{
			"created-on": aws.String(time.Now().UTC().Format(time.RFC3339)),
			"created-by": aws.String(owner),
			"title":      aws.String(name),
		},
	})

	url, perr := req.Presign(expire)
	if perr != nil {
		return nil, perr
	}

	return &PresignOutput{
		Put: url,
		Get: url[:strings.Index(url, command.Mime[6:])+len(command.Mime[6:])],
	}, nil
}
