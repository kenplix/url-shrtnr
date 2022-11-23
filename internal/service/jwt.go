package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/pkg/token"
)

type JWTServiceConfig struct {
	AccessToken     token.Config  `mapstructure:"accessToken"`
	RefreshToken    token.Config  `mapstructure:"refreshToken"`
	InactiveTimeout time.Duration `mapstructure:"inactiveTimeout"`
}

type jwtService struct {
	cache          *redis.Client
	accessServ     token.JWTService
	refreshServ    token.JWTService
	signOutTimeout time.Duration
}

func NewJWTService(cfg JWTServiceConfig, cache *redis.Client) (JWTService, error) {
	accessServ, err := token.NewJWTService(cfg.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create access token service")
	}

	refreshServ, err := token.NewJWTService(cfg.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create refresh token service")
	}

	if accessServ.TokenTTL() >= refreshServ.TokenTTL() {
		return nil, fmt.Errorf(
			"access token TTL [%s] must be less than refresh token TTL [%s]",
			accessServ.TokenTTL(),
			refreshServ.TokenTTL(),
		)
	}

	signOutTimeout := refreshServ.TokenTTL()
	if cfg.InactiveTimeout >= accessServ.TokenTTL() {
		return nil, fmt.Errorf(
			"inactive timeout [%s] must be less than access token TTL [%s] if provided",
			cfg.InactiveTimeout,
			accessServ.TokenTTL(),
		)
	} else if cfg.InactiveTimeout > 0 {
		signOutTimeout = cfg.InactiveTimeout
	}

	if cache == nil {
		return nil, errors.New("cache not provided")
	}

	s := jwtService{
		cache:          cache,
		accessServ:     accessServ,
		refreshServ:    refreshServ,
		signOutTimeout: signOutTimeout,
	}

	return &s, nil
}

func (s *jwtService) CreateTokens(ctx context.Context, userID string) (entity.Tokens, error) {
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

	s.cache.Set(ctx, tokenCacheKey(userID), cacheJSON, s.signOutTimeout)

	tokens := entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func (s *jwtService) ProlongTokens(ctx context.Context, userID string) {
	if s.signOutTimeout == s.refreshServ.TokenTTL() {
		log.Printf("debug: auto-sign-out feature not enabled")
		return
	}

	tokenKey := tokenCacheKey(userID)

	ttl := s.cache.TTL(ctx, tokenKey).Val()
	s.cache.Expire(ctx, tokenKey, s.signOutTimeout)

	log.Printf("debug: user[id:%q]: tokens ttl %s: tokens prolonged on %s", userID, ttl, s.signOutTimeout)
}

func (s *jwtService) ParseAccessToken(tokenString string) (*token.JWTCustomClaims, error) {
	return s.accessServ.ParseToken(tokenString)
}

func (s *jwtService) ParseRefreshToken(tokenString string) (*token.JWTCustomClaims, error) {
	return s.refreshServ.ParseToken(tokenString)
}

func (s *jwtService) ValidateAccessToken(ctx context.Context, claims *token.JWTCustomClaims) error {
	return s.validateToken(ctx, claims, false)
}

func (s *jwtService) ValidateRefreshToken(ctx context.Context, claims *token.JWTCustomClaims) error {
	return s.validateToken(ctx, claims, true)
}

func (s *jwtService) validateToken(ctx context.Context, claims *token.JWTCustomClaims, isRefreshToken bool) error {
	tokenKey := tokenCacheKey(claims.Subject)

	cacheBytes, err := s.cache.Get(ctx, tokenKey).Bytes()
	if err != nil {
		return errors.Wrapf(err, "cache: failed to get %q key", tokenKey)
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
