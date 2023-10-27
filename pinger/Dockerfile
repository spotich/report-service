FROM golang:latest

WORKDIR /pinger

COPY ./ ./

RUN go mod download

CMD ["go", "run", "cmd/client.go"]
