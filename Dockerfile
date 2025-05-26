FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /netpulse ./cmd/netpulse

FROM alpine:3.18
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup
USER appuser
COPY --from=builder /netpulse /usr/local/bin/netpulse
ENTRYPOINT ["netpulse"]