FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v github.com/aws/aws-sdk-go github.com/satori/go.uuid