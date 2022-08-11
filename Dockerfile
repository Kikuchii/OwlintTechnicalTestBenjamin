# syntax=docker/dockerfile:1

FROM postgres:latest

FROM golang:1.19.0-alpine

WORKDIR /

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

COPY  /pkg/*go ./pkg/
COPY  /translate/*go ./translate/

RUN go build -o technicaltestAPI

EXPOSE 8080

CMD ["./technicaltestAPI"]