FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git ca-certificates

COPY . .
RUN go mod download -x 2>&1 | tail -3 || true
RUN CGO_ENABLED=0 GOOS=linux GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod \
    go build -ldflags="-w -s" -o agent ./cmd/agent

FROM golang:1.22-alpine AS prod
WORKDIR /app
COPY --from=builder /build/agent .
EXPOSE 8080
ENTRYPOINT ["./agent"]
