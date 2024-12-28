FROM golang:1.23.4-alpine3.21 as builder 
workdir /app
copy . .

# Build the binary
run go build -o main .

# Path: Dockerfile
FROM alpine:3.19
workdir /app
copy --from=builder /app/main .
copy --from=builder /app/.env .
cmd ["./main"]
