package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testPasswordParseDigest(t *testing.T, h PasswordHasher) {
	require := require.New(t)
	password := "password"

	digest, err := h.Hash(password)
	require.NoError(err)
	pw, err := h.ParseDigest(digest)
	require.NoError(err)
	require.Zero(pw.Compare(password))
}
