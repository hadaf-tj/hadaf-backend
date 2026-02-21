package tokens

import (
	"context"
)

type ITokenIssuer interface {
	// IssueTokens создает Access и Refresh токен
	IssueTokens(ctx context.Context, id int, role string) (string, string, error)
}
