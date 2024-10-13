FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate

RUN CGO_ENABLED=0 GOOS=linux go build -o app server.go

EXPOSE 80

ENV PORT=80

# Run
CMD ["./app"]