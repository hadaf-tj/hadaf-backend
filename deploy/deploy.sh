#!/usr/bin/env bash
set -e

# Перейти в папку приложения
APP_DIR="$(pwd)"
IMAGE_NAME="myapp"
TAG=$(git rev-parse --short HEAD 2>/dev/null || echo "local")

# Собрать docker image
docker build -t ${IMAGE_NAME}:${TAG} .

# На случай старого контейнера
if docker ps -a --format '{{.Names}}' | grep -q "^${IMAGE_NAME}$"; then
  docker stop ${IMAGE_NAME} || true
  docker rm ${IMAGE_NAME} || true
fi

# Запустить контейнер
# --env-file использует .env в каталоге разворачивания
docker run -d \
  --name ${IMAGE_NAME} \
  --env-file .env \
  -p 8000:8000 \
  --network host \
  --restart unless-stopped \
  ${IMAGE_NAME}:${TAG}

# Очистка старых образов (необязательно)
# docker image prune -f

echo "Deployed ${IMAGE_NAME}:${TAG}"