FROM node:18-slim as tailwind

WORKDIR /workdir

COPY package.json package.json
COPY package-lock.json package-lock.json
COPY tailwind.config.js tailwind.config.js

RUN npm install

COPY views views
COPY static/main.css static/main.css

RUN npx tailwindcss -i ./static/main.css -o ./static/tailwindcss.css


FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd
COPY views views

COPY --from=tailwind /workdir/static static
COPY static/images static/images
COPY static/icons static/icons

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

EXPOSE 80

ENV PORT=80

# Run
CMD ["./app"]