FROM golang:latest

ENV SERVER_HOST "0.0.0.0"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]