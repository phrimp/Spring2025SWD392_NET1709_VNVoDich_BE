package services

import (
	"context"
	"encoding/json"
	"google-service/internal/config"
	"google-service/internal/models"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthService struct {
	config       *config.GoogleOAuthConfig
	oauth2Config *oauth2.Config
}

func NewGoogleOAuthService(config *config.GoogleOAuthConfig) *GoogleOAuthService {
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint:     google.Endpoint,
	}
	return &GoogleOAuthService{
		config:       config,
		oauth2Config: &oauth2Config,
	}
}

func (s *GoogleOAuthService) GetAuthURL(state string) string {
	return s.oauth2Config.AuthCodeURL(state)
}

func (s *GoogleOAuthService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth2Config.Exchange(ctx, code)
}

func (s *GoogleOAuthService) GetUserInfo(token *oauth2.Token) (*models.UserInfo, error) {
	client := s.oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userinfo models.UserInfo
	if err := json.Unmarshal(data, &userinfo); err != nil {
		return nil, err
	}
	return &userinfo, nil
}
