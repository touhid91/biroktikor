package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/satori/go.uuid"
)

func isSupportedMime(mime string) bool {
	for _, a := range [5]string{"image/gif", "image/jpeg", "image/jpg", "image/png", "image/svg+xml"} {
		if mime == a {
			return true
		}
	}

	return false
}

func b64Enc(command map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range command {
		fmt.Fprintf(b, "%s:%s;", key, value)
	}

	return base64.StdEncoding.EncodeToString([]byte(b.String()))
}

// Sign creates signed url for s3
func Sign(command *PresignInput) (*PresignOutput, error) {
	var (
		region = os.Getenv("AWS_REGION")
		bucket = os.Getenv("AWS_BUCKET")
		expire = 15 * time.Minute
		dir    = "default"
		id     = uuid.Must(uuid.NewV4()).String()
		name   = id
	)

	if valid := isSupportedMime(command.Mime); !valid {
		return nil, BadArgError{
			Arg:     "mime",
			ErrCode: "unsupported_mime_type",
		}
	}

	if "" != command.Key {
		if i := strings.IndexAny(command.Key, "."); 0 < i {
			dir = command.Key[:i]
		} else {
			dir = command.Key
		}
	}

	if "" != command.Meta.Title {
		name = command.Meta.Title
	}

	owner := strconv.FormatFloat(command.Meta.OwnerID, 'f', 6, 64)
	sess, serr := session.NewSession(&aws.Config{Region: aws.String(region)})
	if serr != nil {
		return nil, serr
	}

	svc := s3.New(sess)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		ContentType: aws.String(command.Mime),
		Bucket:      aws.String(bucket),
		Key: aws.String(fmt.Sprintf("%s/%s.%s", dir, b64Enc(map[string]string{
			"created-by": owner,
			"title":      name,
		}), command.Mime[6:])),
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
