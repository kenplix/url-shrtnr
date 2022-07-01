package config_test

import (
	"testing"
	"time"

	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		fixture string
	}

	type ret struct {
		config *config.Config
		hasErr bool
	}

	setEnv := func(t *testing.T, env map[string]string) {
		t.Helper()
		for key, value := range env {
			t.Setenv(config.EnvPrefix+"_"+key, value)
		}
	}

	testCases := []struct {
		name string
		env  map[string]string
		args args
		ret  ret
	}{
		{
			name: "testing environment",
			env: map[string]string{
				"ENVIRONMENT": "testing",
				"HTTP_PORT":   "1308",
			},
			args: args{
				fixture: "testdata",
			},
			ret: ret{
				config: &config.Config{
					Environment: "testing",
					HTTP: &httpserver.Config{
						Port:            1308,
						ReadTimeout:     1 * time.Second,
						WriteTimeout:    3 * time.Second,
						IdleTimeout:     0 * time.Second,
						ShutdownTimeout: 8 * time.Second,
					},
					Logger: &logger.Config{
						Level:           "debug",
						TimestampFormat: logger.DefaultTimestampFormat,
					},
				},
				hasErr: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setEnv(t, tc.env)

			cfg, err := config.New(tc.args.fixture)
			require.Condition(t, func() bool { return (err != nil) == tc.ret.hasErr })
			assert.Equal(t, tc.ret.config, cfg)
		})
	}
}
