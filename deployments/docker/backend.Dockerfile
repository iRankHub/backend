FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o backend cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend .

EXPOSE 8080
CMD ["./backend"]