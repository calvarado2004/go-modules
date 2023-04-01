package toolkit

import "crypto/rand"

const randonStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Tools is a toolkit for general purpose
type Tools struct{}

// RandomString returns a random string of length n
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randonStringSource)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))

		s[i] = r[x%y]
	}

	return string(s)
}
