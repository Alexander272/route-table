package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
	"github.com/google/uuid"
)

type SessionService struct {
	repo            repository.Session
	user            *UserService
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewSessionService(repo repository.Session, user *UserService, manager auth.TokenManager, accessTTL, refreshTTL time.Duration) *SessionService {
	return &SessionService{
		repo:            repo,
		user:            user,
		tokenManager:    manager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *SessionService) SignIn(ctx context.Context, u models.SignIn) (models.User, string, error) {
	user, err := s.user.Get(ctx, u)
	if err != nil {
		return models.User{}, "", err
	}

	return s.Refresh(ctx, user)
}

func (s *SessionService) Refresh(ctx context.Context, user models.UserWithRole) (models.User, string, error) {
	_, accessToken, err := s.tokenManager.NewJWT(user.Id.String(), user.Role, s.accessTokenTTL)
	if err != nil {
		return models.User{}, "", err
	}
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return models.User{}, "", err
	}

	accessData := models.SessionData{
		UserId:      user.Id.String(),
		Role:        user.Role,
		AccessToken: accessToken,
		Exp:         s.accessTokenTTL,
	}
	if err := s.repo.Create(ctx, refreshToken, accessData); err != nil {
		return models.User{}, "", fmt.Errorf("failed to create session. error: %w", err)
	}

	refreshData := models.SessionData{
		UserId:       user.Id.String(),
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          s.refreshTokenTTL,
	}
	if err := s.repo.Create(ctx, fmt.Sprintf("%s_refresh", user.Id), refreshData); err != nil {
		return models.User{}, "", fmt.Errorf("failed to create session (refresh). error: %w", err)
	}

	retUser := models.User{
		Id:   user.Id,
		Role: user.Role.Role,
	}

	return retUser, accessToken, nil
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

func (s *SessionService) CheckSession(ctx context.Context, u models.UserWithRole, token string) (bool, error) {
	refreshUser, err := s.repo.Get(ctx, fmt.Sprintf("%s_refresh", u.Id.String()))
	if err != nil {
		return false, fmt.Errorf("failed to get session (refresh). error: %w", err)
	}

	user, err := s.repo.Get(ctx, refreshUser.RefreshToken)
	if err != nil && !errors.Is(err, models.ErrSessionEmpty) {
		return false, fmt.Errorf("failed to get session. error: %w", err)
	}

	if user.AccessToken != token && refreshUser.AccessToken != token {
		return false, models.ErrToken
	}

	if user.UserId == "" {
		return true, nil
	}
	return false, nil
}

func (s *SessionService) TokenParse(token string) (user models.UserWithRole, err error) {
	claims, err := s.tokenManager.Parse(token)
	if err != nil {
		return user, err
	}

	r := claims["role"].(map[string]interface{})

	roleId, err := uuid.Parse(r["id"].(string))
	if err != nil {
		return user, fmt.Errorf("failed to parse uuid. error: %w", err)
	}

	op := r["operations"].([]interface{})
	oprations := make([]string, 0, len(op))
	for _, v := range op {
		oprations = append(oprations, v.(string))
	}

	role := models.Role{
		Id:         roleId,
		Title:      r["title"].(string),
		Role:       r["role"].(string),
		Operations: oprations,
	}

	userId, err := uuid.Parse(claims["userId"].(string))
	if err != nil {
		return user, fmt.Errorf("failed to parse uuid. error: %w", err)
	}

	user = models.UserWithRole{
		Id:   userId,
		Role: role,
	}

	return user, nil
}
