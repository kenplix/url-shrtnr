package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// TokenService provides logic for JWT & Refresh tokens generation and parsing.
//
//go:generate mockery --dir . --name TokenService --output ./mocks
type TokenService interface {
	CreateTokens(id string) (Tokens, error)
	ParseAccessToken(token string) (string, error)
	ParseRefreshToken(token string) (string, error)
}

type Service struct {
	accessTokenSigningKey  string
	accessTokenTTL         time.Duration
	refreshTokenSigningKey string
	refreshTokenTTL        time.Duration
}

func NewTokenService(config Config) (*Service, error) {
	s := Service{
		accessTokenTTL:  defaultAccessTokenTTL,
		refreshTokenTTL: defaultRefreshTokenTTL,
	}
	if err := SetConfig(config).apply(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Service) CreateTokens(id string) (Tokens, error) {
	accessToken, err := createToken(id, s.accessTokenTTL, s.accessTokenSigningKey)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := createToken(id, s.refreshTokenTTL, s.refreshTokenSigningKey)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func createToken(id string, ttl time.Duration, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *Service) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %#v", token.Header["alg"])
		}

		return []byte(s.accessTokenSigningKey), nil
	})
	if err != nil {
		return "", err
	}

	return parseToken(token)
}

func (s *Service) ParseRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %#v", token.Header["alg"])
		}

		return []byte(s.refreshTokenSigningKey), nil
	})
	if err != nil {
		return "", err
	}

	return parseToken(token)
}

func parseToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("error get claims from token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("error get subject from token")
	}

	return sub, nil
}
