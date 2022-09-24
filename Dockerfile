# ---- Build ----

FROM golang:1.19.1-bullseye as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o app

# ---- Release ----

FROM debian:bullseye-slim as deploy
RUN apt-get update
COPY --from=builder /app/app .
CMD ["./app"]

# ---- Develop ----

FROM golang:1.19.1-bullseye as develop
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
