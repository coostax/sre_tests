#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o bin/key-server.go main.go

#App stage
FROM alpine:latest
LABEL Name=sretakehome Version=0.0.1
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin /app

CMD exec /app/key-server.go --max-size $MAX_SIZE --srv-port $PORT
