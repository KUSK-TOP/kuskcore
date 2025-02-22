# Build Kusk in a stock Go builder container
FROM golang:1.9-alpine as builder

RUN apk add --no-cache make git

ADD . /go/src/github.com/kusk/kusk
RUN cd /go/src/github.com/kusk/kusk && make kuskd && make kuskcli

# Pull Kusk into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/kusk/kusk/cmd/kuskd/kuskd /usr/local/bin/
COPY --from=builder /go/src/github.com/kusk/kusk/cmd/kuskcli/kuskcli /usr/local/bin/

EXPOSE 1999 46656 46657 9888
