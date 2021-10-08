package hasher

import "testing"

func TestMD5Hasher(t *testing.T) {
	h := NewMD5Hasher([]byte("pepper"))
	testPasswordHasher(t, h)
}
