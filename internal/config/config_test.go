package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"

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
				"ENVIRONMENT": "testing",
				"HTTP_PORT":   "1308",
			},
			args: args{
				fixture: "testdata",
			},
			ret: ret{
				config: config.Config{
					Environment: "testing",
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
					Logger: logger.Config{
						Level: "debug",
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
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
			}

			assert.Equal(t, tc.ret.config, cfg)
		})
	}
}

func shadowEnv(t *testing.T) func() {
	t.Helper()

	environ := map[string]string{}

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, config.EnvPrefix) {
			key, value, _ := strings.Cut(env, "=")
			environ[key] = value

			if err := os.Unsetenv(key); err != nil {
				t.Fatalf("could not shadow env %s: %s", key, err)
			}

			t.Logf("shadow env %s", key)
		}
	}

	return func() {
		for key, value := range environ {
			if err := os.Setenv(key, value); err != nil {
				t.Fatalf("could not restore env %s: %s", key, err)
			}

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
