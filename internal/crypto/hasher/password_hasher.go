package hasher

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, digest string) (int, error)
	ParseDigest(digest string) (PasswordDigest, error)
}
