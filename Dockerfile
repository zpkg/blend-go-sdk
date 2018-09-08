from golang:1-alpine

WORKDIR /go/src/github.com/blend/go-sdk

RUN apk update && \
    apk upgrade && \
    apk add git

RUN go get ./...

ADD . /go/src/github.com/blend/go-sdk
