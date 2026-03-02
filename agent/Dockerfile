FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git ca-certificates

COPY go.mod ./
RUN go mod tidy 2>/dev/null || true
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o agent ./cmd/agent

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 65534 appuser
WORKDIR /app
COPY --from=builder /build/agent .
USER appuser
EXPOSE 8080
ENTRYPOINT ["./agent"]
