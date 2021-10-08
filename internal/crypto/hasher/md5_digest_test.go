package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMD5Digest(t *testing.T) {
	require := require.New(t)
	password := "password"

	h := NewMD5Hasher([]byte("examplepepper"))
	md, err := NewMD5Digest(h, password)
	require.NoError(err)
	require.Zero(md.Compare(password))
}

func TestMD5ParseDigest(t *testing.T) {
	h := NewArgon2Hasher([]byte("examplepepper"))
	testPasswordParseDigest(t, h)
}
