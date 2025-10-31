package rateLimiter

import "context"

// IRateLimiter определяет интерфейс для проверки ограничения по частоте запросов.
type IRateLimiter interface {
	// Allow проверяет, может ли пользователь с данным ключом совершить действие,
	// учитывая лимит и TTL в минутах.
	Allow(ctx context.Context, key string, limit int, ttlMinutes int) (bool, error)

	// ResetAttempts очищает сохраненные попытки
	ResetAttempts(ctx context.Context, key string) error
}
