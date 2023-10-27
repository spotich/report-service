FROM golang:latest

WORKDIR /reporter

COPY ./ ./

RUN go mod download

CMD ["go", "run", "cmd/server.go"]
