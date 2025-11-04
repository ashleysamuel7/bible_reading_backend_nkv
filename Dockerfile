# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
ARG BUILD_TARGET=./main.go
WORKDIR /src

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
# Build to a distinct filename so it never collides with a source directory
RUN go build -o /src/server-bin ${BUILD_TARGET}

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary (explicit filename)
COPY --from=builder /src/server-bin /app/server

# Make sure it is executable and owned by non-root user
RUN chmod +x /app/server \
 && addgroup -S app \
 && adduser -S -G app app \
 && chown app:app /app/server

ENV TZ=UTC

USER app

EXPOSE 8080

CMD ["./server"]
