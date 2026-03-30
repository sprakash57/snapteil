FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/data ./data
COPY --from=builder /app/uploads ./uploads

EXPOSE 4000

CMD ["./server"]
