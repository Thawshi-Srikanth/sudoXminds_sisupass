package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/tokens"
	"sisupass.com/sisupass/internal/types"
)

type OAuthService struct {
	models      *data.Models
	oauthConfig *oauth2.Config
	frontendURL string
}

func NewOAuthService(models *data.Models, oauthConfig *oauth2.Config, frontendURL string) *OAuthService {
	return &OAuthService{
		models:      models,
		oauthConfig: oauthConfig,
		frontendURL: frontendURL,
	}
}

type GoogleOAuthRequest struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type GoogleOAuthResponse struct {
	RedirectURL string       `json:"redirect_url"`
	User        UserResponse `json:"user"`
	Token       string       `json:"token"`
}

func (s *OAuthService) InitiateGoogleOAuth(ctx context.Context, state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *OAuthService) ProcessGoogleOAuthCallback(ctx context.Context, req GoogleOAuthRequest) (*GoogleOAuthResponse, error) {
	token, err := s.oauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange OAuth code: %w", err)
	}

	googleUser, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	existingUser, err := s.models.Users.GetByEmail(googleUser.Email)
	if err == nil {
		authToken, err := s.models.Tokens.New(existingUser.ID, 24*time.Hour, tokens.ScopeAuthentication)
		if err != nil {
			return nil, fmt.Errorf("failed to generate authentication token: %w", err)
		}

		return &GoogleOAuthResponse{
			RedirectURL: fmt.Sprintf("%s/auth/callback?token=%s", s.frontendURL, authToken.Plaintext),
			User: UserResponse{
				ID:        existingUser.ID,
				Username:  existingUser.UserName,
				Email:     existingUser.Email,
				UserType:  existingUser.UserType,
				Activated: existingUser.Activated,
				CreatedAt: existingUser.CreatedAt,
			},
			Token: authToken.Plaintext,
		}, nil
	} else if err != types.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	user, err := s.models.Users.FindOrCreateFromGoogle(googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	authToken, err := s.models.Tokens.New(user.ID, 24*time.Hour, tokens.ScopeAuthentication)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authentication token: %w", err)
	}

	return &GoogleOAuthResponse{
		RedirectURL: fmt.Sprintf("%s/auth/callback?token=%s", s.frontendURL, authToken.Plaintext),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.UserName,
			Email:     user.Email,
			UserType:  user.UserType,
			Activated: user.Activated,
			CreatedAt: user.CreatedAt,
		},
		Token: authToken.Plaintext,
	}, nil
}

func (s *OAuthService) getGoogleUserInfo(accessToken string) (*types.GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user types.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
