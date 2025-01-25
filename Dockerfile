FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/app ./cmd/app/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/build/app .

EXPOSE 8080
CMD ["./app"]
