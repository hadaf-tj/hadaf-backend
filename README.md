# 1. Делаешь файл настроек
cp .env_example .env

# 2. Поднимаешь базу и кеш (нужен Docker)
docker-compose up -d

# 3. Скачиваешь библиотеки и запускаешь сервер
go mod download
go run cmd/shb/main.go
