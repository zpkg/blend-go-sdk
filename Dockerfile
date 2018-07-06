from golang:1-alpine

WORKDIR /go/src/github.com/blend/go-sdk

RUN apk update && \
    apk upgrade && \
    apk add git

RUN go get -u github.com/lib/pq
RUN go get -u golang.org/x/net/http2
RUN go get -u golang.org/x/oauth2
RUN go get -u golang.org/x/oauth2/google
RUN go get -u golang.org/x/lint/golint

ADD . /go/src/github.com/blend/go-sdk