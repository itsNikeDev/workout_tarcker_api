FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server
RUN CGO_ENABLED=0 go build -o migrate ./cmd/migrate

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
EXPOSE 8080
CMD ["./server"]
