package tokens

import (
	"context"
)

type ITokenIssuer interface {
	// IssueTokens создает Access и Refresh токен
	IssueTokens(ctx context.Context, id int, role string, isApproved bool) (string, string, error)
}
