package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/kenplix/url-shrtnr/internal/service"

	"github.com/kenplix/url-shrtnr/internal/config"
	"github.com/kenplix/url-shrtnr/internal/repository"
	"github.com/kenplix/url-shrtnr/pkg/hash"
	"github.com/kenplix/url-shrtnr/pkg/hash/argon2"
	"github.com/kenplix/url-shrtnr/pkg/hash/bcrypt"
	"github.com/kenplix/url-shrtnr/pkg/httpserver"
	"github.com/kenplix/url-shrtnr/pkg/token"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	t.Cleanup(shadowEnv(t))

	type args struct {
		fixture string
	}

	type ret struct {
		config config.Config
		hasErr bool
	}

	testCases := []struct {
		name    string
		environ map[string]string
		args    args
		ret     ret
	}{
		{
			name: "testing environment",
			environ: map[string]string{
				"ENVIRONMENT":                 "testing",
				"HTTP_PORT":                   "1308",
				"JWT_ACCESSTOKEN_PRIVATEKEY":  "<access token private key>",
				"JWT_ACCESSTOKEN_PUBLICKEY":   "<access token public key>",
				"JWT_REFRESHTOKEN_PRIVATEKEY": "<refresh token private key>",
				"JWT_REFRESHTOKEN_PUBLICKEY":  "<refresh token public key>",
			},
			args: args{
				fixture: "testdata",
			},
			ret: ret{
				config: config.Config{
					Environment: config.TestingEnvironment,
					HTTP: httpserver.Config{
						Port:            "1308",
						ReadTimeout:     1 * time.Second,
						WriteTimeout:    3 * time.Second,
						IdleTimeout:     0 * time.Second,
						ShutdownTimeout: 8 * time.Second,
					},
					Database: repository.Config{
						Use: "mongodb",
					},
					Logger: log.Config{
						Level:    "debug",
						Encoding: "console",
						Mode:     "development",
					},
					Hasher: hash.Config{
						Use: "argon2",
						Bcrypt: bcrypt.Config{
							Cost: 4,
						},
						Argon2: argon2.Config{
							Memory:      1024,
							Iterations:  18,
							Parallelism: 2,
							SaltLength:  16,
							KeyLength:   16,
						},
					},
					JWT: service.JWTServiceConfig{
						AccessToken: token.Config{
							PrivateKey: "<access token private key>",
							PublicKey:  "<access token public key>",
							TTL:        20 * time.Minute,
						},
						RefreshToken: token.Config{
							PrivateKey: "<refresh token private key>",
							PublicKey:  "<refresh token public key>",
							TTL:        60 * time.Minute,
						},
						InactiveTimeout: 10 * time.Minute,
					},
				},
				hasErr: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setEnv(t, tc.environ)

			cfg, err := config.Read(tc.args.fixture)
			assert.Falsef(t, (err != nil) != tc.ret.hasErr, "expected error: %t, but got: %v", tc.ret.hasErr, err)
			assert.Equal(t, tc.ret.config, cfg)
		})
	}
}

func shadowEnv(t *testing.T) func() {
	t.Helper()

	environ := map[string]string{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, config.EnvPrefix) {
			continue
		}

		key, value, _ := strings.Cut(env, "=")
		environ[key] = value

		err := os.Unsetenv(key)
		require.NoErrorf(t, err, "failed to shadow env %s: %s", key, err)

		t.Logf("shadow env %s", key)
	}

	return func() {
		for key, value := range environ {
			err := os.Setenv(key, value)
			require.NoErrorf(t, err, "failed to restore env %s: %s", key, err)

			t.Logf("restore env %s", key)
		}
	}
}

func setEnv(t *testing.T, environ map[string]string) {
	t.Helper()

	for key, value := range environ {
		key = config.EnvPrefix + "_" + key
		t.Logf("set testing env %s=%s", key, value)
		t.Setenv(key, value)
	}
}
