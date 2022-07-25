package argon2

import (
	"github.com/alexedwards/argon2id"
)

type Hasher struct {
	params *argon2id.Params
}

func NewHasher(options ...Option) *Hasher {
	var (
		params = *argon2id.DefaultParams
		hasher = Hasher{
			params: &params,
		}
	)

	Preset(options...).apply(&hasher)

	return &hasher
}

func (h *Hasher) HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, h.params)
}

func (h *Hasher) CheckPasswordHash(password, hash string) bool {
	match, _ := argon2id.ComparePasswordAndHash(password, hash)
	return match
}
