package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenService(t *testing.T) {
	t.Parallel()

	type args struct {
		config Config
	}

	type ret struct {
		tokenServ *Service
		hasErr    bool
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
					RefreshTokenSigningKey: "<refresh token signing key>",
				},
			},
			ret: ret{
				tokenServ: nil,
				hasErr:    true,
			},
		},
		{
			name: "empty refresh token signing key",
			args: args{
				config: Config{
					AccessTokenSigningKey:  "<access token signing key>",
					RefreshTokenSigningKey: "",
				},
			},
			ret: ret{
				tokenServ: nil,
				hasErr:    true,
			},
		},
		{
			name: "negative access token TTL",
			args: args{
				config: Config{
					AccessTokenSigningKey:  "<access token signing key>",
					AccessTokenTTL:         testDurationPtr(t, -time.Second),
					RefreshTokenSigningKey: "<refresh token signing key>",
					RefreshTokenTTL:        testDurationPtr(t, time.Second),
				},
			},
			ret: ret{
				tokenServ: nil,
				hasErr:    true,
			},
		},
		{
			name: "negative refresh token TTL",
			args: args{
				config: Config{
					AccessTokenSigningKey:  "<access token signing key>",
					AccessTokenTTL:         testDurationPtr(t, time.Second),
					RefreshTokenSigningKey: "<refresh token signing key>",
					RefreshTokenTTL:        testDurationPtr(t, -time.Second),
				},
			},
			ret: ret{
				tokenServ: nil,
				hasErr:    true,
			},
		},
		{
			name: "ok",
			args: args{
				config: Config{
					AccessTokenSigningKey:  "<access token signing key>",
					AccessTokenTTL:         testDurationPtr(t, time.Second),
					RefreshTokenSigningKey: "<refresh token signing key>",
					RefreshTokenTTL:        testDurationPtr(t, time.Second),
				},
			},
			ret: ret{
				tokenServ: &Service{
					accessTokenSigningKey:  "<access token signing key>",
					accessTokenTTL:         time.Second,
					refreshTokenSigningKey: "<refresh token signing key>",
					refreshTokenTTL:        time.Second,
				},
				hasErr: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenServ, err := NewTokenService(tc.args.config)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokenServ, tokenServ)
		})
	}
}
