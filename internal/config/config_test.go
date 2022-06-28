package config_test

import (
	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
				"HTTP_ADDR":   "addr from env variable",
			},
			args: args{
				fixture: "testdata",
			},
			ret: ret{
				config: &config.Config{
					Environment: "testing",
					HTTP: &httpserver.Config{
						Addr: "addr from env variable",
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
