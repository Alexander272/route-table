package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Alexander272/sealur_proto/api/user_api"
	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewJWT(userId, email string, roles []*user_api.Role, ttl time.Duration) (time.Time, string, error)
	Parse(token string) (jwt.MapClaims, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	jwtKey string
}

func NewManager(jwtKey string) (*Manager, error) {
	if strings.Trim(jwtKey, " ") == "" {
		return nil, errors.New("empty jwt key")
	}
	return &Manager{jwtKey}, nil
}

func (m *Manager) NewJWT(userId, email string, roles []*user_api.Role, ttl time.Duration) (iat time.Time, token string, err error) {
	iat = time.Now()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "exp":    iat.Add(ttl).Unix(),
		"iat":    iat.Unix(),
		"userId": userId,
		"email":  email,
		"roles":  roles,
	})
	token, err = newToken.SignedString([]byte(m.jwtKey))
	if err != nil {
		return iat, token, err
	}
	return iat, token, nil
}

func (m *Manager) Parse(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
