package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
)

type SessionService struct {
	repo            repository.Session
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewSessionService(repo repository.Session, manager auth.TokenManager, accessTTL, refreshTTL time.Duration) *SessionService {
	return &SessionService{
		repo:            repo,
		tokenManager:    manager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *SessionService) SignIn(ctx context.Context, user *models.User) (string, error) {
	// _, accessToken, err := s.tokenManager.NewJWT(user.Id, user.Email, user.Roles, s.accessTokenTTL)
	// if err != nil {
	// 	return "", err
	// }
	// refreshToken, err := s.tokenManager.NewRefreshToken()
	// if err != nil {
	// 	return "", err
	// }

	// accessData := repository.SessionData{
	// 	UserId:      user.Id,
	// 	Roles:       user.Roles,
	// 	AccessToken: accessToken,
	// 	Exp:         s.accessTokenTTL,
	// }
	// if err := s.repo.Create(ctx, user.Id, accessData); err != nil {
	// 	return "", fmt.Errorf("failed to create session. error: %w", err)
	// }

	// refreshData := repository.SessionData{
	// 	UserId:       user.Id,
	// 	Roles:        user.Roles,
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// 	Exp:          s.refreshTokenTTL,
	// }
	// if err := s.repo.Create(ctx, fmt.Sprintf("%s_refresh", user.Id), refreshData); err != nil {
	// 	return "", fmt.Errorf("failed to create session (refresh). error: %w", err)
	// }

	// return accessToken, nil
	return "", nil
}

func (s *SessionService) SingOut(ctx context.Context, userId string) error {
	err := s.repo.Remove(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to delete session. error: %w", err)
	}

	err = s.repo.Remove(ctx, fmt.Sprintf("%s_refresh", userId))
	if err != nil {
		return fmt.Errorf("failed to delete session (refresh). error: %w", err)
	}

	return nil
}

func (s *SessionService) CheckSession(ctx context.Context, u *models.User, token string) (bool, error) {
	user, err := s.repo.Get(ctx, u.Id)
	if err != nil && !errors.Is(err, models.ErrSessionEmpty) {
		return false, fmt.Errorf("failed to get session. error: %w", err)
	}

	refreshUser, err := s.repo.Get(ctx, fmt.Sprintf("%s_refresh", u.Id))
	if err != nil {
		return false, fmt.Errorf("failed to get session (refresh). error: %w", err)
	}

	if user.AccessToken != token && refreshUser.AccessToken != token {
		return false, models.ErrToken
	}

	if user.UserId == "" {
		return true, nil
	}
	return false, nil
}

func (s *SessionService) TokenParse(token string) (user *models.User, err error) {
	claims, err := s.tokenManager.Parse(token)
	if err != nil {
		return nil, err
	}

	// var roles []*models.Role
	// r := claims["roles"].([]interface{})
	// for _, v := range r {
	// 	m := v.(map[string]interface{})
	// 	roles = append(roles, &models.Role{
	// 		Id:      m["id"].(string),
	// 		Service: m["service"].(string),
	// 		Role:    m["role"].(string),
	// 	})
	// }

	// user = &user_api.User{
	// 	Id:    claims["userId"].(string),
	// 	Email: claims["email"].(string),
	// 	Roles: roles,
	// }

	return user, nil
}
