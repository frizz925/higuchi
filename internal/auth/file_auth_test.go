package auth

import (
	"testing"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/stretchr/testify/require"
)

func TestFileAuth(t *testing.T) {
	require := require.New(t)
	password := "password"
	tempFile := "/tmp/auth_passwd.txt"

	h := hasher.NewMD5Hasher([]byte("pepper"))
	ad, err := hasher.NewMD5Digest(h, password)
	require.NoError(err)

	aa := NewFileAuth(h)
	require.NoError(aa.WritePasswordsFile(tempFile, Users{"user": ad}))

	users, err := aa.ReadPasswordsFile(tempFile)
	require.NoError(err)
	require.Contains(users, "user")
	require.Zero(users["user"].Compare(password))
}
