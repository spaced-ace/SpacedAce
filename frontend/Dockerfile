FROM node:18-slim as tailwind

WORKDIR /workdir

COPY package.json package.json
COPY package-lock.json package-lock.json
COPY tailwind.config.js tailwind.config.js

RUN npm install

COPY views views
COPY static/main.css static/main.css

RUN npx tailwindcss -i ./static/main.css -o ./static/tailwindcss.css


FROM golang:1.22 as go-build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY api api
COPY auth auth
COPY cmd cmd
COPY constants constants
COPY context context
COPY models models
COPY render render
COPY service service
COPY utils utils
COPY views views

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

RUN ls

FROM gcr.io/distroless/static-debian12

WORKDIR /workdir

COPY --from=go-build /app/app /workdir/app

COPY --from=tailwind /workdir/static static
COPY static/images static/images
COPY static/icons static/icons
COPY static/htmx.min.js static/htmx.min.js

EXPOSE 80

ENV PORT=80

# Run
CMD ["/workdir/app"]