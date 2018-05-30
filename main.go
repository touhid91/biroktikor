package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func s3Handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	command := PresignInput{
		mime: r.FormValue("mime"),
		key:  r.FormValue("key"),
	}

	if len(command.mime) < 1 {
		io.WriteString(w, "param missing: mime")
	}
	fmt.Print(command)

	url, err := Presign(&command)
	if nil != err {
		fmt.Print(err)
	}

	io.WriteString(w, url)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ses, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
	}

	svc := s3.New(ses)

	res, _ := svc.GetBucketAcl(&s3.GetBucketAclInput{
		Bucket: aws.String(bucket),
	})

	fmt.Print(res, err)

	io.WriteString(w, "Hello ")
}

func main() {
	http.HandleFunc("/s3", s3Handler)
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
