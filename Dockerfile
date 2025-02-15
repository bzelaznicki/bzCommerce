FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Then copy the rest of the code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY html/ html/
EXPOSE 8080
CMD ["./main"]