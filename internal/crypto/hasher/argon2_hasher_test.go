package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArgon2Hasher(t *testing.T) {
	require := require.New(t)
	password := "password"

	h := NewArgon2Hasher([]byte("examplepepper"))
	digest, err := h.Hash(password)
	require.NoError(err)
	n, err := h.Compare(password, digest)
	require.NoError(err)
	require.Zero(n)
}
