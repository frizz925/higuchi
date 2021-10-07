package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArgon2Password(t *testing.T) {
	require := require.New(t)
	password := "password"

	h := NewArgon2Hasher([]byte("examplepepper"))
	digest, err := h.Hash(password)
	require.NoError(err)
	pw, err := ParseArgon2DigestToPassword(h, digest)
	require.NoError(err)
	require.Zero(pw.Compare(password))
}
