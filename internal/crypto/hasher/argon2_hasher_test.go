package hasher

import "testing"

func TestArgon2Hasher(t *testing.T) {
	h := NewArgon2Hasher([]byte("examplepepper"))
	testPasswordHasher(t, h)
}
