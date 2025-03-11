FROM golang:1.23 as go-build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY . .
RUN sqlc generate

RUN CGO_ENABLED=0 GOOS=linux go build -o app server.go

FROM gcr.io/distroless/static-debian12

WORKDIR /workdir

COPY --from=go-build /app/app /workdir/app

EXPOSE 80

ENV PORT=80

# Run
CMD ["/workdir/app"]
