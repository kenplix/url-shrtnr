package auth

import (
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

	testAccessTokenSigningKey := func(t *testing.T) string {
		t.Helper()

		return "<access token signing key>"
	}

	testRefreshTokenSigningKey := func(t *testing.T) string {
		t.Helper()

		return "<refresh token signing key>"
	}

	testDurationPtr := func(t *testing.T, duration time.Duration) *time.Duration {
		t.Helper()

		return &duration
	}

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "empty access token signing key",
			args: args{
				config: Config{
					AccessTokenSigningKey:  "",
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
				},
			},
			ret: ret{
				tokensServ: nil,
				hasErr:     true,
			},
		},
		{
			name: "empty refresh token signing key",
			args: args{
				config: Config{
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					RefreshTokenSigningKey: "",
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
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					AccessTokenTTL:         testDurationPtr(t, -time.Second),
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
					RefreshTokenTTL:        testDurationPtr(t, time.Second),
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
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					AccessTokenTTL:         testDurationPtr(t, time.Second),
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
					RefreshTokenTTL:        testDurationPtr(t, -time.Second),
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
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					AccessTokenTTL:         testDurationPtr(t, defaultRefreshTokenTTL+time.Nanosecond),
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
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
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
					RefreshTokenTTL:        testDurationPtr(t, defaultAccessTokenTTL-time.Nanosecond),
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
					AccessTokenSigningKey:  testAccessTokenSigningKey(t),
					AccessTokenTTL:         testDurationPtr(t, 15*time.Minute),
					RefreshTokenSigningKey: testRefreshTokenSigningKey(t),
					RefreshTokenTTL:        testDurationPtr(t, 60*time.Minute),
				},
			},
			ret: ret{
				tokensServ: &tokensService{
					accessTokenSigningKey:  testAccessTokenSigningKey(t),
					accessTokenTTL:         15 * time.Minute,
					refreshTokenSigningKey: testRefreshTokenSigningKey(t),
					refreshTokenTTL:        60 * time.Minute,
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
