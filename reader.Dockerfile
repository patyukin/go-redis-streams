FROM golang:latest

WORKDIR /app

COPY docker .

RUN go build -o reader .

CMD ["./reader"]