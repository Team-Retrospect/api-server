FROM golang:alpine

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/gocql/gocql

COPY . .

RUN go build -o main .

CMD ["./main"]