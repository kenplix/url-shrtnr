package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// TokensService provides logic for JWT & Refresh tokens generation and parsing.
//
//go:generate mockery --dir . --name TokensService --output ./mocks
type TokensService interface {
	CreateTokens(id string) (Tokens, error)
	ParseAccessToken(token string) (string, error)
	ParseRefreshToken(token string) (string, error)
}

type tokensService struct {
	accessTokenSigningKey  string
	accessTokenTTL         time.Duration
	refreshTokenSigningKey string
	refreshTokenTTL        time.Duration
}

func NewTokensService(config Config) (TokensService, error) {
	s := tokensService{
		accessTokenTTL:  defaultAccessTokenTTL,
		refreshTokenTTL: defaultRefreshTokenTTL,
	}
	if err := SetConfig(config).apply(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *tokensService) CreateTokens(id string) (Tokens, error) {
	accessToken, err := createToken(id, s.accessTokenTTL, s.accessTokenSigningKey)
	if err != nil {
		return Tokens{}, errors.Wrapf(err, "failed to create access token")
	}

	refreshToken, err := createToken(id, s.refreshTokenTTL, s.refreshTokenSigningKey)
	if err != nil {
		return Tokens{}, errors.Wrapf(err, "failed to create refresh token")
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func createToken(id string, ttl time.Duration, signingKey string) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: now.Add(ttl).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
	})

	return token.SignedString([]byte(signingKey))
}

func (s *tokensService) ParseAccessToken(tokenString string) (string, error) {
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

func (s *tokensService) ParseRefreshToken(tokenString string) (string, error) {
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
