package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testPasswordHasher(t *testing.T, h PasswordHasher) {
	require := require.New(t)
	password := "password"

	digest, err := h.Hash(password)
	require.NoError(err)
	n, err := h.Compare(password, digest)
	require.NoError(err)
	require.Zero(n)
}
