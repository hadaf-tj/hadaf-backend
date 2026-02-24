FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/shb

# --- Stage 2: Runner ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# 1. Создаем системную группу и пользователя без root-прав (назовем его hadaf)
RUN addgroup -S hadafgroup && adduser -S hadafuser -G hadafgroup

WORKDIR /app

# 2. Создаем папку для логов
RUN mkdir -p /app/logs

# Копируем исполняемый файл и конфиги
COPY --from=builder /app/main .
COPY --from=builder /app/internal/configs ./configs
COPY --from=builder /app/docs ./docs

# 3. КРИТИЧЕСКИ ВАЖНО: Отдаем права на папку /app нашему новому пользователю.
# Без этого приложение на Go упадет с ошибкой "permission denied", 
# когда попытается записать файл в /app/logs
RUN chown -R hadafuser:hadafgroup /app

# 4. Переключаемся на безопасного пользователя
USER hadafuser

EXPOSE 8000

CMD ["./main"]