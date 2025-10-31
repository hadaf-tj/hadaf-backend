package cache

import (
	"context"
	"time"
)

// ICache определяет контракт для работы с кэшем (например, Redis).
type ICache interface {
	// Set устанавливает значение по ключу с заданным временем жизни.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get возвращает строковое значение по ключу.
	Get(ctx context.Context, key string) (string, error)

	// Delete удаляет значение по ключу.
	Delete(ctx context.Context, key string) error

	// Increment увеличивает числовое значение по ключу на 1 (например, для счётчиков).
	Increment(ctx context.Context, key string) error
}
