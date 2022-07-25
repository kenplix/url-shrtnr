package argon2

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

func SetMemory(memory uint32) Option {
	return optionFunc(func(h *Hasher) {
		if memory != 0 {
			h.params.Memory = memory
		}
	})
}

func SetIterations(iterations uint32) Option {
	return optionFunc(func(h *Hasher) {
		if iterations != 0 {
			h.params.Iterations = iterations
		}
	})
}

func SetParallelism(parallelism uint8) Option {
	return optionFunc(func(h *Hasher) {
		if parallelism != 0 {
			h.params.Parallelism = parallelism
		}
	})
}

func SetSaltLength(saltLength uint32) Option {
	return optionFunc(func(h *Hasher) {
		if saltLength != 0 {
			h.params.SaltLength = saltLength
		}
	})
}

func SetKeyLength(keyLength uint32) Option {
	return optionFunc(func(h *Hasher) {
		if keyLength != 0 {
			h.params.KeyLength = keyLength
		}
	})
}
