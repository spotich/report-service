FROM golang:latest

WORKDIR /timer

COPY ./ ./

RUN go mod download

CMD ["go", "run", "cmd/server.go"]
