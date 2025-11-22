FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o review-assigner ./cmd/review-assigner

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/review-assigner /review-assigner

EXPOSE 8080

CMD ["/review-assigner"]
