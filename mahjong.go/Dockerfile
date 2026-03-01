FROM golang:1.21-alpine AS builder

WORKDIR /app

# cache modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mahjong-server ./cmd/... 

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/mahjong-server .

EXPOSE 8080

CMD ["./mahjong-server"]