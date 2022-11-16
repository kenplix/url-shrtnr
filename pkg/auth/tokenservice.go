package auth

import (
	"crypto/rsa"
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
	accessServ  tokenService
	refreshServ tokenService
}

func NewTokensService(config Config) (TokensService, error) {
	s := tokensService{
		accessServ: tokenService{
			ttl: defaultAccessTokenTTL,
		},
		refreshServ: tokenService{
			ttl: defaultRefreshTokenTTL,
		},
	}
	if err := SetConfig(config).apply(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *tokensService) CreateTokens(id string) (Tokens, error) {
	accessToken, err := s.accessServ.createToken(id)
	if err != nil {
		return Tokens{}, errors.Wrapf(err, "failed to create access token")
	}

	refreshToken, err := s.refreshServ.createToken(id)
	if err != nil {
		return Tokens{}, errors.Wrapf(err, "failed to create refresh token")
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

type tokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttl        time.Duration
}

func (s *tokensService) ParseAccessToken(tokenString string) (string, error) {
	return s.accessServ.parseToken(tokenString)
}

func (s *tokensService) ParseRefreshToken(tokenString string) (string, error) {
	return s.refreshServ.parseToken(tokenString)
}

func (s *tokenService) createToken(id string) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: now.Add(s.ttl).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
	})

	return token.SignedString(s.privateKey)
}

func (s *tokenService) parseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %#v", token.Header["alg"])
		}

		return s.publicKey, nil
	})
	if err != nil {
		return "", err
	}

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
