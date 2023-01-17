package token

import (
	"crypto/rsa"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

//go:generate mockery --dir . --name JWTService --output ./mocks
type JWTService interface {
	CreateToken(id string) (jwt, uid string, err error)
	ParseToken(jwt string) (*JWTCustomClaims, error)
	TokenTTL() time.Duration
}

type jwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttl        time.Duration
}

func NewJWTService(cfg Config) (JWTService, error) {
	var s jwtService
	if err := SetConfig(cfg).apply(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

type JWTCustomClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

func (c *JWTCustomClaims) Valid() error {
	if c.UID == "" {
		return errors.New("token has empty UID")
	}

	return c.StandardClaims.Valid()
}

func (c *JWTCustomClaims) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", c.UID)
	enc.AddString("sub", c.Subject)

	enc.AddInt64("exp", c.ExpiresAt)
	enc.AddInt64("iat", c.IssuedAt)
	enc.AddInt64("nbf", c.NotBefore)

	return nil
}

func (s *jwtService) CreateToken(id string) (tokenString, uid string, err error) {
	uid = uuid.New().String()
	now := time.Now().UTC()

	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodRS256, &JWTCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			Subject:   id,
			ExpiresAt: now.Add(s.ttl).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}).SignedString(s.privateKey)

	return tokenString, uid, err
}

func (s *jwtService) ParseToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (i any, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %#v", token.Header["alg"])
		}

		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("error get claims from token")
	}

	return claims, nil
}

func (s *jwtService) TokenTTL() time.Duration { return s.ttl }
