package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTokensService(t *testing.T) {
	type args struct {
		config Config
	}

	type ret struct {
		tokensServ TokensService
		hasErr     bool
	}

	testDurationPtr := func(t *testing.T, duration time.Duration) *time.Duration {
		t.Helper()

		return &duration
	}

	testRSAKey := func(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
		t.Helper()

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("failed to generate RSA keypair: %s", err)
		}

		return privateKey, &privateKey.PublicKey
	}

	var (
		accessTokenPrivateKey, accessTokenPublicKey   = testRSAKey(t)
		refreshTokenPrivateKey, refreshTokenPublicKey = testRSAKey(t)

		encodedAccessTokenPrivateKey = encodeRSAPrivateKey(t, accessTokenPrivateKey)
		encodedAccessTokenPublicKey  = encodeRSAPublicKey(t, accessTokenPublicKey)

		encodedRefreshTokenPrivateKey = encodeRSAPrivateKey(t, refreshTokenPrivateKey)
		encodedRefreshTokenPublicKey  = encodeRSAPublicKey(t, refreshTokenPublicKey)
	)

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "empty access token private key",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: "",
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "empty refresh token private key",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: "",
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "empty access token public key",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  "",
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "empty refresh token public key",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  "",
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "negative access token TTL",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, -defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "negative refresh token TTL",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, -defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "access token with greater TTL",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL+1),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultRefreshTokenTTL),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "refresh token with less TTL",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, defaultAccessTokenTTL-1),
					},
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "ok",
			args: args{
				config: Config{
					AccessToken: TokenConfig{
						PrivateKey: encodedAccessTokenPrivateKey,
						PublicKey:  encodedAccessTokenPublicKey,
						TTL:        testDurationPtr(t, 15*time.Minute),
					},
					RefreshToken: TokenConfig{
						PrivateKey: encodedRefreshTokenPrivateKey,
						PublicKey:  encodedRefreshTokenPublicKey,
						TTL:        testDurationPtr(t, 60*time.Minute),
					},
				},
			},
			ret: ret{
				tokensServ: &tokensService{
					accessServ: tokenService{
						privateKey: accessTokenPrivateKey,
						publicKey:  accessTokenPublicKey,
						ttl:        15 * time.Minute,
					},
					refreshServ: tokenService{
						privateKey: refreshTokenPrivateKey,
						publicKey:  refreshTokenPublicKey,
						ttl:        60 * time.Minute,
					},
				},
				hasErr: false,
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokensServ, err := NewTokensService(tc.args.config)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokensServ, tokensServ)
		})
	}
}

func encodeRSAPrivateKey(t *testing.T, pk *rsa.PrivateKey) string {
	t.Helper()

	var buf bytes.Buffer

	err := pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk),
	})
	if err != nil {
		t.Fatalf("failed to encode private pem: %s", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func encodeRSAPublicKey(t *testing.T, pk *rsa.PublicKey) string {
	t.Helper()

	pkBytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		t.Fatalf("failed to marshal public key: %s", err)
	}

	var buf bytes.Buffer

	err = pem.Encode(&buf, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkBytes,
	})
	if err != nil {
		t.Fatalf("failed to encode public pem: %s", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
