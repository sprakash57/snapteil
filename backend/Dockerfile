FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs
RUN mkdir -p uploads

EXPOSE 5000
CMD ["./server"]
