FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/redis-set-ttl-for-key-pattern
COPY . /go/src/github.com/octoblu/redis-set-ttl-for-key-pattern

RUN env CGO_ENABLED=0 go build -o redis-set-ttl-for-key-pattern -a -ldflags '-s' .

CMD ["./redis-set-ttl-for-key-pattern"]
