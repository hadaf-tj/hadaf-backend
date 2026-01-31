FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/shb

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN mkdir -p /app/logs

# Копируем исполняемый файл
COPY --from=builder /app/main .

# Копируем конфиги
COPY --from=builder /app/internal/configs ./configs

# Копируем docs
COPY --from=builder /app/docs ./docs

EXPOSE 8000

CMD ["./main"]