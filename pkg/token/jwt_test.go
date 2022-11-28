package token

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

func TestNewJWTService(t *testing.T) {
	type args struct {
		config Config
	}

	type ret struct {
		jwtServ JWTService
		hasErr  bool
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
		privateKey, publicKey = testRSAKey(t)

		encodedPrivateKey = encodeRSAPrivateKey(t, privateKey)
		encodedPublicKey  = encodeRSAPublicKey(t, publicKey)
	)

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "empty private key",
			args: args{
				config: Config{
					PrivateKey: "",
					PublicKey:  encodedPublicKey,
					TTL:        time.Minute,
				},
			},
			ret: ret{
				jwtServ: nil,
				hasErr:  true,
			},
		},
		{
			name: "empty public key",
			args: args{
				config: Config{
					PrivateKey: encodedPrivateKey,
					PublicKey:  "",
					TTL:        time.Minute,
				},
			},
			ret: ret{
				jwtServ: nil,
				hasErr:  true,
			},
		},
		{
			name: "negative TTL",
			args: args{
				config: Config{
					PrivateKey: encodedPrivateKey,
					PublicKey:  encodedPublicKey,
					TTL:        -time.Minute,
				},
			},
			ret: ret{
				jwtServ: nil,
				hasErr:  true,
			},
		},
		{
			name: "ok",
			args: args{
				config: Config{
					PrivateKey: encodedPrivateKey,
					PublicKey:  encodedPublicKey,
					TTL:        time.Minute,
				},
			},
			ret: ret{
				jwtServ: &jwtService{
					privateKey: privateKey,
					publicKey:  publicKey,
					ttl:        time.Minute,
				},
				hasErr: false,
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jwtServ, err := NewJWTService(tc.args.config)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.jwtServ, jwtServ)
		})
	}
}

func encodeRSAPrivateKey(t *testing.T, privateKey *rsa.PrivateKey) string {
	t.Helper()

	var buf bytes.Buffer

	err := pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		t.Fatalf("failed to encode private pem: %s", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func encodeRSAPublicKey(t *testing.T, publicKey *rsa.PublicKey) string {
	t.Helper()

	pkixPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %s", err)
	}

	var buf bytes.Buffer

	err = pem.Encode(&buf, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixPublicKey,
	})
	if err != nil {
		t.Fatalf("failed to encode public pem: %s", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
