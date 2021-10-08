package auth

import (
	"testing"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/stretchr/testify/require"
)

func TestArgon2Auth(t *testing.T) {
	require := require.New(t)
	password := "password"
	tempFile := "/tmp/argon2_passwd.txt"

	h := hasher.NewArgon2Hasher([]byte("pepper"))
	ad, err := hasher.NewArgon2Digest(h, password)
	require.NoError(err)

	aa := NewArgon2Auth(h)
	require.NoError(aa.WritePasswordsFile(tempFile, Argon2Users{"user": ad}))

	users, err := aa.ReadPasswordsFile(tempFile)
	require.NoError(err)
	require.Contains(users, "user")
	require.Zero(users["user"].Compare(password))
}
