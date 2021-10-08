package auth

import "github.com/frizz925/higuchi/internal/crypto/hasher"

type Users map[string]hasher.PasswordDigest
