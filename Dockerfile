FROM golang:alpine AS build-env

# Install certs
RUN apk add --no-cache ca-certificates

ADD . /go/src/github.com/CoreDumped-ETSISI/etsisi-telegram-bot

RUN cd /go/src/github.com/CoreDumped-ETSISI/etsisi-telegram-bot && go build -o app

FROM alpine

# Import certificates
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

WORKDIR /app
COPY --from=build-env /go/src/github.com/CoreDumped-ETSISI/etsisi-telegram-bot/app /app/
ENTRYPOINT ./app
