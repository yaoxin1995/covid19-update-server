FROM golang:latest

ENV SERVER_HOST "0.0.0.0"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]