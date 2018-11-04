FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/plopezm/cloud-kaiser

COPY glide.lock glide.yaml ./
COPY vendor vendor
COPY core core
COPY query-service query-service
COPY repository-service repository-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .