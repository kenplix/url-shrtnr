package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/pkg/token"
)

type tokensService struct {
	cache       *redis.Client
	accessServ  token.JWTService
	refreshServ token.JWTService
}

func NewTokensService(cache *redis.Client, accessServ, refreshServ token.JWTService) (TokensService, error) {
	if cache == nil {
		return nil, errors.New("cache not provided")
	}

	if accessServ == nil {
		return nil, errors.New("access token service not provided")
	}

	if refreshServ == nil {
		return nil, errors.New("refresh token service not provided")
	}

	if accessServ.TokenTTL() >= refreshServ.TokenTTL() {
		return nil, fmt.Errorf(
			"access token TTL [%v] must be less than refresh token TTL [%v]",
			accessServ.TokenTTL(),
			refreshServ.TokenTTL(),
		)
	}

	s := tokensService{
		cache:       cache,
		accessServ:  accessServ,
		refreshServ: refreshServ,
	}

	return &s, nil
}

func (s *tokensService) CreateTokens(ctx context.Context, userID string) (entity.Tokens, error) {
	accessToken, accessTokenUID, err := s.accessServ.CreateToken(userID)
	if err != nil {
		return entity.Tokens{}, errors.Wrapf(err, "failed to create access token")
	}

	refreshToken, refreshTokenUID, err := s.refreshServ.CreateToken(userID)
	if err != nil {
		return entity.Tokens{}, errors.Wrapf(err, "failed to create refresh token")
	}

	cacheJSON, err := json.Marshal(entity.TokensUIDs{
		AccessTokenUID:  accessTokenUID,
		RefreshTokenUID: refreshTokenUID,
	})
	if err != nil {
		return entity.Tokens{}, errors.Wrapf(err, "failed to marshal tokens UIDs")
	}

	s.cache.Set(ctx, tokenCacheKey(userID), cacheJSON, s.refreshServ.TokenTTL())

	tokens := entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func (s *tokensService) ParseAccessToken(tokenString string) (*token.JWTCustomClaims, error) {
	return s.accessServ.ParseToken(tokenString)
}

func (s *tokensService) ParseRefreshToken(tokenString string) (*token.JWTCustomClaims, error) {
	return s.refreshServ.ParseToken(tokenString)
}

func (s *tokensService) ValidateAccessToken(ctx context.Context, claims *token.JWTCustomClaims) error {
	return s.validateToken(ctx, claims, false)
}

func (s *tokensService) ValidateRefreshToken(ctx context.Context, claims *token.JWTCustomClaims) error {
	return s.validateToken(ctx, claims, true)
}

func (s *tokensService) validateToken(ctx context.Context, claims *token.JWTCustomClaims, isRefreshToken bool) error {
	cacheBytes, err := s.cache.Get(ctx, tokenCacheKey(claims.Subject)).Bytes()
	if err != nil {
		return errors.Wrapf(err, "cache: failed to get %q key", tokenCacheKey(claims.Subject))
	}

	var tokensUIDs entity.TokensUIDs

	err = json.Unmarshal(cacheBytes, &tokensUIDs)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall tokens UIDs")
	}

	tokenUID := tokensUIDs.AccessTokenUID
	if isRefreshToken {
		tokenUID = tokensUIDs.RefreshTokenUID
	}

	if tokenUID != claims.UID {
		return errors.New("token not found")
	}

	return nil
}

func tokenCacheKey(userID string) string {
	return fmt.Sprintf("token:%s", userID)
}
