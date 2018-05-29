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
		mimetype:  r.FormValue("mimetype"),
		directory: r.FormValue("directory")}
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

	result, err := svc.ListBuckets(nil)
	if err != nil {
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("minin-bucket")})
	if err != nil {
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
	io.WriteString(w, "Hello bitch")
}

func main() {
	http.HandleFunc("/s3", s3Handler)
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
