package bcrypt

import "golang.org/x/crypto/bcrypt"

// Option configures a Hasher.
type Option interface {
	apply(h *Hasher)
}

type optionFunc func(h *Hasher)

func (fn optionFunc) apply(h *Hasher) {
	fn(h)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(h *Hasher) {
		for _, option := range options {
			option.apply(h)
		}
	})
}

func SetCost(cost int) Option {
	return optionFunc(func(h *Hasher) {
		if cost >= bcrypt.MinCost && cost <= bcrypt.MaxCost {
			h.cost = cost
		}
	})
}
