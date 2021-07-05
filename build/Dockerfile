# Builder
FROM golang:1.13.4-alpine as builder

RUN apk update && apk upgrade && \
  apk --update add git gcc make

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
  apk --update --no-cache add tzdata && \
  mkdir /app

WORKDIR /app

EXPOSE 9123

COPY --from=builder /app/VERSION /app
COPY --from=builder /app/charite-data-donation /app

# Add a user to run nginx in non-root mode
RUN addgroup -S charite-data-donation && adduser -S charite-data-donation -G charite-data-donation
RUN chown -R charite-data-donation:charite-data-donation /app/charite-data-donation

USER charite-data-donation

CMD /app/charite-data-donation
