FROM golang:1.14.2-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/plopezm/cloud-kaiser

COPY go.mod go.sum ./
COPY core core
COPY query-service query-service
COPY repository-service repository-service
COPY pusher-service pusher-service
COPY kaiser-service kaiser-service

RUN go install ./...

FROM alpine:3.11
WORKDIR /usr/bin
RUN mkdir -p /usr/bin/logs && chmod -R a+rwx /usr/bin/logs
COPY --from=build /go/bin .
