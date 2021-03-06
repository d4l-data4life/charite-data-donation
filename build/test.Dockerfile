# Tester
FROM golang:1.13.4-alpine

RUN apk update && apk upgrade && \
  apk --update --no-cache add git gcc musl-dev make tzdata && \
  mkdir /app

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD GOOS=linux GOARCH=amd64 make unit-test
