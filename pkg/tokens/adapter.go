package tokens

import (
	"context"
	"shb/internal/models"
)

type ITokenIssuer interface {
	// IssueTokens создает Access и Refresh токен
	IssueTokens(ctx context.Context, user *models.User) (string, string, error)
}
