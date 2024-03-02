FROM golang:1.21.7-alpine3.19 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o ./bin/sender cmd/sender/main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/sender .
COPY config.yaml config.yaml
ENV CONFIG_FILE_PATH=config.yaml

CMD ["./sender"]