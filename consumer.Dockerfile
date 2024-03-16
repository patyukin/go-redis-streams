FROM golang:1.21.7-alpine3.19 AS builder

ENV config=docker

WORKDIR /app

COPY . /app

RUN go mod download &&  \
    go mod tidy && \
    go mod download && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon

ENV CONFIG_FILE_PATH=config.yaml

ENTRYPOINT CompileDaemon --build="go build -o bin/consumer cmd/consumer/main.go" --command=./bin/consumer
