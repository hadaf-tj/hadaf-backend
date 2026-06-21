package oauth

import (
	"context"
	"shb/internal/configs"
	"shb/internal/models"
	"strings"

	"golang.org/x/oauth2"
	gg "golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	cfg *configs.OAuthProviderConfig
}

func NewGoogleProvider(cfg *configs.OAuthProviderConfig) *GoogleProvider {
	return &GoogleProvider{cfg: cfg}
}

func (GoogleProvider) ProviderName() string {
	return "google"
}

func (g GoogleProvider) CallbackPath() string {
	parts := strings.Split(g.cfg.RedirectURL, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (g GoogleProvider) OAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.cfg.ClientID,
		ClientSecret: g.cfg.ClientSecret,
		Scopes:       []string{goauth2.UserinfoProfileScope, goauth2.UserinfoEmailScope},
		Endpoint:     gg.Endpoint,
		RedirectURL:  g.cfg.RedirectURL,
	}
}

func (g GoogleProvider) GetUser(ctx context.Context, tok *oauth2.Token) (models.OAuthUserInfo, error) {
	client := g.OAuth2Config().Client(ctx, tok)

	service, err := goauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return models.OAuthUserInfo{}, err
	}

	userinfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return models.OAuthUserInfo{}, err
	}

	emailVerified := userinfo.VerifiedEmail != nil && *userinfo.VerifiedEmail

	return models.OAuthUserInfo{
		ID:                userinfo.Id,
		Username:          userinfo.Name,
		Email:             userinfo.Email,
		EmailVerified:     emailVerified,
		OAuthProviderName: g.ProviderName(),
		AvatarURL:         &userinfo.Picture,
	}, nil
}
