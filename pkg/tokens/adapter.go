package tokens

import (
	"context"
	"shb/internal/models"
)

type ITokenIssuer interface {
	// IssueTokens создает Access и Refresh токен
	IssueTokens(ctx context.Context, id int, role string) (string, string, error)
	// VerifyToken проверяет токен и возвращает claims (для Access и Refresh)
	VerifyToken(ctx context.Context, token string) (*models.CustomClaims, error)
}
