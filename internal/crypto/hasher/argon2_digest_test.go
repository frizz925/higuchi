package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArgon2Digest(t *testing.T) {
	require := require.New(t)
	password := "password"

	h := NewArgon2Hasher([]byte("examplepepper"))
	ad, err := NewArgon2Digest(h, password)
	require.NoError(err)
	require.Zero(ad.Compare(password))
}

func TestArgon2ParseDigest(t *testing.T) {
	require := require.New(t)
	password := "password"

	h := NewArgon2Hasher([]byte("examplepepper"))
	digest, err := h.Hash(password)
	require.NoError(err)
	pw, err := h.ParseDigest(digest)
	require.NoError(err)
	require.Zero(pw.Compare(password))
}
