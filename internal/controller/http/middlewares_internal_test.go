package http

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Kenplix/url-shrtnr/internal/controller/http/ginctx"

	"github.com/go-playground/assert/v2"
	ut "github.com/go-playground/universal-translator"

	"github.com/Kenplix/url-shrtnr/internal/service"
	servMocks "github.com/Kenplix/url-shrtnr/internal/service/mocks"
)

func TestTranslatorMiddleware(t *testing.T) {
	type args struct {
		localeParameter      []string
		acceptLanguageHeader []string
	}

	type ret struct {
		chosenLocale string
	}

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "default locale",
			args: args{
				localeParameter:      []string{"ua"},
				acceptLanguageHeader: []string{},
			},
			ret: ret{
				chosenLocale: "en",
			},
		},
		{
			name: "different locales",
			args: args{
				localeParameter:      []string{},
				acceptLanguageHeader: []string{"ru", "en", "ca"},
			},
			ret: ret{
				chosenLocale: "ru",
			},
		},
		{
			name: "priority choice",
			args: args{
				localeParameter:      []string{"ru"},
				acceptLanguageHeader: []string{"en"},
			},
			ret: ret{
				chosenLocale: "ru",
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := testGinContext(t, rec)

			query := url.Values{"locale": tc.args.localeParameter}
			c.Request.URL.RawQuery = query.Encode()

			for _, acceptLanguage := range tc.args.acceptLanguageHeader {
				c.Request.Header.Add("Accept-Language", acceptLanguage)
			}

			var (
				jwtServ   = servMocks.NewJWTService(t)
				authServ  = servMocks.NewAuthService(t)
				usersServ = servMocks.NewUsersService(t)
			)

			h, err := NewHandler(testLogger(t), &service.Services{
				JWT:   jwtServ,
				Auth:  authServ,
				Users: usersServ,
			})
			if err != nil {
				t.Fatalf("failed to create handler: %s", err)
			}

			translatorMiddleware(h.unitrans)(c)

			translator := c.MustGet(ginctx.TranslatorContext).(ut.Translator)
			assert.Equal(t, tc.ret.chosenLocale, translator.Locale())
		})
	}
}
