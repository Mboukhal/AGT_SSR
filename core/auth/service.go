package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

const (
	SessionName = "session"
	UserIDKey   = "user_id"
	EmailKey    = "email"
	UsernameKey = "username"
)

type Service struct {
	Store    *sessions.CookieStore
	OAuthCfg *oauth2.Config
}

func NewService(store *sessions.CookieStore) *Service {
	tenantID := os.Getenv("MICROSOFT_TENANT_ID")
	if tenantID == "" {
		tenantID = "common"
	}

	return &Service{
		Store: store,
		OAuthCfg: &oauth2.Config{
			ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
			ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("MICROSOFT_REDIRECT_URL"),
			Scopes:       []string{"openid", "profile", "email", "User.Read"},
			Endpoint:     microsoft.AzureADEndpoint(tenantID),
		},
	}
}

func (s *Service) AuthURL(state string) string {
	return s.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *Service) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.OAuthCfg.Exchange(ctx, code)
}

func (s *Service) GetMicrosoftUser(ctx context.Context, token *oauth2.Token) (*MicrosoftUser, error) {
	client := s.OAuthCfg.Client(ctx, token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Microsoft API error: %s", string(body))
	}

	var user MicrosoftUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

type MicrosoftUser struct {
	ID          string `json:"id"`
	Email       string `json:"mail"`
	UserEmail   string `json:"userPrincipalName"`
	DisplayName string `json:"displayName"`
}

func (m *MicrosoftUser) EmailOrUPN() string {
	if m.Email != "" {
		return m.Email
	}
	return m.UserEmail
}

func (s *Service) CreateSession(r *http.Request, user *MicrosoftUser) (*sessions.Session, error) {
	session, err := s.Store.New(r, SessionName)
	if err != nil {
		return nil, err
	}
	session.Values[UserIDKey] = user.ID
	session.Values[EmailKey] = user.EmailOrUPN()
	session.Values[UsernameKey] = user.DisplayName
	return session, nil
}

func (s *Service) GetUserFromSession(session *sessions.Session) (string, string, bool) {
	userID, ok := session.Values[UserIDKey].(string)
	if !ok {
		return "", "", false
	}
	username, _ := session.Values[UsernameKey].(string)
	return userID, username, true
}
