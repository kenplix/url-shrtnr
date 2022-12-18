package validator_test

import (
	"strings"
	"testing"

	"github.com/kenplix/url-shrtnr/internal/controller/http/validator"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

func TestUsernameValidation(t *testing.T) {
	type testSchema struct {
		Username string `binding:"username"`
	}

	type args struct {
		schema testSchema
	}

	type ret struct {
		hasErr bool
	}

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "username too short",
			args: args{
				schema: testSchema{
					Username: strings.Repeat("x", 4),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username too long",
			args: args{
				schema: testSchema{
					Username: strings.Repeat("x", 33),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username has digit as first character",
			args: args{
				schema: testSchema{
					Username: "1" + strings.Repeat("x", 4),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username has underscore as first character",
			args: args{
				schema: testSchema{
					Username: "_" + strings.Repeat("x", 4),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username has double underscore inside",
			args: args{
				schema: testSchema{
					Username: "xx__xx",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username has underscore as last character",
			args: args{
				schema: testSchema{
					Username: strings.Repeat("x", 4) + "_",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "username has special character inside",
			args: args{
				schema: testSchema{
					Username: "xx_$_xx",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "ok",
			args: args{
				schema: testSchema{
					Username: "kenplix",
				},
			},
			ret: ret{
				hasErr: false,
			},
		},
	}

	t.Parallel()

	initValidator(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(tc.args.schema)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
			}
		})
	}

	t.Run("wrong field type for this binding", func(t *testing.T) {
		type testSchema struct {
			Username int `binding:"username"`
		}

		assert.Panics(t, func() {
			_ = binding.Validator.ValidateStruct(testSchema{})
		})
	})
}

func TestPasswordValidation(t *testing.T) {
	type testSchema struct {
		Password string `binding:"password"`
	}

	type args struct {
		schema testSchema
	}

	type ret struct {
		hasErr bool
	}

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "password too short",
			args: args{
				schema: testSchema{
					Password: strings.Repeat("x", 7),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "password too long",
			args: args{
				schema: testSchema{
					Password: strings.Repeat("x", 65),
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "password without uppercase letter",
			args: args{
				schema: testSchema{
					Password: "1we$rty2",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "password without lowercase letter",
			args: args{
				schema: testSchema{
					Password: "1WE$RTY2",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "password without digit",
			args: args{
				schema: testSchema{
					Password: "!wE$Rty*",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "password without special character",
			args: args{
				schema: testSchema{
					Password: "1wE3Rty2",
				},
			},
			ret: ret{
				hasErr: true,
			},
		},
		{
			name: "ok",
			args: args{
				schema: testSchema{
					Password: "1wE$Rty2",
				},
			},
			ret: ret{
				hasErr: false,
			},
		},
	}

	t.Parallel()

	initValidator(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(tc.args.schema)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
			}
		})
	}

	t.Run("wrong field type for this binding", func(t *testing.T) {
		type testSchema struct {
			Password int `binding:"password"`
		}

		assert.Panics(t, func() {
			_ = binding.Validator.ValidateStruct(testSchema{})
		})
	})
}

func initValidator(t *testing.T) {
	t.Helper()

	_, err := validator.Init(testLogger(t))
	if err != nil {
		t.Fatal("failed to initialize validator")
	}
}

func testLogger(t *testing.T) *zap.Logger {
	logger, err := log.NewLogger(log.SetLevel(zap.DebugLevel.String()))
	if err != nil {
		t.Fatalf("failed to create testing logger: %s", err)
	}

	return logger
}
