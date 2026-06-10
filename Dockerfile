# --- Build stage ---
FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /task-api ./cmd/server

# --- Runtime stage ---
FROM alpine:3.21

RUN adduser -D -u 1001 app \
    && mkdir -p /data/uploads \
    && chown -R app /data/uploads
USER app
WORKDIR /home/app

COPY --from=build /task-api ./task-api

EXPOSE 8080
CMD ["./task-api"]
